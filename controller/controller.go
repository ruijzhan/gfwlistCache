package controller

import (
	"bufio"
	"io"
	"os"
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
	listSet   gfwlist.Interface
}

func New(conf *Config) *GFWCache {
	gc := &GFWCache{
		logServer: logserver.New(conf.LogServerBindIP, conf.LogServerBindPort,
			logserver.HasPrefix("dns,packet question"), logserver.HasSuffix(":A:IN"), logserver.NoDuplicate()),
		queue:   workqueue.NewDelayingQueue(),
		listSet: gfwlist.New(),
	}

	// 初始化字典树
	lf, err := os.Open(conf.ListFile)
	if err != nil {
		panic(err)
	}
	defer lf.Close()
	rder := bufio.NewReader(lf)
	listNames := make(map[string]struct{})
	for {
		line, _, err := rder.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		tokens := strings.Split(string(line), " ")
		listName, domain := tokens[0], tokens[1]
		listNames[listName] = struct{}{}
		gc.listSet.Insert(listName, domain)
	}

	//初始化 ros 缓存字典
	keys := make([]string, 0, len(listNames))
	for k := range listNames {
		keys = append(keys, k)
	}

	gc.rosCache = addresslist.New(conf.RouterOSAddr, conf.RouterOSUser, conf.RouterOSPasswd, keys)

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

	if lists := gc.listSet.ListsContain(domain); len(lists) > 0 {
		for _, list := range lists {
			gc.queue.Add(list + "$" + domain)
		}
	}

	// if gfwlist.Has(domain) && !gc.rosCache.Has("listName", domain) {
	// 	gc.queue.Add(domain)
	// }
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
	listDomain := obj.(string)
	tokens := strings.Split(listDomain, "$")
	listName, domain := tokens[0], tokens[1]
	err := gc.rosCache.Add(listName, domain, timeout(domain))

	if err != nil {
		if err != addresslist.ErrAlreadyHaveSuchEntry {
			klog.Errorf("add %s to %s error %v", domain, listName, err)
			gc.queue.AddAfter(obj, time.Second*5)
			return true
		}
	} else {
		klog.Infof("%s %s", listName, domain)
	}
	gc.queue.Done(obj)
	return true
}
