package addresslist

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/ruijzhan/routeros"
	addresslist "github.com/ruijzhan/routeros/ip/firewall/address_list"
)

const cacheName = "dns_cache"

var (
	ErrAlreadyHaveSuchEntry = errors.New("from RouterOS device: failure: already have such entry")
)

type AddressList interface {
	Synced() bool
	Has(string) bool
	Add(string) error
	Stop()
}

func New(apiAddr, user, passwd string) AddressList {
	cli, err := routeros.Dial(apiAddr, user, passwd)
	if err != nil {
		panic(err)
	}

	l := &addressList{
		cli:    cli,
		cached: make(map[string]bool),
	}

	go l.sync()

	return l
}

type addressList struct {
	cli    *routeros.Client
	synced int32
	cached map[string]bool
	mtx    sync.RWMutex
}

func (l *addressList) sync() {
	list, err := addresslist.List(l.cli, addresslist.WithListName(cacheName))
	if err != nil {
		panic(err)
	}
	l.mtx.Lock()
	defer l.mtx.Unlock()
	for _, e := range list {
		l.cached[e.Address] = true
	}
	atomic.StoreInt32(&l.synced, 1)
}

func (l *addressList) Synced() bool {
	return atomic.LoadInt32(&l.synced) == 1
}

func (l *addressList) Has(domain string) bool {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	return l.cached[domain]
}

func (l *addressList) Add(domain string) error {
	err := addresslist.Add(l.cli, cacheName, domain, routeros.MAX_TIMEOUT, "")
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if err == nil {
		l.cached[domain] = true
	} else if err.Error() == ErrAlreadyHaveSuchEntry.Error() {
		l.cached[domain] = true
		return ErrAlreadyHaveSuchEntry
	}
	return err
}

func (l *addressList) Stop() {
	l.cli.Close()
	l.cached = nil
}
