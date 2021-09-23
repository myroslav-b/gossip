package traficmanager

import (
	"testing"
	"time"
)

//Schrödinger test
func TestTraficManager(t *testing.T) {
	uptime := 2  //sec
	sleep := 1   //sec
	tikMs := 100 //duration in msec
	repeat := 10 //number of repetitions

	tm := new(TraficManager)
	tm.Init(uint32(uptime), uint32(uptime), uint32(sleep), uint32(sleep))

	u, s := 0, 0
	for i := 0; i < repeat*(uptime+sleep)*1000/tikMs; i++ {
		b := tm.Speak()
		if b {
			u++
		} else {
			s++
		}
		time.Sleep(time.Duration(tikMs) * time.Millisecond)
	}

	if u/s != uptime/sleep {
		t.Errorf("Bad test: %v %v %v %v", u, s, uptime, sleep)
	}
}
