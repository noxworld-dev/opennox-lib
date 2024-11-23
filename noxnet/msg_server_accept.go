package noxnet

import (
	"encoding/binary"
	"io"

	"github.com/noxworld-dev/opennox-lib/binenc"
)

func init() {
	RegisterServerMessage(&MsgServerAccept{})
	RegisterClientMessage(&MsgClientAccept{})
}

type MsgServerAccept struct {
	Unk0   byte
	Unk1   byte
	ID     uint32
	XorKey byte
}

func (*MsgServerAccept) NetOp() Op {
	return MSG_ACCEPTED
}

func (p *MsgServerAccept) EncodeSize() int {
	return 7
}

func (p *MsgServerAccept) Encode(data []byte) (int, error) {
	if len(data) < 7 {
		return 0, io.ErrShortBuffer
	}
	data[0] = p.Unk0
	data[1] = p.Unk1
	binary.LittleEndian.PutUint32(data[2:6], p.ID)
	data[6] = p.XorKey
	return 7, nil
}

func (p *MsgServerAccept) Decode(data []byte) (int, error) {
	if len(data) < 7 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Unk0 = data[0]
	p.Unk1 = data[1]
	p.ID = binary.LittleEndian.Uint32(data[2:6])
	p.XorKey = data[6]
	return 7, nil
}

type MsgClientAccept struct {
	Unk0         byte   // 0
	Unk1         byte   // 0
	PlayerName   string // 2-67
	PlayerClass  byte   // 68
	IsFemale     byte   // 69
	Unk70        [29]byte
	ScreenWidth  uint32 // 99-102
	ScreenHeight uint32 // 103-106
	Serial       string // 107-128
	Unk129       [26]byte
}

func (*MsgClientAccept) NetOp() Op {
	return MSG_ACCEPTED
}

func (p *MsgClientAccept) EncodeSize() int {
	return 155
}

func (p *MsgClientAccept) Encode(data []byte) (int, error) {
	if len(data) < 155 {
		return 0, io.ErrShortBuffer
	}
	data[0] = p.Unk0
	data[1] = p.Unk1
	binenc.CStringSet16(data[2:68], p.PlayerName)
	data[68] = p.PlayerClass
	data[69] = p.IsFemale
	copy(data[70:99], p.Unk70[:])
	binary.LittleEndian.PutUint32(data[99:103], p.ScreenWidth)
	binary.LittleEndian.PutUint32(data[103:107], p.ScreenHeight)
	binenc.CStringSet(data[107:129], p.Serial)
	copy(data[129:155], p.Unk129[:])
	return 155, nil
}

func (p *MsgClientAccept) Decode(data []byte) (int, error) {
	if len(data) < 155 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Unk0 = data[0]
	p.Unk1 = data[1]
	p.PlayerName = binenc.CString16(data[2:68])
	p.PlayerClass = data[68]
	p.IsFemale = data[69]
	copy(p.Unk70[:], data[70:99])
	p.ScreenWidth = binary.LittleEndian.Uint32(data[99:103])
	p.ScreenHeight = binary.LittleEndian.Uint32(data[103:107])
	p.Serial = binenc.CString(data[107:129])
	copy(p.Unk129[:], data[129:155])
	return 155, nil
}
