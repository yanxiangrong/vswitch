package porttable

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"vswitch/pkg/util/log"
)

type Status int

const (
	DOWN Status = 0
	UP   Status = 1
)

const (
	maxPortAmount = 64
)

type PortTable struct {
	conn     [maxPortAmount + 1]net.Conn
	status   [maxPortAmount + 1]Status
	upAmount int
	mutex    *sync.Mutex
}

func New() *PortTable {
	log.Trace("Create new PortTable")
	var t PortTable
	t.upAmount = 0
	t.mutex = new(sync.Mutex)

	return &t
}

func (t *PortTable) getNewPort() (int, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	p := rand.Intn(maxPortAmount) + 1
	if t.status[p] == DOWN {
		return p, nil
	}
	for i := p + 1; i != p; i++ {
		if i > maxPortAmount {
			i = 1
		}
		if t.status[i] == DOWN {
			return i, nil
		}
	}

	return -1, fmt.Errorf("no free ports")
}

func (t *PortTable) Add(conn net.Conn) (int, error) {
	port, err := t.getNewPort()
	if err != nil {
		return -1, nil
	}

	t.mutex.Lock()
	t.conn[port] = conn
	t.status[port] = UP
	t.upAmount++
	t.mutex.Unlock()

	return port, nil
}

func (t *PortTable) Find(port int) (net.Conn, error) {
	conn := t.conn[port]
	status := t.status[port]
	if status == UP {
		return conn, nil
	} else {
		return nil, fmt.Errorf("port is down")
	}
}

func (t *PortTable) Delete(port int) {
	t.mutex.Lock()
	_ = t.conn[port].Close()
	t.conn[port] = nil
	t.status[port] = DOWN
	t.upAmount--
	t.mutex.Unlock()
}

func (t *PortTable) NextUpPort(port int) int {
	for i := port + 1; i != port; i++ {
		if i > maxPortAmount {
			i = 1
		}
		if t.status[i] == UP {
			return i
		}
	}
	return port
}
