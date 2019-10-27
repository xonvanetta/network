package packet

import (
	"bytes"
	"net"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/xonvanetta/network"
)

func TestRead(t *testing.T) {
	reader := bytes.NewBuffer(nil)

	err := Write(reader, network.Connecting, nil)
	assert.NoError(t, err)

	packet, err := Read(reader)
	assert.NoError(t, err)
	assert.Equal(t, network.Connecting, packet.GetType())
}

func TestReadOneByeReaderErrReader(t *testing.T) {
	reader := bytes.NewBuffer(nil)

	err := Write(reader, network.Connecting, nil)
	assert.NoError(t, err)

	packet, err := Read(iotest.DataErrReader(iotest.OneByteReader(reader)))
	assert.NoError(t, err)
	assert.Equal(t, network.Connecting, packet.GetType())
}

func TestReadTwoPackets(t *testing.T) {
	server, client := net.Pipe()

	done := make(chan struct{})
	go func() {
		packet, err := Read(client)
		assert.NoError(t, err)
		assert.Equal(t, network.Connecting, packet.GetType())
		packet, err = Read(client)
		assert.NoError(t, err)
		assert.Equal(t, network.Connecting, packet.GetType())
		close(done)
	}()

	err := Write(server, network.Connecting, nil)
	assert.NoError(t, err)
	err = Write(server, network.Connecting, nil)
	assert.NoError(t, err)
	<-done
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
	reader := &reusableReader{}

	err := Write(reader, network.Connecting, nil)
	assert.NoError(b, err)

	for n := 0; n < b.N; n++ {
		_, err := Read(reader)
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
