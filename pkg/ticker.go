package pkg

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Notice struct{}

func newTimeTicker(ctx context.Context, nC chan <-Notice, start, end string, duration time.Duration) {
	if duration <= 0 {
		return
	}
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
		gap = sTime.Add(duration).Sub(n)
	}else if t2 < 0 {
		// 已经结束
		gap = sTime.Add(time.Hour * 24).Sub(n) + duration
	}else {
		// 已经开始
		next := sTime.Add((t1 / duration - 1) * -duration)
		gap = next.Sub(n)
	}

	go func() {
		var timer *time.Timer
		for {
			log.Println("时间开始")
			timer = time.NewTimer(gap)
			select {
			case <-timer.C:
				gap = duration
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
