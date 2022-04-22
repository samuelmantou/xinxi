package pkg

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"time"
)

type PTType int

const (
	PTTypeStart PTType = iota
	PTTypePaiDan
	PTTypeZj
	PTTypeKj
)

type Cfg struct {
	Change int `json:"change"`
	Pd int `json:"pd"`
	Zj int `json:"zj"`
	Kj int `json:"kj"`
	Start string `json:"start"`
	End string `json:"end"`
}

type PinTuan struct {
	cfg *Cfg
	lastPosition int
	pdC chan struct{}
	zjC chan struct{}
	changeNextC chan struct{}
	db *gorm.DB
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

func (p *PinTuan) Start() {
	p.getDistPidArr()
	pArr := p.getDistPidArr()
	for _, d := range pArr {
		p.startWait(d.Id)
		p.startMiss(d.Id)
		p.startLost(d.Id)
	}
}

func (p *PinTuan) Zj() {
	pArr := p.getDistPidArr()
	for _, d := range pArr {
		p.zj(d.Id)
	}
}

func (p *PinTuan) Kj()  {
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 3
	win := rand.Intn(max - min + 1) + min

	pArr := p.getDistPidArr()
	for _, d := range pArr {
		p.kj(d.Id, win)
	}
}

func (p *PinTuan) TimeTicker() {
	go p.pdTicker()
	go p.zjTicker()
	go p.changePositionTicker()
}

func (p *PinTuan) Run() {
	for {
		select {
		case <-p.pdC:
			log.Println("pd")
			p.Start()
		case <-p.zjC:
			log.Println("zj")
			p.Zj()
		}
	}
}

func New(cfg *Cfg, db *gorm.DB) *PinTuan {
	return &PinTuan{
		cfg: cfg,
		db: db,
		lastPosition: 1,
		pdC: make(chan struct{}, 100),
		zjC: make(chan struct{}, 100),
		changeNextC: make(chan struct{}, 100),
	}
}