package packet

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/nettest"
)

func TestSendTwoPackets(t *testing.T) {
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

	readWriter := NewReadWriter(client, server)

	done := make(chan struct{})
	go func() {
		err := readWriter.WritePacket(Connecting, nil)
		assert.NoError(t, err)
		err = readWriter.WritePacket(Connecting, nil)
		assert.NoError(t, err)
		err = readWriter.WritePacket(Connecting, nil)
		assert.NoError(t, err)
		close(done)
	}()

	<-done

	packet, err := readWriter.ReadPacket()
	assert.NoError(t, err)
	assert.Equal(t, Connecting, packet.GetType())
	//fmt.Println(packet)
	packet, err = readWriter.ReadPacket()
	assert.NoError(t, err)
	assert.Equal(t, Connecting, packet.GetType())

	err = readWriter.WritePacket(Disconnect, nil)
	assert.NoError(t, err)

	packet, err = readWriter.ReadPacket()
	assert.NoError(t, err)
	assert.Equal(t, Connecting, packet.GetType())

	packet, err = readWriter.ReadPacket()
	assert.NoError(t, err)
	assert.Equal(t, Disconnect, packet.GetType())

	//fmt.Println(packet)
}
