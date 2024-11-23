package noxnet

import (
	"encoding/binary"
	"io"
)

func init() {
	RegisterMessage(&MsgImportant{}, false)
	RegisterMessage(&MsgImportantAck{}, false)
}

type MsgImportant struct {
	ID uint32
}

func (*MsgImportant) NetOp() Op {
	return MSG_IMPORTANT
}

func (*MsgImportant) EncodeSize() int {
	return 4
}

func (m *MsgImportant) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:4], m.ID)
	return 4, nil
}

func (m *MsgImportant) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.ID = binary.LittleEndian.Uint32(data[0:4])
	return 4, nil
}

type MsgImportantAck struct {
	ID uint32
}

func (*MsgImportantAck) NetOp() Op {
	return MSG_IMPORTANT_ACK
}

func (*MsgImportantAck) EncodeSize() int {
	return 4
}

func (m *MsgImportantAck) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:4], m.ID)
	return 4, nil
}

func (m *MsgImportantAck) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.ID = binary.LittleEndian.Uint32(data[0:4])
	return 4, nil
}
