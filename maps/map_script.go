package maps

import (
	"encoding/binary"
	"fmt"
	"io"
)

func init() {
	RegisterSection(&Script{})
	RegisterSection(&ScriptData{})
}

type Script struct {
	Data []byte
}

func (*Script) MapSection() string {
	return "ScriptObject"
}

func (sect *Script) MarshalBinary() ([]byte, error) {
	data := make([]byte, 2+4+len(sect.Data))
	binary.LittleEndian.PutUint16(data[0:], 1) // version
	binary.LittleEndian.PutUint32(data[2:], uint32(len(sect.Data)))
	copy(data[6:], sect.Data)
	return data, nil
}

func (sect *Script) UnmarshalBinary(data []byte) error {
	if len(data) < 2 {
		return io.ErrUnexpectedEOF
	}
	vers := binary.LittleEndian.Uint16(data[0:])
	if vers != 1 {
		return fmt.Errorf("unsupported version of script section: %d", vers)
	}
	if len(data) < 6 {
		return io.ErrUnexpectedEOF
	}
	size := binary.LittleEndian.Uint32(data[2:])
	data = data[6:]
	if size > uint32(len(data)) {
		return io.ErrUnexpectedEOF
	}
	sect.Data = make([]byte, size)
	copy(sect.Data, data)
	return nil
}

type ScriptData struct {
	Data []byte
}

func (*ScriptData) MapSection() string {
	return "ScriptData"
}

func (sect *ScriptData) MarshalBinary() ([]byte, error) {
	data := make([]byte, 2+1+len(sect.Data))
	binary.LittleEndian.PutUint16(data[0:], 1) // version
	if len(sect.Data) != 0 {
		data[2] = 1
	}
	copy(data[3:], sect.Data)
	return data, nil
}

func (sect *ScriptData) UnmarshalBinary(data []byte) error {
	if len(data) < 2 {
		return io.ErrUnexpectedEOF
	}
	vers := binary.LittleEndian.Uint16(data[0:])
	if vers != 1 {
		return fmt.Errorf("unsupported version of script data section: %d", vers)
	}
	if len(data) < 3 {
		return io.ErrUnexpectedEOF
	}
	if data[2] == 0 {
		sect.Data = nil
		return nil
	}
	data = data[3:]
	sect.Data = make([]byte, len(data))
	copy(sect.Data, data)
	return nil
}
