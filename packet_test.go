package network

import (
	"net"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func TestPacket(t *testing.T) {
	server, client := net.Pipe()
	defer func() {
		client.Close()
		server.Close()
	}()

	start := make(chan struct{})
	done := make(chan struct{})

	go func() {
		close(start)
		packet, err := readPacket(server)
		assert.NoError(t, err)
		assert.Equal(t, Ping, packet.GetType())
		assert.Equal(t, "570465ec-2944-46c1-a019-b6f59d46e5ef", packet.GetUUID())
		close(done)
	}()
	<-start

	err := writePacket(client, "570465ec-2944-46c1-a019-b6f59d46e5ef", Ping, &Packet{})
	assert.NoError(t, err)

	<-done
}

func TestLargePacket(t *testing.T) {
	server, client := net.Pipe()
	start := make(chan struct{})
	done := make(chan struct{})

	go func() {
		close(start)
		defer close(done)
		reader := iotest.DataErrReader(iotest.OneByteReader(server))
		//reader = iotest.TimeoutReader(reader)
		packet, err := readPacket(reader)
		assert.NoError(t, err)
		assert.Equal(t, Ping, packet.GetType())
		assert.Equal(t, "570465ec-2944-46c1-a019-b6f59d46e5ef", packet.GetUUID())
	}()
	<-start

	err := writePacket(client, "570465ec-2944-46c1-a019-b6f59d46e5ef", Ping, &Packet{})
	assert.NoError(t, err)
	client.Close()
	server.Close()

	<-done
}

//var r *pb.Packet
//
//func BenchmarkPacketRead(b *testing.B) {
//	b.StopTimer()
//	server, client := net.Pipe()
//	//defer func() {
//	//	client.Close()
//	//	server.Close()
//	//}()
//	start := make(chan struct{})
//
//	var packet *pb.Packet
//
//	go func() {
//		close(start)
//		for {
//			_, err := writePacket(client, "asdasd", 10, &pb.Packet{})
//			assert.NoError(b, err)
//		}
//	}()
//	<-start
//
//	b.StartTimer()
//	for n := 0; n < b.N; n++ {
//		// always record the result of Fib to prevent
//		// the compiler eliminating the function call.
//		var err error
//		packet, err = readPacket(server)
//		assert.NoError(b, err)
//	}
//	// always store the result to a package level variable
//	// so the compiler cannot eliminate the Benchmark itself.
//	r = packet
//}
//
//func BenchmarkPacketReadOld(b *testing.B) {
//	server, client := net.Pipe()
//	//defer func() {
//	//	client.Close()
//	//	server.Close()
//	//}()
//	start := make(chan struct{})
//
//	var packet *pb.Packet
//
//	go func() {
//		close(start)
//		for {
//			_, err := oldwritePacket(client, "asdasd", 10, &pb.Packet{})
//			assert.NoError(b, err)
//		}
//	}()
//	<-start
//
//	for n := 0; n < b.N; n++ {
//		// always record the result of Fib to prevent
//		// the compiler eliminating the function call.
//		var err error
//		packet, err = oldreadPacket(server)
//		assert.NoError(b, err)
//	}
//	// always store the result to a package level variable
//	// so the compiler cannot eliminate the Benchmark itself.
//	r = packet
//}
