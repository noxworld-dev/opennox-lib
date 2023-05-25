package maps

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/noxworld-dev/opennox-lib/binenc"
)

func init() {
	RegisterSection(&Script{})
	RegisterSection(&ScriptData{})
}

func ReadScript(r io.Reader) (*Script, error) {
	d, err := NewReader(r)
	if err != nil {
		return nil, err
	}
	return d.ReadScript()
}

func (r *Reader) ReadScript() (*Script, error) {
	data, err := r.readScriptRaw()
	if err != nil {
		return nil, err
	} else if len(data) == 0 {
		return nil, nil
	}
	r.m.Script = new(Script)
	err = r.m.Script.UnmarshalBinary(data)
	return r.m.Script, err
}

func (r *Reader) readScriptRaw() ([]byte, error) {
	for {
		sect, err := r.nextSection()
		if err == io.EOF {
			return nil, nil
		} else if err != nil {
			return nil, err
		}
		if sect == "ScriptObject" {
			return io.ReadAll(r.r)
		}
	}
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
	return sect.Decode(binenc.NewReader(data))
}

func (sect *Script) Decode(r *binenc.Reader) error {
	vers, ok := r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if vers != 1 {
		return fmt.Errorf("unsupported version of script section: %d", vers)
	}
	sect.Data, ok = r.ReadBytes32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
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
	return sect.Decode(binenc.NewReader(data))
}

func (sect *ScriptData) Decode(r *binenc.Reader) error {
	vers, ok := r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if vers != 1 {
		return fmt.Errorf("unsupported version of script data section: %d", vers)
	}
	has, ok := r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if has == 0 {
		sect.Data = nil
		return nil
	}
	sect.Data = r.ReadAllBytes()
	return nil
}
