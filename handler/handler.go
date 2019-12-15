package handler

import (
	"fmt"
	"sync"

	"github.com/golang/protobuf/ptypes/any"
)

var (
	mutex    sync.RWMutex
	handlers = make(callbacks)
)

type Callback = func(event Event) error

type callbacks map[uint64][]Callback

func Add(packetType uint64, callback Callback) {
	mutex.Lock()
	handlers[packetType] = append(handlers[packetType], callback)
	fmt.Println(handlers)
	mutex.Unlock()
}

func Do(packetType uint64, uuid string, any *any.Any) error {
	event := &event{
		uuid: uuid,
		any:  any,
	}
	mutex.RLock()
	callbacks, ok := handlers[packetType]
	mutex.RUnlock()
	if !ok {
		return nil
	}

	for _, callback := range callbacks {
		err := callback(event)
		if err != nil {
			return fmt.Errorf("framework: failed to perform callback on packetType: %d, error: %w", packetType, err)
		}
	}

	return nil
}
