package pkg

import "time"

func (p *PinTuan) pdTicker() {
	for {
		time.Sleep(time.Second * time.Duration(p.getRound()))
		if p.InTimeRange() {
			p.runC<- struct{}{}
		}
	}
}

func (p *PinTuan) insertTicker()  {
	for {
		time.Sleep(time.Second * time.Duration(p.getInsert()))
		if p.InTimeRange() {
			p.insertC<- struct{}{}
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
			time.Sleep(time.Second * time.Duration(p.getChange()))
			p.lastPosition = p.getPosition(p.lastPosition + 1)
		}
	}
}
