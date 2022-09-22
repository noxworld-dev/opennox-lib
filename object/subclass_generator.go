package object

import "encoding/json"

var GeneratorClassNames = []string{
	"GENERATOR_NW",
	"GENERATOR_NE",
	"GENERATOR_SE",
	"GENERATOR_SW",
}

func (c SubClass) AsGenerator() GeneratorClass {
	return GeneratorClass(c)
}

func ParseGeneratorClass(s string) (GeneratorClass, error) {
	v, err := parseEnum("generator class", s, GeneratorClassNames)
	return GeneratorClass(v), err
}

func ParseGeneratorClassSet(s string) (GeneratorClass, error) {
	v, err := parseEnumSet("generator class", s, GeneratorClassNames)
	return GeneratorClass(v), err
}

var _ enum[GeneratorClass] = GeneratorClass(0)

type GeneratorClass uint32

const (
	GeneratorNW = GeneratorClass(1 << iota)
	GeneratorNE
	GeneratorSE
	GeneratorSW
)

func (c GeneratorClass) Has(c2 GeneratorClass) bool {
	return c&c2 != 0
}

func (c GeneratorClass) HasAny(c2 GeneratorClass) bool {
	return c&c2 != 0
}

func (c GeneratorClass) Split() []GeneratorClass {
	return splitBits(c)
}

func (c GeneratorClass) String() string {
	return stringBits(uint32(c), GeneratorClassNames)
}

func (c GeneratorClass) MarshalJSON() ([]byte, error) {
	var arr []string
	for _, s := range c.Split() {
		arr = append(arr, s.String())
	}
	return json.Marshal(arr)
}
