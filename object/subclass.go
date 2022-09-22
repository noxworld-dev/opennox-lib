package object

import "encoding/json"

var SubClassNames = [][]string{
	ArmorClassNames,
	WeaponClassNames,
	MonsterClassNames,
	FoodClassNames,
	OtherClassNames,
	BookClassNames,
	MissileClassNames,
	GeneratorClassNames,
	ExitClassNames,
}

func ParseSubClass(s string) (SubClass, error) {
	v, err := parseEnumMulti("subclass", s, SubClassNames)
	return SubClass(v), err
}

func ParseSubClassSet(s string) (SubClass, error) {
	v, err := parseEnumSetMulti("subclass", s, SubClassNames)
	return SubClass(v), err
}

var _ enum[SubClass] = SubClass(0)

type SubClass uint32

func (c SubClass) Has(c2 SubClass) bool {
	return c&c2 != 0
}

func (c SubClass) HasAny(c2 SubClass) bool {
	return c&c2 != 0
}

func (c SubClass) Split() []SubClass {
	return splitBits(c)
}

func (c SubClass) String() string {
	return stringBitsRaw(uint32(c))
}

func (c SubClass) MarshalJSON() ([]byte, error) {
	var arr []string
	for _, s := range c.Split() {
		arr = append(arr, s.String())
	}
	return json.Marshal(arr)
}
