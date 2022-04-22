package pkg

import "time"

func (p *PinTuan) pdTicker() {
	for {
		time.Sleep(time.Second * time.Duration(p.cfg.Pd))
		if p.InTimeRange() {
			p.pdC<- struct{}{}
		}
	}
}

func (p *PinTuan) zjTicker()  {
	for {
		time.Sleep(time.Second * time.Duration(p.cfg.Zj))
		if p.InTimeRange() {
			p.zjC<- struct{}{}
		}
	}
}

func (p *PinTuan) changePositionTicker()  {
	defer func() {
		go p.changePositionTicker()
	}()
	for {
		select {
		case <-p.changeNextC:
			p.lastPosition = p.getPosition(p.lastPosition + 1)
			return
		default:
			time.Sleep(time.Second * time.Duration(p.cfg.Change))
		}
	}
}