package network

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newClient(t *testing.T) (*handlers, net.Conn, func()) {
	conn, client := net.Pipe()
	handlers := newHandlers()
	connection := connection{
		Conn:     conn,
		handlers: handlers,
		uuid:     "",
		closed:   false,
	}

	go connection.handle()

	return handlers, client, func() {
		conn.Close()
	}
}

func TestNewConnection(t *testing.T) {
	handlers, conn, closer := newClient(t)
	defer closer()

	var uuid string
	wait := make(chan struct{})
	handlers.addHandler(Connecting, func(packet *Packet) error {
		assert.Equal(t, Connecting, packet.GetType())
		uuid = packet.GetUUID()
		close(wait)
		return nil
	})

	connection := newConnection(conn, newHandlers())
	err := connection.disconnect()
	<-wait
	assert.NoError(t, err)
	assert.Equal(t, uuid, connection.uuid)
}

func TestHandlers(t *testing.T) {
	//handlers := newHandlers()
	//handlers.addHandler(Ping, func(packet *Packet) error {
	//	fmt.Println("called")
	//	return nil
	//})
	//
	//conn, client := net.Pipe()
	//connection := newConnection(conn, handlers)
	//
	//err := writePacket(client, connection.uuid, Ping, nil)
	//assert.NoError(t, err)
}
