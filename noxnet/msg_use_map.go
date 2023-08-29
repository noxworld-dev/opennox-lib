package noxnet

import (
	"encoding/binary"
	"io"
)

func init() {
	RegisterMessage(&MsgUseMap{})
}

type MsgUseMap struct {
	MapName     string
	CRC         uint32
	T           uint32
	MapNameJunk []byte
}

func (*MsgUseMap) NetOp() Op {
	return MSG_USE_MAP
}

func (*MsgUseMap) EncodeSize() int {
	return 40
}

func (m *MsgUseMap) Encode(data []byte) (int, error) {
	if len(data) < 40 {
		return 0, io.ErrShortBuffer
	}
	i := cStringSet(data[0:32], m.MapName)
	if len(m.MapNameJunk) != 0 {
		copy(data[i+1:], m.MapNameJunk)
	}
	binary.LittleEndian.PutUint32(data[32:36], m.CRC)
	binary.LittleEndian.PutUint32(data[36:40], m.T)
	return 40, nil
}

func (m *MsgUseMap) Decode(data []byte) (int, error) {
	if len(data) < 40 {
		return 0, io.ErrUnexpectedEOF
	}
	m.MapName = cString(data[0:32])
	m.MapNameJunk = nil
	if i := len(m.MapName); !allZeros(data[i+1 : 32]) {
		m.MapNameJunk = make([]byte, 32-i-1)
		copy(m.MapNameJunk, data[i+1:32])
	}
	m.CRC = binary.LittleEndian.Uint32(data[32:36])
	m.T = binary.LittleEndian.Uint32(data[36:40])
	return 40, nil
}
