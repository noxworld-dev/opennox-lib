package noxnet

import (
	"encoding/binary"
	"io"

	"github.com/noxworld-dev/opennox-lib/binenc"
)

func init() {
	RegisterMessage(&MsgServerJoin{})
}

type MsgServerJoin struct {
	Unk0        byte     // 0
	PlayerName  string   // 1-50
	PlayerClass byte     // 51
	PlayerLevel byte     // 52
	Serial      string   // 53-76
	Version     uint32   // 77-80
	Team        uint32   // 81-84
	Unk85       [10]byte // 85-88
	Unk95       byte     // 95
	Unk96       byte     // 96
}

func (*MsgServerJoin) NetOp() Op {
	return MSG_SERVER_JOIN
}

func (*MsgServerJoin) EncodeSize() int {
	return 97
}

func (p *MsgServerJoin) Encode(data []byte) (int, error) {
	if len(data) < 97 {
		return 0, io.ErrShortBuffer
	}
	data[0] = p.Unk0
	binenc.CStringSet16(data[1:51], p.PlayerName)
	binenc.CStringSet(data[53:77], p.Serial)
	binary.LittleEndian.PutUint32(data[77:81], p.Version)
	binary.LittleEndian.PutUint32(data[81:85], p.Team)
	copy(data[85:95], p.Unk85[:])
	data[95] = p.Unk95
	data[96] = p.Unk96
	return 97, nil
}

func (p *MsgServerJoin) Decode(data []byte) (int, error) {
	if len(data) < 97 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Unk0 = data[0]
	p.PlayerName = binenc.CString16(data[1:53])
	p.Serial = binenc.CString(data[53:77])
	p.Version = binary.LittleEndian.Uint32(data[77:81])
	p.Team = binary.LittleEndian.Uint32(data[81:85])
	copy(p.Unk85[:], data[85:95])
	p.Unk95 = data[95]
	p.Unk96 = data[96]
	return 97, nil
}
