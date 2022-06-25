package main

import (
	"sync"
	"time"
)

type Device struct {
	Addr string `json:"Addr"`
	// IP        string       `json:"ip"`
	// Port      string       `json:"port"`
	ID        string    `json:"-"`
	MM        string    `json:"-"`
	LastAlive time.Time `json:"LastAlive"`
	Lifetime  uint32    `json:"Lifetime,string"`
}

type Pool struct {
	sync.RWMutex
	DevicesMap  map[string]*Device
	DevicesList []Device // json
	count       uint32
}

func (pool *Pool) Init() {
	pool.DevicesMap = make(map[string]*Device)
	pool.count = 0

	go func() {
		for {
			time.Sleep(1 * time.Second)
			pool.ListCache()
		}
	}()

}

func (pool *Pool) Update(d *Device) {
	pool.Lock()
	r := pool.Find(d.ID)
	if r == nil {
		pool.DevicesMap[d.ID] = d
		pool.count += 1
	} else {
		r.LastAlive = time.Now()
	}
	pool.Unlock()
}

func (pool *Pool) Find(ID string) *Device {
	device, ok := pool.DevicesMap[ID]
	if ok {
		return device
	}
	return nil
}

func (pool *Pool) ListCache() {
	pool.RLock()

	var list []Device
	for k, v := range pool.DevicesMap {
		now := time.Now()
		subM := now.Sub(v.LastAlive)
		if subM.Seconds() > float64(LIFETIME) {
			delete(pool.DevicesMap, k)
			continue
		}

		list = append(list, *v)
	}
	pool.DevicesList = list
	pool.RUnlock()
}
