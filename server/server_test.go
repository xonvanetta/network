package server

import (
	"net"
	"testing"

	"github.com/xonvanetta/network/connection"
	"github.com/xonvanetta/network/connection/event"

	"github.com/xonvanetta/network/connection/packet"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/nettest"
)

func newServer(t *testing.T) (*server, packet.ReadWriter, func()) {
	listener, err := nettest.NewLocalListener("tcp")
	assert.NoError(t, err)

	events := event.NewHandler()
	server := New(events).(*server)

	done := make(chan struct{})
	events.Add(connection.Connecting, func(event event.Event) error {
		close(done)
		return nil
	})

	stop := make(chan struct{})
	go func() {
		server.Start(listener)
		close(stop)
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	assert.NoError(t, err)

	<-done

	return server, packet.NewReadWriter(conn), func() {
		err = listener.Close()
		assert.NoError(t, err)
		<-stop
	}
}

func TestServerStart(t *testing.T) {
	server, _, closer := newServer(t)

	allConnections := server.pool.All()
	assert.Len(t, allConnections, 1)

	closer()
}

func TestServerPong(t *testing.T) {
	_, conn, closer := newServer(t)

	err := conn.Write(2, nil)
	assert.NoError(t, err)

	closer()
}
