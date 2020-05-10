package client

import (
	"fmt"
	"net"
	"time"

	"github.com/xonvanetta/network/connection/event"

	"github.com/golang/protobuf/proto"
	"github.com/xonvanetta/network/connection"
)

type Client interface {
	Send(packetType uint64, packet proto.Message) error
	Connect(addr string) error
	Disconnect() error
}

//each second check latency

type client struct {
	conn connection.Handler

	events event.Handler

	latency time.Duration
}

func New(events event.Handler) Client {
	client := &client{
		events: events,
	}
	//Todo: ping handler

	return client
}

//func (c *client) setup() {
//	handler.Add(packet.Ping, c.ping)
//}
//
//func (c *client) ping(event event.Event) error {
//	return nil
//}

//func (c *client) Handle(packetType uint64, handlerFunc Callback) {
//	if handlerFunc == nil {
//		panic(fmt.Errorf("framework: no handlefunc provided"))
//	}
//	c.handlers.addHandler(packetType, handlerFunc)
//}

func (c *client) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("client: failed to connect: %w", err)
	}

	c.conn, err = connection.New(conn, c.events)
	if err != nil {
		return fmt.Errorf("client: failed to create new conneciton: %w", err)
	}
	return nil
}

func (c *client) Disconnect() error {
	return c.conn.Close()
}

func (c *client) Latency() time.Duration {
	return c.latency
}

func (c *client) Send(packetType uint64, any proto.Message) error {
	if c.conn == nil {
		return fmt.Errorf("client: not connected")
	}
	return c.conn.Write(packetType, any)
}
