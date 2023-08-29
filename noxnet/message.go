package noxnet

import (
	"fmt"
	"io"
	"reflect"
)

var (
	byOp       = make(map[Op]reflect.Type)
	serverOnly = make(map[Op]reflect.Type)
	clientOnly = make(map[Op]reflect.Type)
)

func RegisterMessage(p Message) {
	op := p.NetOp()
	if _, ok := byOp[op]; ok {
		panic("already registered")
	}
	if _, ok := serverOnly[op]; ok {
		panic("already registered")
	}
	if _, ok := clientOnly[op]; ok {
		panic("already registered")
	}
	byOp[op] = reflect.TypeOf(p).Elem()
}

func RegisterServerMessage(p Message) {
	op := p.NetOp()
	if _, ok := serverOnly[op]; ok {
		panic("already registered")
	}
	if _, ok := byOp[op]; ok {
		panic("already registered")
	}
	serverOnly[op] = reflect.TypeOf(p).Elem()
}

func RegisterClientMessage(p Message) {
	op := p.NetOp()
	if _, ok := clientOnly[op]; ok {
		panic("already registered")
	}
	if _, ok := byOp[op]; ok {
		panic("already registered")
	}
	clientOnly[op] = reflect.TypeOf(p).Elem()
}

type Encoded interface {
	EncodeSize() int
	Encode(data []byte) (int, error)
	Decode(data []byte) (int, error)
}

type Message interface {
	NetOp() Op
	Encoded
}

func EncodePacketSize(p Message) int {
	return 1 + p.EncodeSize()
}

func EncodePacket(data []byte, p Message) (int, error) {
	if sz := EncodePacketSize(p); len(data) < sz {
		return 0, io.ErrShortBuffer
	}
	data[0] = byte(p.NetOp())
	n, err := p.Encode(data[1:])
	if err != nil {
		return 0, err
	}
	return 1 + n, nil
}

func AppendPacket(data []byte, p Message) ([]byte, error) {
	sz := EncodePacketSize(p)
	orig := data
	i := len(orig)
	data = append(data, make([]byte, sz)...)
	buf := data[i : i+sz]
	_, err := EncodePacket(buf, p)
	if err != nil {
		return orig, err
	}
	return data, nil
}

func DecodeAnyPacket(fromServer bool, data []byte) (Message, int, error) {
	if len(data) == 0 {
		return nil, 0, io.EOF
	}
	op := Op(data[0])
	rt, ok := byOp[op]
	if !ok {
		if fromServer {
			rt, ok = serverOnly[op]
		} else {
			rt, ok = clientOnly[op]
		}
	}
	if !ok {
		return nil, 0, fmt.Errorf("unsupported packet: %v", op)
	}
	p := reflect.New(rt).Interface().(Message)
	n, err := p.Decode(data[1:])
	return p, 1 + n, err
}

func DecodePacket(data []byte, p Message) (int, error) {
	if len(data) == 0 {
		return 0, io.EOF
	}
	if got, exp := Op(data[0]), p.NetOp(); got != exp {
		return 0, fmt.Errorf("expected packet: %v, got: %v", exp, got)
	}
	return p.Decode(data[1:])
}
