package maps

import (
	"encoding/binary"
	"fmt"
	"image"
	"io"
)

func init() {
	RegisterSection(&DestructableWalls{})
}

type DestructableWall struct {
	Pos image.Point
}

func (w *DestructableWall) MarshalBinary() ([]byte, error) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint32(data[0:], uint32(int32(w.Pos.X)))
	binary.LittleEndian.PutUint32(data[4:], uint32(int32(w.Pos.Y)))
	return data, nil
}

func (w *DestructableWall) UnmarshalBinary(data []byte) error {
	if len(data) < 8 {
		return io.ErrUnexpectedEOF
	}
	w.Pos.X = int(binary.LittleEndian.Uint32(data[0:]))
	w.Pos.Y = int(binary.LittleEndian.Uint32(data[4:]))
	return nil
}

type DestructableWalls struct {
	Walls []DestructableWall
}

func (*DestructableWalls) MapSection() string {
	return "DestructableWalls"
}

func (sect *DestructableWalls) MarshalBinary() ([]byte, error) {
	data := make([]byte, 4, 4+8*len(sect.Walls))
	binary.LittleEndian.PutUint16(data[0:], 1)
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

func (sect *DestructableWalls) UnmarshalBinary(data []byte) error {
	*sect = DestructableWalls{}
	if len(data) < 2 {
		return io.ErrUnexpectedEOF
	}
	vers := binary.LittleEndian.Uint16(data)
	data = data[2:]
	if vers != 1 {
		return fmt.Errorf("unsupported version of destructable walls section: %d", vers)
	}
	if len(data) < 2 {
		return io.ErrUnexpectedEOF
	}
	n := int(binary.LittleEndian.Uint16(data))
	data = data[2:]
	sect.Walls = make([]DestructableWall, 0, n)
	for i := 0; i < n; i++ {
		var w DestructableWall
		if err := w.UnmarshalBinary(data); err != nil {
			return err
		}
		data = data[8:]
		sect.Walls = append(sect.Walls, w)
	}
	return nil
}
