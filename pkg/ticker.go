package pkg

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Notice struct{}

func newTimeTicker(ctx context.Context, nC chan <-Notice, start, end string, insert time.Duration) {
	n := time.Now()
	nDay := n.Format("2006-01-02")
	layout := "2006-01-02 15:04"
	sTime, err := time.ParseInLocation(layout, fmt.Sprintf("%s %s", nDay, start), time.Local)
	if err != nil {
		log.Println(err)
	}
	eTime, err := time.ParseInLocation(layout, fmt.Sprintf("%s %s", nDay, end), time.Local)
	if err != nil {
		log.Println(err)
	}
	t1 := sTime.Sub(n)
	t2 := eTime.Sub(n)
	var gap time.Duration
	if t1 > 0 {
		// 还没开始
		gap = sTime.Add(insert).Sub(n)
	}else if t2 < 0 {
		// 已经结束
		gap = sTime.Add(time.Hour * 24).Sub(n) + insert
	}else {
		// 已经开始
		next := sTime.Add((t1 / insert - 1) * -insert)
		gap = next.Sub(n)
	}

	go func() {
		var timer *time.Timer
		for {
			log.Println("时间开始")
			timer = time.NewTimer(gap)
			select {
			case <-timer.C:
				gap = insert
				if time.Now().Add(gap).Sub(eTime) > 0 {
					gap += sTime.Add(time.Hour * 24).Sub(eTime)
				}
				timer.Stop()
				nC<- struct{}{}
				log.Println("时间结束")
			case <-ctx.Done():
				log.Println("回收运行")
				timer.Stop()
				return
			}
		}
	}()
}

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
