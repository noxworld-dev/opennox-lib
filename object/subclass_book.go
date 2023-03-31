package object

import "encoding/json"

var BookClassNames = []string{
	"SPELL_BOOK",
	"FIELD_GUIDE",
	"ABILITY_BOOK",
}

func (c SubClass) AsBook() BookClass {
	return BookClass(c)
}

func ParseBookClass(s string) (BookClass, error) {
	v, err := parseEnum("book class", s, BookClassNames)
	return BookClass(v), err
}

func ParseBookClassSet(s string) (BookClass, error) {
	v, err := parseEnumSet("book class", s, BookClassNames)
	return BookClass(v), err
}

var _ enum[BookClass] = BookClass(0)

type BookClass uint32

const (
	BookSpell      = BookClass(1 << iota) // 0x1
	BookFieldGuide                        // 0x2
	BookAbility                           // 0x4
)

func (c BookClass) Has(c2 BookClass) bool {
	return c&c2 != 0
}

func (c BookClass) HasAny(c2 BookClass) bool {
	return c&c2 != 0
}

func (c BookClass) Split() []BookClass {
	return splitBits(c)
}

func (c BookClass) String() string {
	return stringBits(uint32(c), BookClassNames)
}

func (c BookClass) MarshalJSON() ([]byte, error) {
	var arr []string
	for _, s := range c.Split() {
		arr = append(arr, s.String())
	}
	return json.Marshal(arr)
}
