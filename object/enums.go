package object

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type enum[T ~uint32] interface {
	json.Marshaler
	Has(c2 T) bool
	HasAny(c2 T) bool
	Split() []T
	String() string
}

func splitBits[T ~uint32](v T) (out []T) {
	for i := 0; i < 32; i++ {
		v2 := T(1 << i)
		if v&v2 != 0 {
			out = append(out, v2)
		}
	}
	return
}

func stringBits(v uint32, names []string) string {
	var out []string
	for i := 0; i < 32; i++ {
		v2 := uint32(1 << i)
		if v&v2 != 0 {
			if i < len(names) {
				out = append(out, names[i])
			} else {
				out = append(out, "0x"+strconv.FormatUint(uint64(v2), 16))
			}
		}
	}
	return strings.Join(out, " | ")
}

func stringBitsRaw(v uint32) string {
	var out []string
	for i := 0; i < 32; i++ {
		v2 := uint32(1 << i)
		if v&v2 != 0 {
			out = append(out, "0x"+strconv.FormatUint(uint64(v2), 16))
		}
	}
	return strings.Join(out, " | ")
}

func parseEnum(ename string, s string, names []string) (uint32, error) {
	s = strings.ToUpper(s)
	if s == "" || s == "NULL" {
		return 0, nil
	}
	for i, v := range names {
		if s == v {
			return 1 << i, nil
		}
	}
	return 0, fmt.Errorf("invalid %s name: %q", ename, s)
}

func parseEnumSet(ename string, s string, names []string) (uint32, error) {
	var (
		out  uint32
		last error
	)
	for _, w := range strings.Split(s, "+") {
		w = strings.TrimSpace(w)
		v, err := parseEnum(ename, w, names)
		if err != nil {
			last = err
		}
		out |= v
	}
	return out, last
}

func parseEnumMulti(ename string, s string, lists [][]string) (uint32, error) {
	s = strings.ToUpper(s)
	if s == "" || s == "NULL" {
		return 0, nil
	}
	for _, names := range lists {
		for i, name := range names {
			if s == name {
				return 1 << i, nil
			}
		}
	}
	return 0, fmt.Errorf("invalid %s name: %q", ename, s)
}

func parseEnumSetMulti(ename string, s string, lists [][]string) (uint32, error) {
	var (
		out  uint32
		last error
	)
	for _, w := range strings.Split(s, "+") {
		w = strings.TrimSpace(w)
		v, err := parseEnumMulti(ename, w, lists)
		if err != nil {
			last = err
		}
		out |= v
	}
	return out, last
}
