package noxnet

import (
	"fmt"
	"io"
)

func init() {
	RegisterMessage(&MsgServerError{})
}

type ConnectError byte

func (e ConnectError) Name() string {
	switch e {
	case ErrLowPing:
		return "ErrLowPing"
	case ErrHighPing:
		return "ErrHighPing"
	case ErrLowLevel:
		return "ErrLowLevel"
	case ErrHighLevel:
		return "ErrHighLevel"
	case ErrClosed:
		return "ErrClosed"
	case ErrBanned:
		return "ErrBanned"
	case ErrWrongPassword:
		return "ErrWrongPassword"
	case ErrIllegalClass:
		return "ErrIllegalClass"
	case ErrTimeOut:
		return "ErrTimeOut"
	case ErrFindFailed:
		return "ErrFindFailed"
	case ErrNeedRefresh:
		return "ErrNeedRefresh"
	case ErrFull:
		return "ErrFull"
	case ErrDupSerial:
		return "ErrDupSerial"
	case ErrWrongVer:
		return "ErrWrongVer"
	}
	return fmt.Sprintf("ConnectError(%d)", int(e))
}

const (
	ErrLowPing = ConnectError(iota)
	ErrHighPing
	ErrLowLevel
	ErrHighLevel
	ErrClosed
	ErrBanned
	ErrWrongPassword
	ErrIllegalClass
	ErrTimeOut
	ErrFindFailed
	ErrNeedRefresh
	ErrFull
	ErrDupSerial
	ErrWrongVer
)

type MsgServerError struct {
	Err ConnectError
}

func (*MsgServerError) NetOp() Op {
	return MSG_SERVER_ERROR
}

func (*MsgServerError) EncodeSize() int {
	return 1
}

func (p *MsgServerError) Encode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrShortBuffer
	}
	data[0] = byte(p.Err)
	return 1, nil
}

func (p *MsgServerError) Decode(data []byte) (int, error) {
	if len(data) < 9 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Err = ConnectError(data[0])
	return 1, nil
}
