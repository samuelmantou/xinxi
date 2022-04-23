package pkg

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type Cfg struct {
	Change int `json:"change"`
	Round int `json:"round"`
	Winner int `json:"winner"`
	Old int `json:"old"`
	Start string `json:"start"`
	End string `json:"end"`
}

type PinTuan struct {
	cfg *Cfg
	Gap time.Duration
	LastStatus int
	lastPosition int
	runC chan struct{}
	insertC chan struct{}
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

func (p *PinTuan) Insert() {
	pArr := p.getDistPidArr()
	for _, i := range pArr {
		p.insertNew(i.Id)
		p.insertLost(i.Id)
	}
}

func (p *PinTuan) Open() {
	//pArr := p.getDistPidArr()
	//for _, d := range pArr {
	//	p.zj(d.Id)
	//}
}

//func (p *PinTuan) Kj()  {
//	rand.Seed(time.Now().UnixNano())
//	min := 0
//	max := 3
//	win := rand.Intn(max - min + 1) + min
//
//	pArr := p.getDistPidArr()
//	for _, d := range pArr {
//		p.kj(d.Id, win)
//	}
//}

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