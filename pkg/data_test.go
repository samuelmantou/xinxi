package pkg

import (
	"log"
	"runtime"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	p := []int{1,2,3,4}
	j := 0
	for i := 0; i < len(p) / 4; i++ {
		for j = i * 4; j < i * 4 + 4; j++ {
			log.Println(j)
		}
	}
	for ; j < len(p); j++ {
		log.Println(j)
	}
}

func TestName2(t *testing.T) {
	log.Println(runtime.GOOS)
}

func TestChannel(t *testing.T) {
	tC := make(chan struct{})

	go func() {
		for {
			time.Sleep(time.Second * 1)
			tC<- struct{}{}
		}
	}()

	for {
		select {
		case <-tC:
			log.Println(1)
		default:
		}
	}
}