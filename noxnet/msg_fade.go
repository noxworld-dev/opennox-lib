package noxnet

import (
	"io"
)

func init() {
	RegisterMessage(&MsgFadeBegin{}, false)
}

type MsgFadeBegin struct {
	Out  byte
	Menu byte
}

func (*MsgFadeBegin) NetOp() Op {
	return MSG_FADE_BEGIN
}

func (*MsgFadeBegin) EncodeSize() int {
	return 2
}

func (p *MsgFadeBegin) Encode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrShortBuffer
	}
	data[0] = p.Out
	data[1] = p.Menu
	return 2, nil
}

func (p *MsgFadeBegin) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Out = data[0]
	p.Menu = data[1]
	return 2, nil
}
