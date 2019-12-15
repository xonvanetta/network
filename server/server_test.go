package server

import (
	"net"
	"testing"

	"github.com/xonvanetta/network/client"

	"github.com/xonvanetta/network/packet"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/nettest"
)

func newServer(t *testing.T) (*server, net.Conn, func()) {
	listener, err := nettest.NewLocalListener("tcp")
	assert.NoError(t, err)

	server := client.New(listener).(*server)

	stop := make(chan struct{})
	go func() {
		server.Start()
		close(stop)
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	assert.NoError(t, err)

	return server, conn, func() {
		err = listener.Close()
		assert.NoError(t, err)
		<-stop
	}
}

func TestServerStart(t *testing.T) {
	server, _, closer := newServer(t)

	closer()

	allConnections := All()
	assert.Len(t, allConnections, 1)
}

func TestServerPong(t *testing.T) {
	_, conn, closer := newServer(t)

	err := packet.Write(conn, packet.Pong, nil)
	assert.NoError(t, err)

	closer()
}
