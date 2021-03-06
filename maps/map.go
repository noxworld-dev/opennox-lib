package maps

import (
	"encoding"
	"errors"
	"image"
	"reflect"
)

var mapSections = make(map[string]reflect.Type)

func RegisterSection(sect Section) {
	name := sect.MapSection()
	if _, ok := mapSections[name]; ok {
		panic("already registered: " + name)
	}
	mapSections[name] = reflect.TypeOf(sect).Elem()
}

type RawSection struct {
	Name string
	Data []byte
}

func (sect RawSection) Supported() bool {
	_, ok := mapSections[sect.Name]
	return ok
}

func (sect RawSection) Decode() (Section, error) {
	typ, ok := mapSections[sect.Name]
	if !ok {
		return nil, errors.New("unsupported section: " + sect.Name)
	}
	s := reflect.New(typ).Interface().(Section)
	err := s.UnmarshalBinary(sect.Data)
	return s, err
}

type Section interface {
	MapSection() string
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type Info struct {
	Filename string `json:"name"`
	Size     int    `json:"size"`
	MapInfo
}

type Map struct {
	Info

	crc      uint32
	wallOffX uint32
	wallOffY uint32

	Intro      *MapIntro
	Ambient    *AmbientData
	Walls      *WallMap
	Floor      *FloorMap
	Script     *Script
	ScriptData *ScriptData
	Unknown    []RawSection
}

// GridBoundingBox returns a bounding box for all walls and tiles on the map.
// Returned rectangle uses grid coordinates, not pixel coordinates.
func (m *Map) GridBoundingBox() image.Rectangle {
	var r image.Rectangle
	r.Min.X = -1
	r.Min.Y = -1
	for _, w := range m.Walls.Walls {
		p := w.Pos
		if r.Min.X == -1 || r.Min.X > int(p.X) {
			r.Min.X = int(p.X)
		}
		if r.Min.Y == -1 || r.Min.Y > int(p.Y) {
			r.Min.Y = int(p.Y)
		}
		if r.Max.X < int(p.X) {
			r.Max.X = int(p.X)
		}
		if r.Max.Y < int(p.Y) {
			r.Max.Y = int(p.Y)
		}
	}
	// TODO: tiles
	return r
}
