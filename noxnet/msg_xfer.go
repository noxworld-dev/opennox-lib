package noxnet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

func init() {
	RegisterMessage(&MsgXfer{})
}

type XferCode byte

const (
	XferStart  = XferCode(0)
	XferAccept = XferCode(1)
	XferData   = XferCode(2)
	XferAck    = XferCode(3)
	XferClose  = XferCode(4)
	XferCode5  = XferCode(5)
	XferCode6  = XferCode(6)
)

type SubXfer interface {
	XferOp() XferCode
	EncodeSize() int
	Encode(data []byte) (int, error)
	Decode(data []byte) (int, error)
}

type MsgXferStart struct {
	Act   byte
	Unk1  byte
	Size  uint32
	Type  FixedString
	Token byte
	Unk5  [3]byte
}

func (*MsgXferStart) XferOp() XferCode {
	return XferStart
}

func (*MsgXferStart) EncodeSize() int {
	return 138
}

func (m *MsgXferStart) Encode(data []byte) (int, error) {
	if len(data) < 138 {
		return 0, io.ErrShortBuffer
	}
	data[0] = m.Act
	data[1] = m.Unk1
	binary.LittleEndian.PutUint32(data[2:6], m.Size)
	m.Type.Encode(data[6:134])
	data[134] = m.Token
	copy(data[135:138], m.Unk5[:])
	return 138, nil
}

func (m *MsgXferStart) Decode(data []byte) (int, error) {
	if len(data) < 138 {
		return 0, io.ErrUnexpectedEOF
	}
	m.Act = data[0]
	m.Unk1 = data[1]
	m.Size = binary.LittleEndian.Uint32(data[2:6])
	m.Type.Decode(data[6:134])
	m.Token = data[134]
	copy(m.Unk5[:], data[135:138])
	return 138, nil
}

type MsgXferState struct {
	Code   XferCode
	Stream byte
	Token  byte
}

func (m *MsgXferState) XferOp() XferCode {
	return m.Code
}

func (*MsgXferState) EncodeSize() int {
	return 2
}

func (m *MsgXferState) Encode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrShortBuffer
	}
	data[0] = m.Stream
	data[1] = m.Token
	return 2, nil
}

func (m *MsgXferState) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	m.Stream = data[0]
	m.Token = data[1]
	return 2, nil
}

type MsgXferData struct {
	Stream byte
	Token  byte
	Chunk  uint16
	Data   []byte
}

func (*MsgXferData) XferOp() XferCode {
	return XferData
}

func (m *MsgXferData) EncodeSize() int {
	return 6 + len(m.Data)
}

func (m *MsgXferData) Encode(data []byte) (int, error) {
	if len(data) < 6+len(m.Data) {
		return 0, io.ErrShortBuffer
	}
	if len(m.Data) > math.MaxUint16 {
		return 0, errors.New("xfer packet too large")
	}
	data[0] = m.Stream
	data[1] = m.Token
	binary.LittleEndian.PutUint16(data[2:4], m.Chunk)
	binary.LittleEndian.PutUint16(data[4:6], uint16(len(m.Data)))
	copy(data[6:], m.Data)
	return 6 + len(m.Data), nil
}

func (m *MsgXferData) Decode(data []byte) (int, error) {
	if len(data) < 6 {
		return 0, io.ErrUnexpectedEOF
	}
	m.Stream = data[0]
	m.Token = data[1]
	m.Chunk = binary.LittleEndian.Uint16(data[2:4])
	sz := int(binary.LittleEndian.Uint16(data[4:6]))
	data = data[6:]
	if sz > len(data) {
		return 0, io.ErrUnexpectedEOF
	}
	m.Data = make([]byte, sz)
	copy(m.Data, data[:sz])
	return 6 + sz, nil
}

type MsgXferAck struct {
	Stream byte
	Token  byte
	Chunk  uint16
}

func (*MsgXferAck) XferOp() XferCode {
	return XferAck
}

func (*MsgXferAck) EncodeSize() int {
	return 4
}

func (m *MsgXferAck) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	data[0] = m.Stream
	data[1] = m.Token
	binary.LittleEndian.PutUint16(data[2:4], m.Chunk)
	return 4, nil
}

func (m *MsgXferAck) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.Stream = data[0]
	m.Token = data[1]
	m.Chunk = binary.LittleEndian.Uint16(data[2:4])
	return 4, nil
}

type MsgXferClose struct {
	Stream byte
}

func (*MsgXferClose) XferOp() XferCode {
	return XferClose
}

func (*MsgXferClose) EncodeSize() int {
	return 1
}

func (m *MsgXferClose) Encode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrShortBuffer
	}
	data[0] = m.Stream
	return 1, nil
}

func (m *MsgXferClose) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	m.Stream = data[0]
	return 1, nil
}

type MsgXfer struct {
	Msg SubXfer
}

func (*MsgXfer) NetOp() Op {
	return MSG_XFER_MSG
}

func (m *MsgXfer) EncodeSize() int {
	return 1 + m.Msg.EncodeSize()
}

func (m *MsgXfer) Encode(data []byte) (int, error) {
	if len(data) < m.EncodeSize() {
		return 0, io.ErrShortBuffer
	}
	data[0] = byte(m.Msg.XferOp())
	n, err := m.Msg.Encode(data[1:])
	if err != nil {
		return 0, err
	}
	return 1 + n, nil
}

func (m *MsgXfer) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	typ := XferCode(data[0])
	switch typ {
	case XferStart:
		m.Msg = &MsgXferStart{}
	case XferAccept, XferCode5, XferCode6:
		m.Msg = &MsgXferState{Code: typ}
	case XferData:
		m.Msg = &MsgXferData{}
	case XferAck:
		m.Msg = &MsgXferAck{}
	case XferClose:
		m.Msg = &MsgXferClose{}
	default:
		return 0, fmt.Errorf("unexpected xfer message subtype: 0x%x", typ)
	}
	n, err := m.Msg.Decode(data[1:])
	if err != nil {
		return 0, err
	}
	return 1 + n, nil
}
