package maps

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func init() {
	RegisterSection(&MapInfo{})
	RegisterSection(&MapIntro{})
	RegisterSection(&AmbientData{})
}

type MapInfo struct {
	Format        int           `json:"format,omitempty"`
	Summary       string        `json:"summary,omitempty"`        // 0 [64]
	Description   string        `json:"description,omitempty"`    // 64 [512]
	Version       string        `json:"version,omitempty"`        // 576 [16]
	Author        string        `json:"author,omitempty"`         // 592 [64]
	Email         string        `json:"email,omitempty"`          // 656 [64]
	Author2       string        `json:"author_2,omitempty"`       // 720 [128]
	Email2        string        `json:"email_2,omitempty"`        // 848 [128]
	Field7        string        `json:",omitempty"`               // 976 [256]
	Copyright     string        `json:"copyright,omitempty"`      // 1232 [128]
	Date          string        `json:"date_str,omitempty"`       // 1360 [32]
	Flags         uint32        `json:"flags,omitempty"`          // 1392
	MinPlayers    byte          `json:"min_players,omitempty"`    // 1396
	MaxPlayers    byte          `json:"max_players,omitempty"`    // 1397
	QuestIntro    string        `json:"quest_intro,omitempty"`    // 1398
	QuestGraphics string        `json:"quest_graphics,omitempty"` // 1430
	Trailing      MapInfoCompat `json:"trailing,omitempty"`
}

type MapInfoCompat struct {
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
	Author      string `json:"author,omitempty"`
	Email       string `json:"email,omitempty"`
	Author2     string `json:"author_2,omitempty"`
	Email2      string `json:"email_2,omitempty"`
	Field7      string `json:",omitempty"`
	Copyright   string `json:"copyright,omitempty"`
	Date        string `json:"date_str,omitempty"`
}

func (*MapInfo) MapSection() string {
	return "MapInfo"
}

func (info *MapInfo) MarshalBinary() ([]byte, error) {
	if info.Format > 3 {
		return nil, fmt.Errorf("unsupported version: %d", info.Format)
	}
	if info.Format < 1 {
		data := make([]byte, 2)
		binary.LittleEndian.PutUint16(data, uint16(info.Format))
		return data, nil
	}
	sz := 2 + 64 + 512 + 16 + 64 + 192 + 64 + 192 + 128 + 128 + 32 + 4
	switch info.Format {
	case 2:
		sz += 2
	case 3:
		sz += 2 + len(info.QuestIntro) + len(info.QuestGraphics)
	}
	out := make([]byte, sz)
	data := out

	binary.LittleEndian.PutUint16(data, uint16(info.Format))
	data = data[2:]

	for _, f := range []struct {
		val   string
		trail string
		max   int
	}{
		{info.Summary, info.Trailing.Summary, 64},
		{info.Description, info.Trailing.Description, 512},
		{info.Version, info.Trailing.Version, 16},
		{info.Author, info.Trailing.Author, 64},
		{info.Email, info.Trailing.Email, 192},
		{info.Author2, info.Trailing.Author2, 64},
		{info.Email2, info.Trailing.Email2, 192},
		{info.Field7, info.Trailing.Field7, 128},
		{info.Copyright, info.Trailing.Copyright, 128},
		{info.Date, info.Trailing.Date, 32},
	} {
		dst := data[:f.max]
		data = data[f.max:]
		n := copy(dst, f.val)
		if n+1 < len(dst) {
			dst = dst[n+1:]
			if f.trail != "" {
				copy(dst, f.trail)
			}
		}
	}
	binary.LittleEndian.PutUint32(data, info.Flags)
	data = data[4:]
	if info.Format == 2 {
		data[0] = info.MinPlayers
		data[1] = info.MaxPlayers
		data = data[2:]
	}
	if info.Format < 3 {
		return out, nil
	}
	data[0] = byte(len(info.QuestIntro))
	data = data[1:]
	n := copy(data, info.QuestIntro)
	data = data[n:]

	data[0] = byte(len(info.QuestGraphics))
	data = data[1:]
	n = copy(data, info.QuestGraphics)
	data = data[n:]
	return out, nil
}

