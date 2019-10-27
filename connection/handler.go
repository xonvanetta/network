package connection

import (
	"fmt"
	"sync"

	"github.com/golang/protobuf/ptypes/any"
)

var (
	mutex    sync.RWMutex
	handlers = make(callbacks)
)

//Todo: own package
type Event interface {
	UUID() string
	Any() *any.Any
}

type HandlerCallback = func(event Event) error

type callbacks map[uint64][]HandlerCallback

func Add(packetType uint64, callback HandlerCallback) {
	mutex.Lock()
	handlers[packetType] = append(handlers[packetType], callback)
	mutex.Unlock()
}

func do(packetType uint64, uuid string, any *any.Any) error {
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

type event struct {
	uuid string
	any  *any.Any
}

func (e event) UUID() string {
	return e.uuid
}

func (e event) Any() *any.Any {
	return e.any
}
