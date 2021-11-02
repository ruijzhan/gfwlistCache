package logserver

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func Test_logserver_Logs(t *testing.T) {
	ls := &logserver{
		logCh: make(chan string),
	}
	count := 10000
	log := "dns local query: #8627344 gateway-carry.icloud.com. A"
	go func() {
		time.Sleep(time.Second)
		for i := 0; i < count; i++ {
			ls.logCh <- log
		}
		close(ls.logCh)
	}()
	var count2 int32
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			for range ls.Logs() {
				atomic.AddInt32(&count2, 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if count != int(count2) {
		t.Fatal()
	}
}
