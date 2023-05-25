package maps

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/noxworld-dev/opennox-lib/types"
)

func init() {
	RegisterSection(&Waypoints{})
}

type WaypointLink struct {
	ID    uint32
	Flags byte
}

func (w *WaypointLink) MarshalBinary() ([]byte, error) {
	data := make([]byte, 5)
	binary.LittleEndian.PutUint32(data[0:], w.ID)
	data[4] = w.Flags
	return data, nil
}

func (w *WaypointLink) UnmarshalBinary(data []byte) error {
	if len(data) < 5 {
		return io.ErrUnexpectedEOF
	}
	w.ID = binary.LittleEndian.Uint32(data[0:])
	w.Flags = data[4]
	return nil
}

type Waypoint struct {
	ID    uint32
	Pos   types.Pointf
	Name  string
	Flags uint32
	Links []WaypointLink
}

func (w *Waypoint) EncodingSize() int {
	return 13 + len(w.Name) + 5 + 5*len(w.Links)
}

func (w *Waypoint) MarshalBinary() ([]byte, error) {
	data := make([]byte, 13+len(w.Name)+5, w.EncodingSize())
	i := 0
	binary.LittleEndian.PutUint32(data[i:], w.ID)
	i += 4
	binary.LittleEndian.PutUint32(data[i:], math.Float32bits(w.Pos.X))
	i += 4
	binary.LittleEndian.PutUint32(data[i:], math.Float32bits(w.Pos.Y))
	i += 4
	data[i] = byte(len(w.Name))
	i++
	i += copy(data[13:], w.Name)
	binary.LittleEndian.PutUint32(data[i:], w.Flags)
	i += 4
	data[i] = byte(len(w.Links))
	for _, l := range w.Links {
		b, err := l.MarshalBinary()
		if err != nil {
			return nil, err
		}
		data = append(data, b...)
	}
	return data, nil
}

func (w *Waypoint) UnmarshalBinary(data []byte) error {
	if len(data) < 13 {
		return io.ErrUnexpectedEOF
	}
	w.ID = binary.LittleEndian.Uint32(data[0:])
	data = data[4:]
	w.Pos.X = math.Float32frombits(binary.LittleEndian.Uint32(data[0:]))
	w.Pos.Y = math.Float32frombits(binary.LittleEndian.Uint32(data[4:]))
	data = data[8:]
	sz := int(data[0])
	data = data[1:]
	if len(data) < sz {
		return io.ErrUnexpectedEOF
	}
	w.Name = string(data[:sz])
	data = data[sz:]
	if len(data) < 5 {
		return io.ErrUnexpectedEOF
	}
	w.Flags = binary.LittleEndian.Uint32(data[0:])
	data = data[4:]
	sz = int(data[0])
	data = data[1:]
	for i := 0; i < sz; i++ {
		var l WaypointLink
		if err := l.UnmarshalBinary(data); err != nil {
			return err
		}
		data = data[5:]
		w.Links = append(w.Links, l)
	}
	return nil
}

type Waypoints struct {
	Waypoints []Waypoint
}

func (*Waypoints) MapSection() string {
	return "WayPoints"
}

func (sect *Waypoints) MarshalBinary() ([]byte, error) {
	data := make([]byte, 6)
	binary.LittleEndian.PutUint16(data[0:], 4)
	binary.LittleEndian.PutUint32(data[2:], uint32(len(sect.Waypoints)))
	for _, w := range sect.Waypoints {
		b, err := w.MarshalBinary()
		if err != nil {
			return nil, err
		}
		data = append(data, b...)
	}
	return data, nil
}

func (sect *Waypoints) UnmarshalBinary(data []byte) error {
	*sect = Waypoints{}
	if len(data) < 2 {
		return io.ErrUnexpectedEOF
	}
	vers := binary.LittleEndian.Uint16(data)
	data = data[2:]
	if vers != 4 {
		return fmt.Errorf("unsupported version of waypoints section: %d", vers)
	}
	if len(data) < 4 {
		return io.ErrUnexpectedEOF
	}
	n := int(binary.LittleEndian.Uint32(data))
	data = data[4:]
	sect.Waypoints = make([]Waypoint, 0, n)
	for i := 0; i < n; i++ {
		var w Waypoint
		if err := w.UnmarshalBinary(data); err != nil {
			return err
		}
		sect.Waypoints = append(sect.Waypoints, w)
		data = data[w.EncodingSize():]
	}
	return nil
}
