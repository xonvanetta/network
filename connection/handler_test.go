package connection

import (
	"net"
	"testing"

	"github.com/xonvanetta/network/connection/packet"

	"github.com/xonvanetta/network/connection/utils"

	"github.com/xonvanetta/network/connection/event"

	"github.com/stretchr/testify/assert"
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
	counter := utils.NewCounter()

	handler := event.NewHandler()
	handler.Add(Connecting, func(event event.Event) error {
		counter.Inc()
		return nil
	})

	handler.Add(Disconnecting, func(event event.Event) error {
		counter.Inc()
		return nil
	})

	c, _ := net.Pipe()
	conn, err := New(c, handler)
	assert.NoError(t, err)

	err = conn.Close()
	assert.NoError(t, err)
	counter.Verify(t, 2)
}

func TestRead(t *testing.T) {
	counter := utils.NewCounter()

	handler := event.NewHandler()
	handler.Add(2, func(event event.Event) error {
		counter.Inc()
		return nil
	})

	c, s := net.Pipe()
	conn, err := New(c, handler)
	assert.NoError(t, err)

	writer := packet.NewWriter(s)
	err = writer.Write(2, nil)
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
