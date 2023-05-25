package maps

import (
	"encoding/binary"
	"fmt"
	"image"
	"io"
)

func init() {
	RegisterSection(&SecretWalls{})
}

type SecretWall struct {
	Pos       image.Point
	OpenWait  uint32
	Flags     byte
	State     byte
	OpenDelay byte
	LastOpen  uint32
	R2        uint32
}

func (w *SecretWall) MarshalBinary() ([]byte, error) {
	data := make([]byte, 23)
	binary.LittleEndian.PutUint32(data[0:], uint32(int32(w.Pos.X)))
	binary.LittleEndian.PutUint32(data[4:], uint32(int32(w.Pos.Y)))
	binary.LittleEndian.PutUint32(data[8:], w.OpenWait)
	data[12] = w.Flags
	data[13] = w.State
	data[14] = w.OpenDelay
	binary.LittleEndian.PutUint32(data[15:], w.LastOpen)
	binary.LittleEndian.PutUint32(data[19:], w.R2)
	return data, nil
}

func (w *SecretWall) UnmarshalBinary(data []byte) error {
	if len(data) < 23 {
		return io.ErrUnexpectedEOF
	}
	w.Pos.X = int(binary.LittleEndian.Uint32(data[0:]))
	w.Pos.Y = int(binary.LittleEndian.Uint32(data[4:]))
	w.OpenWait = binary.LittleEndian.Uint32(data[8:])
	w.Flags = data[12]
	w.State = data[13]
	w.OpenDelay = data[14]
	w.LastOpen = binary.LittleEndian.Uint32(data[15:])
	w.R2 = binary.LittleEndian.Uint32(data[19:])
	return nil
}

type SecretWalls struct {
	Walls []SecretWall
}

func (*SecretWalls) MapSection() string {
	return "SecretWalls"
}

func (sect *SecretWalls) MarshalBinary() ([]byte, error) {
	data := make([]byte, 4, 4+23*len(sect.Walls))
	binary.LittleEndian.PutUint16(data[0:], 2)
	binary.LittleEndian.PutUint16(data[2:], uint16(len(sect.Walls)))
	for _, w := range sect.Walls {
		b, err := w.MarshalBinary()
		if err != nil {
			return nil, err
		}
		data = append(data, b...)
	}
	return data, nil
}

func (sect *SecretWalls) UnmarshalBinary(data []byte) error {
	*sect = SecretWalls{}
	if len(data) < 2 {
		return io.ErrUnexpectedEOF
	}
	vers := binary.LittleEndian.Uint16(data)
	data = data[2:]
	if vers != 2 {
		return fmt.Errorf("unsupported version of secret walls section: %d", vers)
	}
	if len(data) < 2 {
		return io.ErrUnexpectedEOF
	}
	n := int(binary.LittleEndian.Uint16(data))
	data = data[2:]
	sect.Walls = make([]SecretWall, 0, n)
	for i := 0; i < n; i++ {
		var w SecretWall
		if err := w.UnmarshalBinary(data); err != nil {
			return err
		}
		data = data[23:]
		sect.Walls = append(sect.Walls, w)
	}
	return nil
}
