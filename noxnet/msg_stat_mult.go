package noxnet

import (
	"encoding/binary"
	"io"
	"math"
)

func init() {
	RegisterMessage(&MsgStatMult{}, false)
}

type MsgStatMult struct {
	Health   float32
	Mana     float32
	Speed    float32
	Strength float32
}

func (*MsgStatMult) NetOp() Op {
	return MSG_STAT_MULTIPLIERS
}

func (*MsgStatMult) EncodeSize() int {
	return 16
}

func (p *MsgStatMult) Encode(data []byte) (int, error) {
	if len(data) < 16 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:], math.Float32bits(p.Health))
	binary.LittleEndian.PutUint32(data[4:], math.Float32bits(p.Mana))
	binary.LittleEndian.PutUint32(data[8:], math.Float32bits(p.Strength))
	binary.LittleEndian.PutUint32(data[12:], math.Float32bits(p.Speed))
	return 16, nil
}

func (p *MsgStatMult) Decode(data []byte) (int, error) {
	if len(data) < 16 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Health = math.Float32frombits(binary.LittleEndian.Uint32(data[0:]))
	p.Mana = math.Float32frombits(binary.LittleEndian.Uint32(data[4:]))
	p.Strength = math.Float32frombits(binary.LittleEndian.Uint32(data[8:]))
	p.Speed = math.Float32frombits(binary.LittleEndian.Uint32(data[12:]))
	return 16, nil
}
