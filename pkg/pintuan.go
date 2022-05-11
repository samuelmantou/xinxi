package pkg

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
	"xinxi/pkg/model"
)

type Cfg struct {
	Url string `yaml:"url"`
}

type Reward func(winIds, refundIds []string)

type PinTuan struct {
	db *gorm.DB
	reload chan struct{}
	cfg *Cfg
	Gap time.Duration
	LastStatus int
	lastPosition int
	openC chan Notice
	insertC chan Notice
	changeNextC chan struct{}
	cancelInsertFn context.CancelFunc
	cancelOpenFn context.CancelFunc
	lock sync.Mutex
}

func (p *PinTuan) getConfigValue(key string) string {
	var c model.ConfigData
	if err := p.db.Where("`root_key` = 'xinxi'  AND `key` = ?", key).Find(&c).Error; err != nil {
		log.Println(err)
	}
	return c.Value
}
func (p *PinTuan) getStart() string {
	return p.getConfigValue("xinxi2_pin_tuan_start")
}

func (p *PinTuan) getEnd() string {
	return p.getConfigValue("xinxi2_pin_tuan_end")
}

func (p *PinTuan) getRound() int {
	m := p.getConfigValue("xinxi2_pin_tuan_pre")
	mm, _ := strconv.Atoi(m)
	return mm
}

func (p *PinTuan) getInsert() int {
	m := p.getConfigValue("xinxi2_pin_tuan_gap")
	mm, _ := strconv.Atoi(m)
	return mm
}

func (p *PinTuan) getChange() int {
	r := p.getRound()
	return r / 3
}

func (p *PinTuan) InTimeRange() bool {
	n := time.Now()
	nDay := n.Format("2006-01-02")
	layout := "2006-01-02 15:04"
	sTime, err := time.ParseInLocation(layout, fmt.Sprintf("%s %s", nDay, p.getStart()), time.Local)
	if err != nil {
		log.Println(err)
	}
	eTime, err := time.ParseInLocation(layout, fmt.Sprintf("%s %s", nDay, p.getEnd()), time.Local)
	if err != nil {
		log.Println(err)
	}
	if n.Sub(sTime) > 0 && eTime.Sub(n) > 0 {
		return true
	}
	return false
}

func (p *PinTuan) getPosition(i int) int {
	if i > 4 {
		return 1
	}
	return i
}

func (p *PinTuan) Insert() {
	p.lock.Lock()
	defer p.lock.Unlock()
	pArr := p.getDistPidArr()
	for _, i := range pArr {
		p.insertNew(i.Id)
		p.insertLost(i.Id)
	}
}

func (p *PinTuan) Change() {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.lastPosition == 4 {
		p.lastPosition = 1
	}else{
		p.lastPosition++
	}
}

func (p *PinTuan) reward(winIds, refundIds []string) {
	go func() {
		data := url.Values{
			"ids": winIds,
			"refundIds": refundIds,
		}
		_, err := http.PostForm(p.cfg.Url, data)
		if err != nil {
			log.Println(err)
		}
	}()
}

func (p *PinTuan) Open() {
	p.lock.Lock()
	defer p.lock.Unlock()
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 4
	win := rand.Intn(max - min + 1) + min
	pArr := p.getDistPidArr()
	for _, d := range pArr {
		p.open(d.Id, win, p.reward)
	}
}

func (p *PinTuan) Reload() {
	if p.cancelInsertFn != nil {
		p.cancelInsertFn()
	}
	var ctx context.Context
	ctx, p.cancelInsertFn = context.WithCancel(context.Background())
	newTimeTicker(ctx, p.insertC, p.getStart(), p.getEnd(), 10 * time.Second)

	if p.cancelOpenFn != nil {
		p.cancelOpenFn()
	}
	ctx, p.cancelOpenFn = context.WithCancel(context.Background())
	newTimeTicker(ctx, p.openC, p.getStart(), p.getEnd(), time.Duration(p.getRound()) * time.Minute)
}

func (p *PinTuan) Run() {
	start := p.getStart()
	end := p.getEnd()
	round := p.getRound()
	insert := p.getInsert()
	go func() {
		for {
			time.Sleep(time.Minute)
			s := p.getStart()
			if start != s {
				start = s
				p.Reload()
				continue
			}
			e := p.getEnd()
			if end != e {
				end = e
				p.Reload()
				continue
			}
			r := p.getRound()
			if round != r {
				round = r
				p.Reload()
				continue
			}
			i := p.getInsert()
			if insert != i {
				insert = i
				p.Reload()
				continue
			}
		}
	}()
	for {
		select {
		case <-p.openC:
			log.Println("open")
			p.Open()
		case <-p.insertC:
			log.Println("insert")
			p.Insert()
		case <-p.reload:
			log.Println("reload")
		}
	}
}

func New(cfg *Cfg, db *gorm.DB) *PinTuan {
	p := &PinTuan{
		cfg: cfg,
		db: db,
		reload: make(chan struct{}),
		lastPosition: 1,
		openC: make(chan Notice),
		insertC: make(chan Notice),
		changeNextC: make(chan struct{}),
	}
	go func() {
		for {
			time.Sleep(time.Minute * 2)
			log.Println("change")
			p.Change()
		}
	}()
	return p
}
