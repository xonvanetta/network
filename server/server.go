package server

import (
	"fmt"
	"net"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/xonvanetta/network/connection"
	"github.com/xonvanetta/network/connection/event"
)

type Server interface {
	Send(uuid string, packetType uint64, packet proto.Message) error
	SendAll(packetType uint64, packet proto.Message) error
	Start(listener net.Listener)
}

type server struct {
	pool   Pool
	events event.Handler
}

func New(events event.Handler) Server {
	server := &server{
		pool:   NewPool(),
		events: events,
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
	s.events.Add(connection.Disconnecting, func(event event.Event) error {
		s.pool.Remove(event.UUID())
		return nil
	})
}

//func (s *server) pong(event event.Event) error {
//	conn := s.pool.Get(event.UUID())
//	fmt.Println(conn.SetDeadline(time.Now()))
//	//todo: set lastPong on conn
//	return nil
//}
//
//func (s *server) pingConnections() error {
//	for uuid, conn := range s.pool.All() {
//		err := conn.Write(packet.Ping, nil)
//		//Todo: handle disconnects
//		if err != nil {
//			return fmt.Errorf("framework: failed to ping: %s %s", uuid, err)
//		}
//	}
//
//	return nil
//}

func (s *server) Send(uuid string, packetType uint64, any proto.Message) error {
	conn := s.pool.Get(uuid)
	if conn == nil {
		return fmt.Errorf("server: connection not found: %s", uuid)
	}

	return conn.Write(packetType, any)
}

func (s *server) SendAll(packetType uint64, any proto.Message) error {
	for _, conn := range s.pool.All() {
		err := conn.Write(packetType, any)
		if err != nil {
			return fmt.Errorf("server: failed to write packet: %w", err)
		}
	}
	return nil
}

func (s *server) Start(listener net.Listener) {
	logrus.Infof("server: started server on %s", listener.Addr())
	for {
		conn, err := listener.Accept()
		if err != nil && strings.Contains(err.Error(), "use of closed network connection") {
			return
		}

		if err != nil {
			logrus.Printf("server: failed to accept tcp :%s\n", err)
			return
		}

		go func(c net.Conn) {
			conn, err := connection.New(c, s.events)
			if err != nil {
				logrus.Errorf("server: failed to create new connection :%s", err)

				return
			}
			s.pool.Add(conn.UUID(), conn)
		}(conn)
	}
}
