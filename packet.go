package network

import (
	"bytes"
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

var (
	defaultSize = 4096
)

func writePacket(writer io.Writer, uuid string, packetType uint64, packet proto.Message) error {
	var packetAny *any.Any
	if packet != nil {
		var err error
		packetAny, err = ptypes.MarshalAny(packet)
		if err != nil {
			return err
		}
	}

	p := &Packet{
		Type:    packetType,
		UUID:    uuid,
		Message: packetAny,
	}

	buffer := proto.NewBuffer(nil)
	err := buffer.EncodeMessage(p)
	if err != nil {
		return err
	}
	_, err = writer.Write(buffer.Bytes())
	return err
}

func read(reader io.Reader) ([]byte, error) {
	buf := make([]byte, defaultSize)
	var packetLength int
	var bytesRead int
	var buffer bytes.Buffer
	for {
		n, err := reader.Read(buf)
		if err != nil && n == 0 {
			return nil, err
		}
		bytesRead += n

		if packetLength == 0 {
			//varintLength is not part of the actual message so we need to add it to the real packet length
			n, varintLength := proto.DecodeVarint(buf)
			if varintLength == 0 {
				return nil, fmt.Errorf("packet: varint var longer than the real packet")
			}
			packetLength = varintLength + int(n)
		}

		_, err = buffer.Write(buf[:n])
		if err != nil {
			return nil, err
		}

		if bytesRead == packetLength {
			break
		}
		buf = make([]byte, packetLength-bytesRead)
	}
	return buffer.Bytes(), nil
}

func readPacket(reader io.Reader) (*Packet, error) {
	b, err := read(reader)
	if err != nil {
		return nil, err
	}
	buffer := proto.NewBuffer(b)

	p := &Packet{}
	err = buffer.DecodeMessage(p)
	return p, err
}

/**
THIS BELOW IS FASTER THAN THIS ABOVE: SAVE FOR LEGACY
**/
//
//func oldreadPacket(conn net.Conn) (*pb.Packet, error) {
//	buffer, err := oldread(conn)
//	if err != nil {
//		return nil, err
//	}
//	//[8 10 18 6 97 115 100 97 115 100 26 22 10 20 116 121 112 101 46 103 111 111 103 108 101 97 112 105 115 46 99 111 109 47]
//
//	fmt.Println(buffer)
//
//	p := &pb.Packet{}
//	err = proto.Unmarshal(buffer, p)
//	if err != nil {
//		return nil, err
//	}
//
//	return p, nil
//}
//
////Validate packetType above 10
//func oldwritePacket(conn net.Conn, uid string, packetType uint64, packet proto.Message) (int, error) {
//	any, err := ptypes.MarshalAny(packet)
//	if err != nil {
//		return 0, fmt.Errorf("failed to marshal any: %s", err)
//	}
//
//	p := &pb.Packet{
//		Type:    packetType,
//		UUID:    uid,
//		Message: any,
//	}
//
//	b, err := proto.Marshal(p)
//	if err != nil {
//		return 0, err
//	}
//	return conn.Write(createPacket(b))
//}
//
//func createPacket(b []byte) []byte {
//	buf := make([]byte, packetLengthByte)
//	binary.LittleEndian.PutUint64(buf, uint64(len(b)))
//
//	return append(buf, b...)
//}
//
////Benchmark append vs buffer.write
////try to use bufio reader
//func oldread(conn net.Conn) ([]byte, error) {
//	buffer := make([]byte, defaultSize)
//	var buf []byte
//	var packetLength int
//	var readBytes int
//
//	for {
//		read, err := conn.Read(buffer)
//		if err != nil {
//			return nil, err
//		}
//
//		if packetLength == 0 {
//			packetLength = int(binary.LittleEndian.Uint64(buffer[0:packetLengthByte]))
//			buffer = buffer[packetLengthByte:]
//			read -= packetLengthByte
//		}
//		readBytes += read
//
//		//fmt.Printf("read %d \n", read)
//		//fmt.Printf("readBytes %d = %d packetLength\n", readBytes, packetLength)
//		if readBytes == packetLength {
//			return append(buf, buffer[:read]...), nil
//		}
//		buf = buffer
//		buffer = make([]byte, packetLength)
//	}
//}
