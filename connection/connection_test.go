package connection

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xonvanetta/network"
	"github.com/xonvanetta/network/packet"
)

//type errConn struct {
//	net.Conn
//	err error
//}
//
//func (c errConn) Read(p []byte) (int, error) {
//	if c.err != nil {
//		return 0, c.err
//	}
//	return c.Conn.Read(p)
//}

func TestNew(t *testing.T) {
	counter := &counter{}
	Add(network.Ping, func(event Event) error {
		counter.Inc()
		return nil
	})

	conn, server := net.Pipe()
	conn = New(conn)

	err := packet.Write(server, network.Ping, nil)
	assert.NoError(t, err)

	err = conn.Close()
	assert.NoError(t, err)
	counter.Verify(t, 1)
}

//func TestHandleError(t *testing.T) {
//	conn := connection{}
//
//	conn.error()
//}
