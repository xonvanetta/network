package server

import (
	"net"
	"sync"
)

type Pool interface {
	Add(uuid string, conn net.Conn)
	All() map[string]net.Conn
	Get(uuid string) net.Conn
	Remove(uuid string)
}

type pool struct {
	//total       int
	connections map[string]net.Conn
	mutex       *sync.RWMutex
}

func NewPool() Pool {
	return &pool{
		connections: make(map[string]net.Conn),
		mutex:       &sync.RWMutex{},
	}
}

func (p *pool) Add(uuid string, conn net.Conn) {
	p.mutex.Lock()
	p.connections[uuid] = conn
	p.mutex.Unlock()
}

func (p *pool) Get(uuid string) net.Conn {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.connections[uuid]
}

func (p *pool) Remove(uuid string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	_, ok := p.connections[uuid]
	if !ok {
		return
	}
	delete(p.connections, uuid)
}

func (p *pool) All() map[string]net.Conn {
	p.mutex.RLock()
	m := make(map[string]net.Conn)
	for uid, conn := range p.connections {
		m[uid] = conn
	}
	p.mutex.RUnlock()
	return m
}
