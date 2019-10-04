package network

import (
	"fmt"
	"io"
	"net"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type handlers struct {
	handlers map[uint64]HandlerCallbacks
}

func newHandlers() *handlers {
	return &handlers{
		handlers: make(map[uint64]HandlerCallbacks),
	}
}

type HandlerCallback = func(packet *Packet) error

type HandlerCallbacks []HandlerCallback

func (hf *HandlerCallbacks) add(handlerFunc HandlerCallback) {
	*hf = append(*hf, handlerFunc)
}

func (h *handlers) addHandler(packetType uint64, callback HandlerCallback) {
	callbacks := h.handlers[packetType]
	callbacks.add(callback)
	h.handlers[packetType] = callbacks
}

func (h *handlers) do(packetType uint64, packet *Packet) error {
	callbacks, ok := h.handlers[packetType]
	if !ok {
		return nil
	}

	for _, callback := range callbacks {
		err := callback(packet)
		if err != nil {
			return fmt.Errorf("framework: failed to perform callback on packetType: %d, error: %w", packetType, err)
		}
	}

	return nil
}

type connection struct {
	net.Conn
	//Todo: pass handlers as reference?
	handlers *handlers
	uuid     string
	closed   bool
}

func newConnection(conn net.Conn, handlers *handlers) *connection {
	connection := &connection{
		Conn:     conn,
		uuid:     uuid.New().String(),
		handlers: handlers,
	}

	err := writePacket(conn, connection.uuid, Connecting, &Packet{UUID: connection.uuid})
	if err != nil {
		connection.handleError(err)
	}
	go connection.handle()

	return connection
}

func (c *connection) handle() {
	for {
		if c.closed {
			return
		}
		packet, err := readPacket(c)
		if err != nil {
			c.handleError(err)
			continue
		}

		//Todo: do some real deadline setter
		//err = c.SetDeadline(time.Now().Add(time.Second * 5))
		//if err != nil {
		//	logrus.Errorf("network: failed to set deadline: %s", err)
		//}

		err = c.handlers.do(packet.Type, packet)
		if err != nil {
			logrus.Errorf("network: failed to do something: %s", err)
		}
	}
}

func (c *connection) handleError(err error) {
	switch t := err.(type) {
	case *net.OpError:
		//Todo: check this
		if t.Timeout() {
			logrus.Errorf("network: closing connection due to timeout: %s", t)
			err := c.disconnect()
			if err != nil {
				logrus.Errorf("network: failed to disconnect due timeout: %s", err)
			}
			return
		}
		if t.Temporary() {
			return
		}

	default:
		if err == io.EOF {
			//Client didnt disconnect, they just rage quit, must handle disconnect
			err := c.disconnect()
			if err != nil {
				logrus.Errorf("network: failed to disconnect: %s", err)
				return
			}
		}

	}
}

func (c *connection) disconnect() error {
	packet := &Packet{UUID: c.uuid}
	err := c.handlers.do(Disconnect, packet)
	if err != nil {
		return fmt.Errorf("framework: failed to perform disconnect callbacks: %s", err)
	}

	err = c.Close()
	if err != nil {
		return fmt.Errorf("framework: failed to close connection: %s", err)
	}
	c.closed = true
	return nil
}
