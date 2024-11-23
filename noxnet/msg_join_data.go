package noxnet

import (
	"encoding/binary"
	"io"
)

func init() {
	RegisterMessage(&MsgJoinData{}, false)
}

type MsgJoinData struct {
	NetCode uint16
	Unk2    uint32
}

func (*MsgJoinData) NetOp() Op {
	return MSG_JOIN_DATA
}

func (*MsgJoinData) EncodeSize() int {
	return 6
}

func (m *MsgJoinData) Encode(data []byte) (int, error) {
	if len(data) < 6 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint16(data[0:2], m.NetCode)
	binary.LittleEndian.PutUint32(data[2:6], m.Unk2)
	return 6, nil
}

func (m *MsgJoinData) Decode(data []byte) (int, error) {
	if len(data) < 6 {
		return 0, io.ErrUnexpectedEOF
	}
	m.NetCode = binary.LittleEndian.Uint16(data[0:2])
	m.Unk2 = binary.LittleEndian.Uint32(data[2:6])
	return 6, nil
}
