package noxnet

import (
	"io"

	"github.com/noxworld-dev/opennox-lib/binenc"
)

func init() {
	RegisterMessage(&MsgServerPass{}, false)
}

type MsgServerPass struct {
	Unk0 byte   // 0
	Pass string // 1-19
}

func (*MsgServerPass) NetOp() Op {
	return MSG_SERVER_PASSWORD
}

func (*MsgServerPass) EncodeSize() int {
	return 19
}

func (p *MsgServerPass) Encode(data []byte) (int, error) {
	if len(data) < 19 {
		return 0, io.ErrShortBuffer
	}
	data[0] = p.Unk0
	binenc.CStringSet16(data[1:19], p.Pass)
	return 19, nil
}

func (p *MsgServerPass) Decode(data []byte) (int, error) {
	if len(data) < 19 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Unk0 = data[0]
	p.Pass = binenc.CString16(data[1:19])
	return 19, nil
}
