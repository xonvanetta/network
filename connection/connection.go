package connection

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/xonvanetta/network/handler"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/xonvanetta/network/packet"
)

type Handler interface {
	UUID() string
}

type connection struct {
	net.Conn

	wg       sync.WaitGroup
	uuid     string
	lastPing time.Time
	closed   bool
}

func New(conn net.Conn) *connection {
	connection := &connection{
		Conn: conn,
		uuid: uuid.New().String(),
	}

	connection.wg.Add(1)
	go func() {
		connection.read()
		connection.wg.Done()
	}()

	err := handler.Do(packet.Connecting, connection.uuid, nil)
	if err != nil {
		logrus.Errorf("failed to do packet connecting for %s: %s", connection.uuid, err)
	}
	return connection
}

func (c *connection) UUID() string {
	return c.uuid
}

func (c *connection) Close() error {
	err := handler.Do(packet.Disconnect, c.uuid, nil)
	if err != nil {
		return fmt.Errorf("failed to disconnect: %s", err)
	}
	err = c.Conn.Close()
	c.wg.Wait()
	return err
}

func (c *connection) read() {
	for {
		pk, err := packet.Read(c)
		if err != nil {
			if err == io.EOF || err == io.ErrClosedPipe {
				c.Close()
				return
			}
			c.error(err)
			continue
		}

		//Todo: do some real deadline setter
		//fmt.Println(c.uuid, pk.GetType(), pk.GetMessage())

		err = handler.Do(pk.GetType(), c.uuid, pk.GetMessage())
		if err != nil {
			logrus.Errorf("network: failed to do something: %s", err)
		}
	}
}

func (c *connection) error(err error) {
	switch t := err.(type) {
	case *net.OpError:
		//Todo: check this
		if t.Timeout() {
			logrus.Errorf("network: closing connection due to timeout: %s", t)
			err := c.Conn.Close()
			if err != nil {
				logrus.Errorf("network: failed to disconnect due timeout: %s", err)
			}
			return
		}
		if t.Temporary() {
			return
		}

	default:
		logrus.Warnf("network: unhandled error occured: %s", err)
	}
}
