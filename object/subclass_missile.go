package object

import "encoding/json"

var MissileClassNames = []string{
	"MISSILE_COUNTERSPELL",
	"MAGIC",
}

func (c SubClass) AsMissile() MissileClass {
	return MissileClass(c)
}

func ParseMissileClass(s string) (MissileClass, error) {
	v, err := parseEnum("missile class", s, MissileClassNames)
	return MissileClass(v), err
}

func ParseMissileClassSet(s string) (MissileClass, error) {
	v, err := parseEnumSet("missile class", s, MissileClassNames)
	return MissileClass(v), err
}

var _ enum[MissileClass] = MissileClass(0)

type MissileClass uint32

const (
	MissileMissileCounterSpell = MissileClass(1 << iota) // 0x1
	MissileMagic                                         // 0x2
)

func (c MissileClass) Has(c2 MissileClass) bool {
	return c&c2 != 0
}

func (c MissileClass) HasAny(c2 MissileClass) bool {
	return c&c2 != 0
}

func (c MissileClass) Split() []MissileClass {
	return splitBits(c)
}

func (c MissileClass) String() string {
	return stringBits(uint32(c), MissileClassNames)
}

func (c MissileClass) MarshalJSON() ([]byte, error) {
	var arr []string
	for _, s := range c.Split() {
		arr = append(arr, s.String())
	}
	return json.Marshal(arr)
}
