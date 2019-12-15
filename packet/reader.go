package packet

import (
	"fmt"
	"io"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/gogo/protobuf/proto"
)

var (
	defaultSize = 4096

	NilPacket Packet
)

type readWriter struct {
	packets [][]byte

	io.Reader
	io.Writer
}

type ReadWriter interface {
	io.Reader
	io.Writer
	ReadPacket() (Packet, error)
	WritePacket(packetType uint64, packet proto.Message) error
}

func NewReadWriter(w io.Writer, r io.Reader) ReadWriter {
	return &readWriter{
		Reader: r,
		Writer: w,
	}
}

func New(rw io.ReadWriter) ReadWriter {
	return &readWriter{
		Reader: rw,
		Writer: rw,
	}
}

func pop(slice [][]byte) ([]byte, [][]byte) {
	return slice[len(slice)-1], slice[:len(slice)-1]
}

func (rw readWriter) WritePacket(packetType uint64, packet proto.Message) error {
	var packetAny *any.Any
	if packet != nil {
		var err error
		packetAny, err = ptypes.MarshalAny(packet)
		if err != nil {
			return err
		}
	}

	p := Packet{
		Type:    packetType,
		Message: packetAny,
	}

	message, err := proto.Marshal(&p)
	if err != nil {
		return err
	}

	varint := proto.EncodeVarint(uint64(proto.Size(&p)))

	//buffer := proto.NewBuffer(nil)
	//err := buffer.EncodeMessage(&p)
	//if err != nil {
	//	return err
	//}
	_, err = rw.Writer.Write(append(varint, message...))
	return err
}

func (rw *readWriter) ReadPacket() (Packet, error) {
	b, err := rw.readPacket()
	if err != nil {
		return NilPacket, err
	}

	packet := Packet{}
	err = proto.NewBuffer(b).DecodeMessage(&packet)
	return packet, err
}

func (rw *readWriter) readPacket() ([]byte, error) {
	buf := make([]byte, defaultSize)
	var packetLength uint64
	var bytesRead uint64

	var buffer []byte
	for {
		n, err := rw.Read(buf)
		if err != nil && n == 0 {
			return nil, err
		}
		bytesRead += uint64(n)

		if packetLength == 0 {
			packetLength = readPacketLength(buf)
		}

		if packetLength > bytesRead {
			buffer = append(buffer, buf[:n]...)
			buf = make([]byte, packetLength-bytesRead)
			continue
		}

		buffer = append(buffer, buf[:n]...)

		if bytesRead == packetLength {
			break
		}

		//Two or more packets got buffered

		//fmt.Println(bytesRead, packetLength)
		bufferPacketsLength := bytesRead - packetLength
		var bufferPackets []byte
		bufferPackets = append(bufferPackets, buf[:bufferPacketsLength]...)

		for {
			nextPacket := readPacketLength(bufferPackets)
			var packet []byte
			packet = append(packet, bufferPackets[:nextPacket]...)
			rw.packets = append(rw.packets, packet)

			bufferPacketsLength -= nextPacket

			if bufferPacketsLength == 0 {
				break
			}
		}
		break
	}
	return buffer, nil
}

func (rw *readWriter) Read(buf []byte) (n int, err error) {
	if len(rw.packets) != 0 {
		var packet []byte
		packet, rw.packets = pop(rw.packets)
		copy(buf, packet)
		return len(packet), nil
	}

	return rw.Reader.Read(buf)
}

func readPacketLength(buf []byte) uint64 {
	//fmt.Println(buf)
	//varintLength is not part of the actual message so we need to add it to the real packet length
	n, varintLength := proto.DecodeVarint(buf)
	if varintLength == 0 {
		panic(fmt.Errorf("packet: varint is longer than the real packet"))
	}
	return uint64(varintLength) + n
}
