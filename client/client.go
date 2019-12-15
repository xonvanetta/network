package client

import (
	"fmt"
	"net"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/xonvanetta/network/connection"
	"github.com/xonvanetta/network/handler"
	"github.com/xonvanetta/network/packet"
)

type Client interface {
	Send(packetType uint64, packet proto.Message) error
	Connect(addr string) error
	Disconnect() error
	UUID() string
}

//each second check latency

type client struct {
	latency time.Duration
	conn    net.Conn
	uuid    string
}

func New() Client {
	client := &client{}
	//Todo: ping handler

	return client
}

func (c *client) setup() {
	handler.Add(packet.Ping, c.ping)
}

func (c *client) ping(event handler.Event) error {
	return nil
}

//func (c *client) setUUID(packet *pb.Packet) error {
//	c.uuid = packet.GetUUID()
//	return nil
//}

func (c *client) UUID() string {
	return c.uuid
}

//func (c *client) Handle(packetType uint64, handlerFunc Callback) {
//	if handlerFunc == nil {
//		panic(fmt.Errorf("framework: no handlefunc provided"))
//	}
//	c.handlers.addHandler(packetType, handlerFunc)
//}

func (c *client) Connect(addr string) error {
	var err error
	c.conn, err = net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("framework: failed to connect: %s", err)
	}

	connection.New(c.conn)

	return nil
}

func (c *client) Disconnect() error {
	return c.conn.Close()
}

func (c *client) Latency() time.Duration {
	return c.latency
}

//check if there is any diff in state then update?
func (c *client) Send(packetType uint64, any proto.Message) error {
	if c.conn == nil {
		return fmt.Errorf("framework: client not connected")
	}
	return packet.Write(c.conn, packetType, any)
}
