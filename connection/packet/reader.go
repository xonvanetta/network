package packet

import (
	"io"

	"google.golang.org/protobuf/encoding/protowire"

	"github.com/golang/protobuf/proto"
)

var (
	defaultSize = 4096

	NilPacket Packet
)

type reader struct {
	buf []byte

	reader io.Reader
}

type Reader interface {
	Read() (Packet, error)
}

func NewReader(r io.Reader) Reader {
	return &reader{
		reader: r,
	}
}

func (r *reader) Read() (Packet, error) {
	b, err := r.lowLevelRead()
	if err != nil {
		return NilPacket, err
	}

	//todo: use pool to buffer them in memory
	packet := Packet{}
	err = proto.NewBuffer(b).Unmarshal(&packet)
	return packet, err
}

func (r *reader) lowLevelRead() ([]byte, error) {
	b := make([]byte, defaultSize)
	var packetLength uint64

	var buffer []byte
	for {
		if r.buf != nil {
			b = r.buf
			r.buf = nil
		} else {
			n, err := r.reader.Read(b)
			if err != nil && n == 0 {
				return nil, err
			}
			b = b[:n]
		}

		buffer = append(buffer, b...)

		if packetLength == 0 {
			if len(buffer) < protowire.SizeFixed64() {
				continue
			}
			packetLength, _ = protowire.ConsumeFixed64(buffer[:protowire.SizeFixed64()])
			buffer = buffer[protowire.SizeFixed64():]
		}

		bytesRead := uint64(len(buffer))
		//fmt.Printf("packetLength: %d, bytesRead: %d, buffer: %v, b: %v\n", packetLength, bytesRead, len(buffer), len(b))

		//packet is longer than what we did read into the buffer
		if packetLength > bytesRead {
			b = make([]byte, packetLength-bytesRead)
			continue
		}

		//two or more packets are in the buffer
		if bytesRead > packetLength {
			r.buf = buffer[packetLength:]

			buffer = buffer[:packetLength]
		}
		//packet is the same size as what we read from the io stream
		return buffer, nil
	}

}
