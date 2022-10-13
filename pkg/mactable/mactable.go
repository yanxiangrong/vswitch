package mactable

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
	"vswitch/pkg/util/log"
)

type MacTable struct {
	record    map[uint64]int
	timestamp map[uint64]int64
	ticker    *time.Ticker
	macAmount int
	mutex     *sync.Mutex
}

const (
	expireSeconds   = 60 * 5
	onceCheckAmount = 50
	timePerExpire   = 1 * time.Second
	maxMacAmount    = 8192
)

func hwAddrToUint64(mac net.HardwareAddr) uint64 {
	b := make([]byte, 8)
	copy(b[2:], mac[:6])

	return binary.BigEndian.Uint64(b)
}

func int64ToHwAddr(i uint64) net.HardwareAddr {
	mac := make([]byte, 6)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	copy(mac, b[2:])

	return mac
}

func New() *MacTable {
	log.Trace("Create new MacTable")
	var t MacTable
	t.ticker = time.NewTicker(timePerExpire)
	t.mutex = new(sync.Mutex)
	t.record = make(map[uint64]int)
	t.timestamp = make(map[uint64]int64)

	go t.expireLoop()

	return &t
}

func (t *MacTable) Update(mac net.HardwareAddr, port int) {
	if len(t.record) >= maxMacAmount {
		t.forceExpire()
	}

	i := hwAddrToUint64(mac)

	t.mutex.Lock()
	t.record[i] = port
	t.timestamp[i] = time.Now().Unix()
	t.mutex.Unlock()
}

func (t *MacTable) Find(mac net.HardwareAddr) (int, error) {
	i := hwAddrToUint64(mac)

	port, ok := t.record[i]
	if ok {
		return port, nil
	} else {
		return -1, fmt.Errorf("mac not found")
	}
}

func (t *MacTable) expireLoop() {
	count := 0
	expired := 0
	currenTs := int64(0)
	for {
		if len(t.record) < onceCheckAmount {
			currenTs = (<-t.ticker.C).Unix()

			if len(t.record) < onceCheckAmount {
				continue
			}

			count = onceCheckAmount
			expired = 0
		}

		t.mutex.Lock()
		for mac := range t.timestamp {
			if count == 0 {
				if expired <= onceCheckAmount/4 {
					t.mutex.Unlock()
					currenTs = (<-t.ticker.C).Unix()
					t.mutex.Lock()
				}
				count = onceCheckAmount
				expired = 0
			}

			if int(currenTs-t.timestamp[mac]) > expireSeconds {
				log.Trace(fmt.Sprintf("Delete old mac mapping: %s -> %d", int64ToHwAddr(mac), t.record[mac]))
				expired++
				delete(t.timestamp, mac)
				delete(t.record, mac)
			}

			count--
		}
		t.mutex.Unlock()

	}
}

func (t *MacTable) forceExpire() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	currenTs := time.Now().Unix()
	minKey := uint64(0)
	minVal := currenTs

	hasDel := false
	for mac, val := range t.timestamp {
		if int(currenTs-val) > expireSeconds/2 {
			log.Trace(fmt.Sprintf("Force delete old mac mapping: %s -> %d", int64ToHwAddr(mac), t.record[mac]))
			hasDel = true
			delete(t.timestamp, mac)
			delete(t.record, mac)
		}

		if val < minVal {
			minVal = val
			minKey = mac
		}
	}

	if !hasDel {
		log.Trace(fmt.Sprintf("Force delete old mac mapping: %s -> %d", int64ToHwAddr(minKey), t.record[minKey]))
		delete(t.timestamp, minKey)
		delete(t.record, minKey)
	}
}
