package packet

import (
	"math/rand"
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	reader, writer := setup()

	err := writer.Write(15, nil)
	assert.NoError(t, err)

	packet, err := reader.Read()

	assert.NoError(t, err)
	assert.Equal(t, uint64(15), packet.GetType())
}

func TestWrite8k(t *testing.T) {
	reader, writer := setup()

	randomBytes := make([]byte, 8192)
	_, err := rand.Read(randomBytes)
	assert.NoError(t, err)

	err = writer.Write(1234, &any.Any{Value: randomBytes})
	assert.NoError(t, err)

	packet, err := reader.Read()
	assert.NoError(t, err)
	assert.Equal(t, uint64(1234), packet.GetType())
}

func TestWriteTwoPackets(t *testing.T) {
	reader, writer := setup()

	err := writer.Write(1, nil)
	assert.NoError(t, err)
	err = writer.Write(8, nil)
	assert.NoError(t, err)

	packet, err := reader.Read()
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), packet.GetType())

	packet, err = reader.Read()
	assert.NoError(t, err)
	assert.Equal(t, uint64(8), packet.GetType())
}
