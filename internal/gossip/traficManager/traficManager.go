package traficmanager

import (
	"math/rand"
	"time"
)

type TraficControler interface {
	Speak() bool
}

type TraficManager struct {
	uptime    uint32
	sleep     uint32
	semaphore bool
}

func (tm *TraficManager) Speak() bool {
	//fmt.Println(tm.semaphore)
	return tm.semaphore
}

func (tm *TraficManager) Init(minUptime, maxUptime, minSleep, maxSleep uint32) {
	rand.Seed(time.Now().Unix())
	tm.uptime = minUptime + uint32(rand.Int31n(int32(maxUptime-minUptime+1)))
	tm.sleep = minSleep + uint32(rand.Int31n(int32(maxSleep-minSleep+1)))
	tm.semaphore = true
	go interrupt(tm.uptime, tm.sleep, &tm.semaphore)
}

func interrupt(uptime, sleep uint32, sem *bool) {
	for {
		if *sem {
			time.Sleep(time.Duration(uptime) * time.Second)
			*sem = !(*sem)
		} else {
			time.Sleep(time.Duration(sleep) * time.Second)
			*sem = !(*sem)
		}
	}
}
