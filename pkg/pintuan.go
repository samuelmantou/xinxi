package pkg

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Cfg struct {
	Change int `yaml:"change"`
	Round int `yaml:"round"`
	Insert int `yaml:"insert"`
	Start string `yaml:"start"`
	End string `yaml:"end"`
	Url string `yaml:"url"`
}

type Reward func(winIds, refundIds []string)

type PinTuan struct {
	db *gorm.DB
	cfg *Cfg
	Gap time.Duration
	LastStatus int
	lastPosition int
	runC chan struct{}
	insertC chan struct{}
	changeNextC chan struct{}
	lock sync.Mutex
}

func (p *PinTuan) InTimeRange() bool {
	n := time.Now()
	nDay := n.Format("2006-01-02")
	layout := "2006-01-02 15:04"
	sTime, err := time.ParseInLocation(layout, fmt.Sprintf("%s %s", nDay, p.cfg.Start), time.Local)
	if err != nil {
		log.Println(err)
	}
	eTime, err := time.ParseInLocation(layout, fmt.Sprintf("%s %s", nDay, p.cfg.End), time.Local)
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

func (p *PinTuan) TimeTicker() {
	go p.pdTicker()
	go p.insertTicker()
	go p.changePositionTicker()
}

func (p *PinTuan) Run() {
	for {
		select {
		case <-p.runC:
			p.Open()
		case <-p.insertC:
			p.Insert()
		}
	}
}

func New(cfg *Cfg, db *gorm.DB) *PinTuan {
	return &PinTuan{
		cfg: cfg,
		db: db,
		lastPosition: 1,
		runC: make(chan struct{}),
		insertC: make(chan struct{}),
		changeNextC: make(chan struct{}),
	}
}