package noxnet

import (
	"encoding/binary"
	"io"
)

func init() {
	RegisterMessage(&MsgTimestamp{})
	RegisterMessage(&MsgFullTimestamp{})
	RegisterMessage(&MsgRateChange{})
}

func TimestampSet16(ts uint32, v uint16) uint32 {
	ts16 := uint16(ts & 0xFFFF)
	overflow := (ts16 >= 0xC000) && (v < 0x4000)
	if !overflow && v < ts16 {
		return ts // out of order
	}
	ts = (ts & 0xFFFF0000) | uint32(v)
	if overflow {
		ts += 0x10000
	}
	return ts
}

type MsgTimestamp struct {
	T uint16
}

func (*MsgTimestamp) NetOp() Op {
	return MSG_TIMESTAMP
}

func (*MsgTimestamp) EncodeSize() int {
	return 2
}

func (m *MsgTimestamp) Encode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint16(data[:2], m.T)
	return 2, nil
}

func (m *MsgTimestamp) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	m.T = binary.LittleEndian.Uint16(data[:2])
	return 2, nil
}

type MsgFullTimestamp struct {
	T uint32
}

func (*MsgFullTimestamp) NetOp() Op {
	return MSG_FULL_TIMESTAMP
}

func (*MsgFullTimestamp) EncodeSize() int {
	return 4
}

func (m *MsgFullTimestamp) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[:4], m.T)
	return 4, nil
}

func (m *MsgFullTimestamp) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.T = binary.LittleEndian.Uint32(data[:4])
	return 4, nil
}

type MsgRateChange struct {
	Rate byte
}

func (*MsgRateChange) NetOp() Op {
	return MSG_RATE_CHANGE
}

func (*MsgRateChange) EncodeSize() int {
	return 1
}

func (p *MsgRateChange) Encode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrShortBuffer
	}
	data[0] = p.Rate
	return 1, nil
}

func (p *MsgRateChange) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Rate = data[0]
	return 1, nil
}
