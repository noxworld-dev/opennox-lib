package noxnet

import (
	"github.com/noxworld-dev/opennox-lib/noxnet/xfer"
)

func init() {
	RegisterMessage(&MsgXfer{})
}

type MsgXfer struct {
	Msg xfer.Msg
}

func (*MsgXfer) NetOp() Op {
	return MSG_XFER_MSG
}

func (m *MsgXfer) EncodeSize() int {
	return xfer.EncodeSize(m.Msg)
}

func (m *MsgXfer) Encode(data []byte) (int, error) {
	return xfer.Encode(data, m.Msg)
}

func (m *MsgXfer) Decode(data []byte) (int, error) {
	msg, n, err := xfer.Decode(data)
	m.Msg = msg
	return n, err
}
