package packet

import (
	"bytes"
	"net"
	"testing"
	"testing/iotest"

	"golang.org/x/net/nettest"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	reader := New(bytes.NewBuffer(nil))

	err := reader.WritePacket(Connecting, nil)
	assert.NoError(t, err)

	packet, err := reader.ReadPacket()
	assert.NoError(t, err)
	assert.Equal(t, Connecting, packet.GetType())
}

func TestReadOneByeReaderErrReader(t *testing.T) {
	writer := bytes.NewBuffer(nil)

	reader := NewReadWriter(writer, iotest.DataErrReader(iotest.OneByteReader(writer)))

	err := reader.WritePacket(Connecting, nil)
	assert.NoError(t, err)

	packet, err := reader.ReadPacket()
	assert.NoError(t, err)
	assert.Equal(t, Connecting, packet.GetType())
}

func TestWriteTwoPackets(t *testing.T) {
	readWriter := NewReadWriter(net.Pipe())
	done := make(chan struct{})
	go func() {
		err := readWriter.WritePacket(Connecting, nil)
		assert.NoError(t, err)
		err = readWriter.WritePacket(Connecting, nil)
		assert.NoError(t, err)
		close(done)
	}()

	packet, err := readWriter.ReadPacket()
	assert.NoError(t, err)
	assert.Equal(t, Connecting, packet.GetType())
	packet, err = readWriter.ReadPacket()
	assert.NoError(t, err)
	assert.Equal(t, Connecting, packet.GetType())
	<-done

}

func TestSendTwoPacketsThenRead(t *testing.T) {
	listener, err := nettest.NewLocalListener("tcp")
	assert.NoError(t, err)

	setup := make(chan struct{})
	var client net.Conn
	go func() {
		client, err = net.Dial("tcp", listener.Addr().String())
		assert.NoError(t, err)
		close(setup)
	}()
	server, err := listener.Accept()
	assert.NoError(t, err)
	<-setup

	readWriter := NewReadWriter(server, client)

	done := make(chan struct{})
	go func() {
		err := readWriter.WritePacket(Connecting, nil)
		assert.NoError(t, err)
		err = readWriter.WritePacket(Connecting, nil)
		assert.NoError(t, err)
		close(done)
	}()

	<-done

	packet, err := readWriter.ReadPacket()
	assert.NoError(t, err)
	assert.Equal(t, Connecting, packet.GetType())
}

type reusableReader struct {
	buf []byte
}

func (r *reusableReader) Read(p []byte) (int, error) {
	return copy(p, r.buf), nil
}

func (r *reusableReader) Write(p []byte) (int, error) {
	r.buf = make([]byte, len(p))
	return copy(r.buf, p), nil
}

// 7100 ns - 7500 ns
// 1750 ns - 1850 ns
func BenchmarkPacketRead(b *testing.B) {
	readWriter := New(&reusableReader{})

	err := readWriter.WritePacket(Connecting, nil)
	assert.NoError(b, err)

	for n := 0; n < b.N; n++ {
		_, err := readWriter.ReadPacket()
		assert.NoError(b, err)
	}
}

//func BenchmarkPacketReadOld(b *testing.B) {
//	server, client := net.Pipe()
//	defer func() {
//		client.Close()
//		server.Close()
//	}()
//	start := make(chan struct{})
//
//	var packet *Packet
//
//	go func() {
//		close(start)
//		for {
//			_, err := oldwritePacket(client, "asdasd", 10, &Packet{})
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
