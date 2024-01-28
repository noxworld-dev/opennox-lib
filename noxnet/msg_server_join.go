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
	Unk0       byte     // 0
	PlayerName string   // 1-52?
	Serial     string   // 53-74
	Unk75      [2]byte  // 75-76
	Version    uint32   // 77-80
	Unk81      [16]byte // 81-96
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
	binenc.CStringSet16(data[1:53], p.PlayerName)
	binenc.CStringSet(data[53:75], p.Serial)
	copy(data[75:77], p.Unk75[:])
	binary.LittleEndian.PutUint32(data[77:81], p.Version)
	copy(data[81:97], p.Unk81[:])
	return 97, nil
}

func (p *MsgServerJoin) Decode(data []byte) (int, error) {
	if len(data) < 97 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Unk0 = data[0]
	p.PlayerName = binenc.CString16(data[1:53])
	p.Serial = binenc.CString(data[53:75])
	copy(p.Unk75[:], data[75:77])
	p.Version = binary.LittleEndian.Uint32(data[77:81])
	copy(p.Unk81[:], data[81:97])
	return 97, nil
}
