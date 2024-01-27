package noxnet

import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"
)

func allZeros(data []byte) bool {
	for _, v := range data {
		if v != 0 {
			return false
		}
	}
	return true
}

func cLen(data []byte) int {
	if i := bytes.IndexByte(data, 0); i >= 0 {
		return i
	}
	return len(data)
}

func cString(data []byte) string {
	data = data[:cLen(data)]
	return string(data)
}

func cStringSet0(data []byte, s string) {
	zeros(data)
	copy(data, s)
	data[len(data)-1] = 0
}

func cStringSet(data []byte, s string) int {
	zeros(data)
	return copy(data, s)
}

func cLen16(data []byte) int {
	for i := 0; i < len(data)/2; i++ {
		v := binary.LittleEndian.Uint16(data[2*i:])
		if v == 0 {
			return i
		}
	}
	return len(data)
}

func cString16(data []byte) string {
	data16 := make([]uint16, len(data)/2)
	for i := range data16 {
		v := binary.LittleEndian.Uint16(data[2*i:])
		if v == 0 {
			data16 = data16[:i]
			break
		}
		data16[i] = v
	}
	return string(utf16.Decode(data16))
}

func cStringSet16(data []byte, s string) {
	zeros(data)
	data16 := utf16.Encode([]rune(s))
	for i, v := range data16 {
		binary.LittleEndian.PutUint16(data[2*i:], v)
	}
}

func zeros(data []byte) {
	for i := range data {
		data[i] = 0
	}
}

type FixedString struct {
	Value string
	Junk  []byte
}

func (s *FixedString) Encode(data []byte) {
	i := cStringSet(data, s.Value)
	if len(s.Junk) != 0 {
		copy(data[i+1:], s.Junk)
	}
}

func (s *FixedString) Decode(data []byte) {
	s.Value = cString(data)
	s.Junk = nil
	if i := len(s.Value); !allZeros(data[i+1:]) {
		s.Junk = make([]byte, len(data)-i-1)
		copy(s.Junk, data[i+1:])
	}
}
