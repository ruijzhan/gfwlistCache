package addresslist

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ruijzhan/routeros"
	addresslist "github.com/ruijzhan/routeros/ip/firewall/address_list"
	"k8s.io/klog/v2"
)

var (
	ErrAlreadyHaveSuchEntry = errors.New("from RouterOS device: failure: already have such entry")
)

type AddressList interface {
	Synced() bool
	Has(listName, domain string) bool
	Add(listName string, domain string, timeout string) error
	Stop()
}

func New(apiAddr, user, passwd string, listNames []string) AddressList {
	cli, err := routeros.Dial(apiAddr, user, passwd)
	if err != nil {
		panic(err)
	}

	l := &addressList{
		cli:    cli,
		cached: make(map[string]map[string]bool),
	}
	for _, listName := range listNames {
		l.cached[listName] = make(map[string]bool)
	}

	go l.sync()

	go func() {
		tk := time.NewTicker(time.Hour)
		for range tk.C {
			err := l.resync()
			if err != nil {
				klog.Errorf("resync cache failed: %v", err)
			} else {
				for _, listName := range listNames {
					klog.Infof("resynced cache: %s: %d entries", listName, len(l.cached[listName]))
				}
			}
		}
	}()

	return l
}

type addressList struct {
	cli    *routeros.Client
	synced int32
	cached map[string]map[string]bool
	mtx    sync.RWMutex
}

func (l *addressList) sync() {
	listNames := make([]string, 0, len(l.cached))
	for k := range l.cached {
		listNames = append(listNames, k)
	}
	l.mtx.Lock()
	defer l.mtx.Unlock()
	for _, listName := range listNames {
		list, err := addresslist.List(l.cli, addresslist.WithListName(listName))
		if err != nil {
			panic(err)
		}
		for _, e := range list {
			l.cached[listName][e.Address] = true
		}
	}
	atomic.StoreInt32(&l.synced, 1)
	for _, listName := range listNames {
		klog.Infof("cache synced: %s %d entries", listName, len(l.cached[listName]))
	}
}

func (l *addressList) resync() error {
	if !l.Synced() {
		return nil
	}
	newCache := make(map[string]map[string]bool)
	listNames := make([]string, 0, len(l.cached))
	for k := range l.cached {
		listNames = append(listNames, k)
	}
	for _, listName := range listNames {
		list, err := addresslist.List(l.cli, addresslist.WithListName(listName))
		if err != nil {
			return err
		}
		newCache[listName] = make(map[string]bool)

		for _, e := range list {
			newCache[listName][e.Address] = true
		}
	}

	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.cached = newCache
	return nil
}

func (l *addressList) Synced() bool {
	return atomic.LoadInt32(&l.synced) == 1
}

func (l *addressList) Has(listName, domain string) bool {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	return l.cached[listName][domain]
}

func (l *addressList) Add(listName, domain, timeout string) error {
	if l.Has(listName, domain) {
		return ErrAlreadyHaveSuchEntry
	}
	l.mtx.Lock()
	defer l.mtx.Unlock()
	err := addresslist.Add(l.cli, listName, domain, timeout, "")
	if err == nil {
		l.cached[listName][domain] = true
	} else if err.Error() == ErrAlreadyHaveSuchEntry.Error() {
		l.cached[listName][domain] = true
		return ErrAlreadyHaveSuchEntry
	}
	return err
}

func (l *addressList) Stop() {
	l.cli.Close()
	l.cached = nil
}