func (info *MapInfo) UnmarshalBinary(data []byte) error {
	*info = MapInfo{}
	if len(data) < 2 {
		return io.ErrUnexpectedEOF
	}
	info.Format = int(binary.LittleEndian.Uint16(data))
	data = data[2:]
	if info.Format > 3 {
		return fmt.Errorf("unsupported version: %d", info.Format)
	}
	if info.Format < 1 {
		return nil
	}
	for _, f := range []struct {
		p   *string
		tp  *string
		max int
	}{
		{&info.Summary, &info.Trailing.Summary, 64},
		{&info.Description, &info.Trailing.Description, 512},
		{&info.Version, &info.Trailing.Version, 16},
		{&info.Author, &info.Trailing.Author, 64},
		{&info.Email, &info.Trailing.Email, 192},
		{&info.Author2, &info.Trailing.Author2, 64},
		{&info.Email2, &info.Trailing.Email2, 192},
		{&info.Field7, &info.Trailing.Field7, 128},
		{&info.Copyright, &info.Trailing.Copyright, 128},
		{&info.Date, &info.Trailing.Date, 32},
	} {
		if len(data) < f.max {
			return io.ErrUnexpectedEOF
		}
		str := data[:f.max]
		data = data[f.max:]
		if i := bytes.IndexByte(str, 0); i >= 0 {
			if trail := bytes.TrimRight(str[i+1:], "\x00"); f.tp != nil && len(trail) > 0 {
				*f.tp = string(trail)
			}
			str = str[:i]
		}
		*f.p = string(str)
	}
	if len(data) < 4 {
		return io.ErrUnexpectedEOF
	}
	info.Flags = binary.LittleEndian.Uint32(data)
	data = data[4:]
	if info.Format == 2 {
		if len(data) < 2 {
			return io.ErrUnexpectedEOF
		}
		info.MinPlayers = data[0]
		info.MaxPlayers = data[1]
		data = data[2:]
	} else {
		info.MinPlayers = 2
		info.MaxPlayers = 16
	}
	if info.Format < 3 {
		return nil
	}
	if len(data) == 0 {
		return io.ErrUnexpectedEOF
	}

	sz := int(data[0])
	data = data[1:]
	if len(data) < sz {
		return io.ErrUnexpectedEOF
	}
	info.QuestIntro = string(data[:sz])
	data = data[sz:]
	if len(data) == 0 {
		return io.ErrUnexpectedEOF
	}

	sz = int(data[0])
	data = data[1:]
	if len(data) < sz {
		return io.ErrUnexpectedEOF
	}
	info.QuestGraphics = string(data[:sz])
	data = data[sz:]
	return nil
}

type MapIntro struct {
	Data string
}

func (*MapIntro) MapSection() string {
	return "MapIntro"
}

func (sect *MapIntro) MarshalBinary() ([]byte, error) {
	data := make([]byte, 6+len(sect.Data))
	binary.LittleEndian.PutUint16(data[0:], 1)
	binary.LittleEndian.PutUint32(data[2:], uint32(len(sect.Data)))
	copy(data[6:], sect.Data)
	return data, nil
}

func (sect *MapIntro) UnmarshalBinary(data []byte) error {
	if len(data) < 2 {
		return io.ErrUnexpectedEOF
	}
	vers := binary.LittleEndian.Uint16(data)
	data = data[2:]
	if vers > 1 {
		return fmt.Errorf("unsupported map intro version: %d", vers)
	}
	if len(data) < 4 {
		return io.ErrUnexpectedEOF
	}
	*sect = MapIntro{}
	size := binary.LittleEndian.Uint32(data)
	data = data[4:]
	if size > uint32(len(data)) {
		return io.ErrUnexpectedEOF
	}
	sect.Data = string(data[:size])
	return nil
}

type ColorRGB struct {
	R, G, B byte
}

type AmbientData struct {
	AmbientColor ColorRGB
}

func (*AmbientData) MapSection() string {
	return "AmbientData"
}

func (sect *AmbientData) MarshalBinary() ([]byte, error) {
	data := make([]byte, 14)
	binary.LittleEndian.PutUint16(data[0:], 1)
	binary.LittleEndian.PutUint32(data[2:], uint32(sect.AmbientColor.R))
	binary.LittleEndian.PutUint32(data[6:], uint32(sect.AmbientColor.G))
	binary.LittleEndian.PutUint32(data[10:], uint32(sect.AmbientColor.B))
	return data, nil
}

func (sect *AmbientData) UnmarshalBinary(data []byte) error {
	if len(data) < 2 {
		return io.ErrUnexpectedEOF
	}
	vers := binary.LittleEndian.Uint16(data)
	data = data[2:]
	if vers != 1 {
		return fmt.Errorf("unsupported ambient data version: %d", vers)
	}
	if len(data) < 12 {
		return io.ErrUnexpectedEOF
	}
	r := binary.LittleEndian.Uint32(data[0:])
	g := binary.LittleEndian.Uint32(data[4:])
	b := binary.LittleEndian.Uint32(data[8:])
	if r > 0xff || g > 0xff || b > 0xff {
		return fmt.Errorf("invalid color value in ambient data: (%d,%d,%d)", r, g, b)
	}
	sect.AmbientColor = ColorRGB{
		R: byte(r), G: byte(g), B: byte(b),
	}
	return nil
}
