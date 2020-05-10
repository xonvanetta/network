package packet

import (
	"bytes"
	"math/rand"
	"net"
	"testing"
	"testing/iotest"

	"golang.org/x/net/nettest"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/stretchr/testify/assert"
)

func TestReadOneByteReaderErrReader(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	reader := NewReader(iotest.DataErrReader(iotest.OneByteReader(writer)))

	randomBytes := make([]byte, 8192)
	_, err := rand.Read(randomBytes)
	assert.NoError(t, err)

	p := Packet{
		Type:    3421,
		Message: &any.Any{Value: randomBytes},
	}

	buffer := new(proto.Buffer)
	err = buffer.EncodeFixed64(uint64(proto.Size(&p)))
	assert.NoError(t, err)

	err = buffer.Marshal(&p)
	assert.NoError(t, err)

	_, err = writer.Write(buffer.Bytes())
	assert.NoError(t, err)

	packet, err := reader.Read()
	assert.NoError(t, err)
	assert.Equal(t, uint64(3421), packet.GetType())
}

func TestReadMultiplePackets(t *testing.T) {
	listener, err := nettest.NewLocalListener("tcp")
	assert.NoError(t, err)

	client, err := net.Dial("tcp", listener.Addr().String())
	assert.NoError(t, err)

	server, err := listener.Accept()
	assert.NoError(t, err)

	reader := NewReader(server)
	writer := NewWriter(client)

	err = writer.Write(1, nil)
	assert.NoError(t, err)
	err = writer.Write(2, nil)
	assert.NoError(t, err)
	err = writer.Write(3, nil)
	assert.NoError(t, err)

	packet, err := reader.Read()
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), packet.GetType())

	packet, err = reader.Read()
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), packet.GetType())

	err = writer.Write(10, nil)
	assert.NoError(t, err)

	packet, err = reader.Read()
	assert.NoError(t, err)
	assert.Equal(t, uint64(3), packet.GetType())

	packet, err = reader.Read()
	assert.NoError(t, err)
	assert.Equal(t, uint64(10), packet.GetType())
}
