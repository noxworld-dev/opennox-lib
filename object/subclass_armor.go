package object

import "encoding/json"

var ArmorClassNames = []string{
	"HELMET", "SHIELD", "BREASTPLATE", "ARM_ARMOR", "PANTS", "BOOTS", "SHIRT", "LEG_ARMOR", "BACK",
}

func (c SubClass) AsArmor() ArmorClass {
	return ArmorClass(c)
}

func ParseArmorClass(s string) (ArmorClass, error) {
	v, err := parseEnum("armor class", s, ArmorClassNames)
	return ArmorClass(v), err
}

func ParseArmorClassSet(s string) (ArmorClass, error) {
	v, err := parseEnumSet("armor class", s, ArmorClassNames)
	return ArmorClass(v), err
}

var _ enum[ArmorClass] = ArmorClass(0)

type ArmorClass uint32

const (
	ArmorHelmet = ArmorClass(1 << iota)
	ArmorShield
	ArmorBreastplate
	ArmorArmArmor
	ArmorPants
	ArmorBoots
	ArmorShirt
	ArmorLegArmor
	ArmorBack
)

func (c ArmorClass) Has(c2 ArmorClass) bool {
	return c&c2 != 0
}

func (c ArmorClass) HasAny(c2 ArmorClass) bool {
	return c&c2 != 0
}

func (c ArmorClass) Split() []ArmorClass {
	return splitBits(c)
}

func (c ArmorClass) String() string {
	return stringBits(uint32(c), ArmorClassNames)
}

func (c ArmorClass) MarshalJSON() ([]byte, error) {
	var arr []string
	for _, s := range c.Split() {
		arr = append(arr, s.String())
	}
	return json.Marshal(arr)
}
