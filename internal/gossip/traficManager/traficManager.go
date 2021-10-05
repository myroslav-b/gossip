package traficmanager

import (
	"math/rand"
	"sync"
	"time"
)

type TraficManager struct {
	sync.Mutex
	uptime    uint32
	sleep     uint32
	semaphore bool
}

func (tm *TraficManager) Speak() bool {
	tm.Lock()
	defer tm.Unlock()
	return tm.semaphore
}

func (tm *TraficManager) Init(minUptime, maxUptime, minSleep, maxSleep uint32) {
	tm.Lock()
	rand.Seed(time.Now().Unix())
	tm.uptime = minUptime + uint32(rand.Int31n(int32(maxUptime-minUptime+1)))
	tm.sleep = minSleep + uint32(rand.Int31n(int32(maxSleep-minSleep+1)))
	tm.semaphore = true
	tm.Unlock()
	go tm.interrupt()
}

func (tm *TraficManager) interrupt() {
	for {
		if tm.semaphore {
			time.Sleep(time.Duration(tm.uptime) * time.Second)
			tm.Lock()
			tm.semaphore = !(tm.semaphore)
			tm.Unlock()
		} else {
			time.Sleep(time.Duration(tm.sleep) * time.Second)
			tm.Lock()
			tm.semaphore = !(tm.semaphore)
			tm.Unlock()
		}
	}
}
