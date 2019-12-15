package server

import (
	"fmt"
	"net"
	"time"

	"github.com/xonvanetta/network/handler"

	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/xonvanetta/network/connection"
	"github.com/xonvanetta/network/packet"
)

type Server interface {
	Send(uuid string, packetType uint64, packet proto.Message) error
	SendAll(packetType uint64, packet proto.Message) error
	Start(listener net.Listener)
}

type server struct {
	pool Pool
}

func New() Server {
	server := &server{
		pool: NewPool(),
	}

	//go func() {
	//	ticker := time.NewTicker(time.Second * 5)
	//	for range ticker.C {
	//		err := server.pingConnections()
	//		if err != nil {
	//			logrus.Errorf("failed to ping connections: %s", err)
	//		}
	//	}
	//}()

	server.setup()
	return server
}

func (s *server) setup() {
	//handler.Add(network.Pong, s.pong)
	//handler.Add(network.Disconnect, s.disconnect)
}

func (s *server) pong(event handler.Event) error {
	conn := s.pool.Get(event.UUID())
	fmt.Println(conn.SetDeadline(time.Now()))
	//todo: set lastPong on conn
	return nil
}

func (s *server) pingConnections() error {
	for uuid, conn := range s.pool.All() {
		err := packet.Write(conn, packet.Ping, nil)
		//Todo: handle disconnects
		if err != nil {
			return fmt.Errorf("framework: failed to ping: %s %s", uuid, err)
		}
	}

	return nil
}

func (s *server) Send(uuid string, packetType uint64, any proto.Message) error {
	conn := s.pool.Get(uuid)
	if conn == nil {
		return fmt.Errorf("framework: connection not found: %s", uuid)
	}

	return packet.Write(conn, packetType, any)
}

func (s *server) SendAll(packetType uint64, any proto.Message) error {
	for _, conn := range s.pool.All() {
		err := packet.Write(conn, packetType, any)
		if err != nil {
			return fmt.Errorf("framework: failed to write packet: %s", err)
		}
	}
	return nil
}

func (s *server) Start(listener net.Listener) {
	logrus.Infof("started server on %s", listener.Addr())
	for {
		conn, err := listener.Accept()
		//Todo: handle error better. ie net op error
		if err != nil {
			logrus.Printf("network: failed to accept tcp :%s\n", err)
			return
		}

		go func(c net.Conn) {
			conn := connection.New(c)
			s.pool.Add(conn.UUID(), conn)
		}(conn)
	}
}
