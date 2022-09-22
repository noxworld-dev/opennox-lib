package object

import "encoding/json"

var ExitClassNames = []string{
	"QUEST_EXIT",
	"QUEST_WARP_EXIT",
}

func (c SubClass) AsExit() ExitClass {
	return ExitClass(c)
}

func ParseExitClass(s string) (ExitClass, error) {
	v, err := parseEnum("exit class", s, ExitClassNames)
	return ExitClass(v), err
}

func ParseExitClassSet(s string) (ExitClass, error) {
	v, err := parseEnumSet("exit class", s, ExitClassNames)
	return ExitClass(v), err
}

var _ enum[ExitClass] = ExitClass(0)

type ExitClass uint32

const (
	ExitQuest = ExitClass(1 << iota)
	ExitQuestWarp
)

func (c ExitClass) Has(c2 ExitClass) bool {
	return c&c2 != 0
}

func (c ExitClass) HasAny(c2 ExitClass) bool {
	return c&c2 != 0
}

func (c ExitClass) Split() []ExitClass {
	return splitBits(c)
}

func (c ExitClass) String() string {
	return stringBits(uint32(c), ExitClassNames)
}

func (c ExitClass) MarshalJSON() ([]byte, error) {
	var arr []string
	for _, s := range c.Split() {
		arr = append(arr, s.String())
	}
	return json.Marshal(arr)
}
