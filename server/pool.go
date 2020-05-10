package server

import (
	"sync"

	"github.com/xonvanetta/network/connection"
)

type Pool interface {
	Add(uuid string, conn connection.Handler)
	All() map[string]connection.Handler
	Get(uuid string) connection.Handler
	Remove(uuid string)
}

type pool struct {
	//total       int
	connections map[string]connection.Handler
	mutex       *sync.RWMutex
}

func NewPool() Pool {
	return &pool{
		connections: make(map[string]connection.Handler),
		mutex:       &sync.RWMutex{},
	}
}

func (p *pool) Add(uuid string, conn connection.Handler) {
	p.mutex.Lock()
	p.connections[uuid] = conn
	p.mutex.Unlock()
}

func (p *pool) Get(uuid string) connection.Handler {
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

func (p *pool) All() map[string]connection.Handler {
	p.mutex.RLock()
	m := make(map[string]connection.Handler)
	for uuid, conn := range p.connections {
		m[uuid] = conn
	}
	p.mutex.RUnlock()
	return m
}
