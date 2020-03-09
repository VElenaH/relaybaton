package relaybaton

import (
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
)

type connectionPool struct {
	conns     [1 << 16]*net.Conn
	mutex     sync.RWMutex
	closeSent [1 << 16]bool
}

func newConnectionPool() *connectionPool {
	cp := &connectionPool{}
	i := 0
	for i < (1 << 16) {
		cp.conns[i] = nil
		i++
	}
	cp.mutex = sync.RWMutex{}
	return cp
}

func (cp *connectionPool) get(index uint16) *net.Conn {
	defer cp.mutex.RUnlock()
	cp.mutex.RLock()
	conn := cp.conns[index]
	return conn
}

func (cp *connectionPool) set(index uint16, conn *net.Conn) {
	cp.mutex.Lock()
	cp.conns[index] = conn
	cp.closeSent[index] = false
	cp.mutex.Unlock()
}

func (cp *connectionPool) delete(index uint16) {
	cp.mutex.Lock()
	if cp.conns[index] != nil {
		err := (*cp.conns[index]).Close()
		if err != nil {
			log.WithField("session", index).Warn(err)
		}
		cp.conns[index] = nil
	}
	cp.mutex.Unlock()
}

func (cp *connectionPool) isCloseSent(index uint16) bool {
	defer cp.mutex.RUnlock()
	cp.mutex.RLock()
	sent := cp.closeSent[index]
	return sent
}

func (cp *connectionPool) setCloseSent(index uint16) {
	cp.mutex.Lock()
	cp.closeSent[index] = true
	cp.mutex.Unlock()
}
