package packet

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func setup() (Reader, Writer) {
	buffer := bytes.NewBuffer(nil)
	return NewReader(buffer), NewWriter(buffer)
}

// 7100 ns - 7500 ns
// 1750 ns - 1850 ns
func BenchmarkPacketRead(b *testing.B) {
	reusableReader := &reusableReader{}
	reader := NewReader(reusableReader)
	writer := NewWriter(reusableReader)

	err := writer.Write(1, nil)
	assert.NoError(b, err)

	for n := 0; n < b.N; n++ {
		_, err := reader.Read()
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
