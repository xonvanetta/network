package event

import (
	"fmt"
	"sync"
)

type Handler interface {
	Add(uint64, Callback)
	Do(uint64, Event) error
}

type handler struct {
	callbacks callbacks
	mutex     *sync.RWMutex
}

func NewHandler() Handler {
	return &handler{
		callbacks: make(callbacks),
		mutex:     &sync.RWMutex{},
	}
}

type Callback = func(event Event) error

type callbacks map[uint64][]Callback

func (h *handler) Add(packetType uint64, callback Callback) {
	h.mutex.Lock()
	h.callbacks[packetType] = append(h.callbacks[packetType], callback)
	h.mutex.Unlock()
}

func (h *handler) Do(packetType uint64, event Event) error {
	h.mutex.RLock()
	callbacks, ok := h.callbacks[packetType]
	h.mutex.RUnlock()
	if !ok {
		return nil
	}

	for _, callback := range callbacks {
		err := callback(event)
		if err != nil {
			return fmt.Errorf("failed to perform callback on packetType: %d, error: %w", packetType, err)
		}
	}

	return nil
}
