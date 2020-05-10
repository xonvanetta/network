package packet

import (
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

type writer struct {
	writer io.Writer
}

type Writer interface {
	Write(packetType uint64, packet proto.Message) error
}

func NewWriter(w io.Writer) Writer {
	return &writer{writer: w}
}

func (w writer) Write(packetType uint64, packet proto.Message) error {
	var packetAny *any.Any
	if packet != nil {
		var err error
		packetAny, err = ptypes.MarshalAny(packet)
		if err != nil {
			return err
		}
	}

	//todo: use pool to buffer them in memory
	p := Packet{
		Type:    packetType,
		Message: packetAny,
	}

	//todo: use pool to buffer them in memory
	buffer := new(proto.Buffer)
	err := buffer.EncodeFixed64(uint64(proto.Size(&p)))
	if err != nil {
		return err
	}

	err = buffer.Marshal(&p)
	if err != nil {
		return err
	}

	//todo: buffer.Reset() use pool for mem usage

	_, err = w.writer.Write(buffer.Bytes())
	return err
}
