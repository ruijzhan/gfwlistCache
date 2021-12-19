package controller

import (
	"strings"
	"time"

	addresslist "github.com/ruijzhan/gfwlistCache/address_list"
	"github.com/ruijzhan/gfwlistCache/gfwlist"
	"github.com/ruijzhan/gfwlistCache/logserver"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

type GFWCache struct {
	rosCache  addresslist.AddressList
	logServer logserver.LogServer
	queue     workqueue.DelayingInterface
}

func New(conf *Config) *GFWCache {
	gc := &GFWCache{
		rosCache: addresslist.New(conf.RouterOSAddr, conf.RouterOSUser, conf.RouterOSPasswd),
		logServer: logserver.New(conf.LogServerBindIP, conf.LogServerBindPort,
			logserver.HasPrefix("dns,packet question"), logserver.HasSuffix(":A:IN"), logserver.NoDuplicate()),
		queue: workqueue.NewDelayingQueue(),
	}

	return gc
}

func (gc *GFWCache) Run(stopCh <-chan struct{}) {
	klog.Infoln("Sync RouterOS address-list")
	for !gc.rosCache.Synced() {
		time.Sleep(time.Second)
	}
	klog.Infoln("RouterOS address-list synced")

	go gc.listenROSLog(stopCh, 3, gc.handleROSLog)

	go gc.runWorker()

	<-stopCh
	gc.queue.ShutDown()
}

func (gc *GFWCache) listenROSLog(stopCh <-chan struct{}, workers int, handler func(string)) {
	gc.logServer.Run(stopCh)
	for i := 0; i < workers; i++ {
		go func() {
			for line := range gc.logServer.Logs() {
				handler(line)
			}
		}()
	}
	<-stopCh
}

func (gc *GFWCache) handleROSLog(line string) {
	domain := strings.Split(line, ":")[1]
	domain = strings.TrimSpace(domain)
	domain = strings.TrimSuffix(domain, ".")
	if gfwlist.Has(domain) && !gc.rosCache.Has(domain) {
		gc.queue.Add(domain)
	}
}

func (gc *GFWCache) runWorker() {
	for gc.worker() {
	}
}

func (gc *GFWCache) worker() bool {
	obj, stopped := gc.queue.Get()
	if stopped {
		return false
	}
	domain := obj.(string)
	err := gc.rosCache.Add(domain, timeout(domain))

	if err != nil {
		if err != addresslist.ErrAlreadyHaveSuchEntry {
			klog.Errorf("add %s error %v", domain, err)
			gc.queue.AddAfter(domain, time.Second*5)
			return true
		}
	} else {
		klog.Infof("%s", domain)
	}
	gc.queue.Done(domain)
	return true
}
