package pkg

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"testing"
	"time"
	"xinxi/pkg/model"
)

func TestA(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/xinxi?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}

	var c model.ConfigData
	if err := db.Where("`root_key` = 'xinxi'  AND `key` = ?", "xinxi2_pin_tuan_start").Find(&c).Error; err != nil {
		log.Println(err)
	}
	log.Println(c)
}

func newTimeTicker(ctx context.Context, start, end string, insert time.Duration) {
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
	var timer *time.Timer

	for {
		log.Println("时间开始")
		timer = time.NewTimer(gap)
		select {
		case <-timer.C:
			log.Println("时间结束")
			gap = insert
			if time.Now().Add(gap).Sub(eTime) > 0 {
				gap += sTime.Add(time.Hour * 24).Sub(eTime)
			}
			timer.Stop()
		case <-ctx.Done():
			log.Println("回收运行")
			timer.Stop()
			return
		}
	}
}

func TestTime(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go newTimeTicker(ctx, "8:00", "20:26", time.Minute * 2)
	time.Sleep(time.Second * 5)
	cancel()
	<-make(chan struct{})
}
