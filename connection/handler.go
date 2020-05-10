package connection

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/xonvanetta/network/connection/event"
	"github.com/xonvanetta/network/connection/packet"
)

type Handler interface {
	UUID() string

	Close() error

	Write(packetType uint64, packet proto.Message) error
}

type handler struct {
	packet packet.ReadWriter
	events event.Handler

	conn net.Conn
	uuid string

	wg    *sync.WaitGroup
	mutex *sync.Mutex

	//lastPing time.Time
	closed bool
}

func New(conn net.Conn, events event.Handler) (Handler, error) {
	h := &handler{
		conn:   conn,
		packet: packet.NewReadWriter(conn),
		uuid:   uuid.New().String(),
		events: events,
		wg:     &sync.WaitGroup{},
		mutex:  &sync.Mutex{},
	}

	err := h.events.Do(Connecting, event.New(h.uuid, nil))
	if err != nil {
		return nil, fmt.Errorf("failed to do event for connecting %s: %s", h.uuid, err)
	}

	h.wg.Add(1)
	go func() {
		h.read()
		h.wg.Done()
	}()

	return h, nil
}

func (h *handler) Write(packetType uint64, packet proto.Message) error {
	return h.packet.Write(packetType, packet)
}

func (h *handler) UUID() string {
	return h.uuid
}

func (h *handler) Close() error {
	err := h.events.Do(Disconnecting, event.New(h.uuid, nil))
	if err != nil {
		return fmt.Errorf("failed to disconnect: %s", err)
	}
	h.setClosed()
	err = h.conn.Close()
	h.wg.Wait()
	return err
}

func (h *handler) setClosed() {
	h.mutex.Lock()
	h.closed = true
	h.mutex.Unlock()
}

func (h *handler) isClosed() bool {
	h.mutex.Lock()
	b := h.closed
	h.mutex.Unlock()
	return b
}

func (h *handler) read() {
	for {
		if h.isClosed() {
			return
		}
		pk, err := h.packet.Read()
		if err != nil {
			h.error(err)
			continue
		}

		e := event.New(h.uuid, pk.GetMessage())
		fmt.Println(pk.GetType(), e)
		err = h.events.Do(pk.GetType(), e)
		if err != nil {
			logrus.Errorf("network: failed to do something: %s", err)
		}
	}
}

func (h *handler) error(err error) {
	if err == io.EOF || err == io.ErrClosedPipe {
		h.Close()
		return
	}
	switch t := err.(type) {
	case *net.OpError:
		if t.Timeout() {
			logrus.Errorf("network: closing connection due to timeout: %s", t)
			err := h.conn.Close()
			if err != nil {
				logrus.Errorf("network: failed to disconnect due timeout: %s", err)
			}
			return
		}
		if t.Temporary() {
			return
		}

	default:
		logrus.Warnf("network: unhandled error occurred: %s", err)
	}
}
