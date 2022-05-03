package pkg

import (
	"context"
	"time"
)

func (p *PinTuan) pdTicker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		time.Sleep(time.Second * time.Duration(p.getRound()))
		p.runC<- struct{}{}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (p *PinTuan) insertTicker(ctx context.Context)  {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		time.Sleep(time.Second * time.Duration(p.getInsert()))
		p.insertC<- struct{}{}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (p *PinTuan) changePositionTicker(ctx context.Context)  {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		time.Sleep(time.Second * time.Duration(p.getChange()))
		p.lastPosition = p.getPosition(p.lastPosition + 1)
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}
