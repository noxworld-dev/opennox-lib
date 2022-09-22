package object

import "encoding/json"

var MonsterClassNames = []string{
	"SMALL_MONSTER", "MEDIUM_MONSTER", "LARGE_MONSTER", "SHOPKEEPER", "NPC", "FEMALE_NPC",
	"UNDEAD", "MONITOR", "MIGRATE", "IMMUNE_POISON", "IMMUNE_FIRE", "IMMUNE_ELECTRICITY",
	"IMMUNE_FEAR", "BOMBER", "NO_TARGET", "NO_SPELL_TARGET", "HAS_SOUL", "WARCRY_STUN",
	"LOOK_AROUND", "WOUNDED_NPC",
}

func (c SubClass) AsMonster() MonsterClass {
	return MonsterClass(c)
}

func ParseMonsterClass(s string) (MonsterClass, error) {
	v, err := parseEnum("monster class", s, MonsterClassNames)
	return MonsterClass(v), err
}

func ParseMonsterClassSet(s string) (MonsterClass, error) {
	v, err := parseEnumSet("monster class", s, MonsterClassNames)
	return MonsterClass(v), err
}

var _ enum[MonsterClass] = MonsterClass(0)

type MonsterClass uint32

const (
	MonsterSmall = MonsterClass(1 << iota)
	MonsterMedium
	MonsterLarge
	MonsterShopkeeper
	MonsterNPC
	MonsterFemaleNPC
	MonsterUndead
	MonsterMonitor
	MonsterMigrate
	MonsterImmunePoison
	MonsterImmuneFire
	MonsterImmuneElectricity
	MonsterImmuneFear
	MonsterBomber
	MonsterNoTarget
	MonsterNoSpellTarget
	MonsterHasSoul
	MonsterWarcryStun
	MonsterLookAround
	MonsterWoundedNPC
)

func (c MonsterClass) Has(c2 MonsterClass) bool {
	return c&c2 != 0
}

func (c MonsterClass) HasAny(c2 MonsterClass) bool {
	return c&c2 != 0
}

func (c MonsterClass) Split() []MonsterClass {
	return splitBits(c)
}

func (c MonsterClass) String() string {
	return stringBits(uint32(c), MonsterClassNames)
}

func (c MonsterClass) MarshalJSON() ([]byte, error) {
	var arr []string
	for _, s := range c.Split() {
		arr = append(arr, s.String())
	}
	return json.Marshal(arr)
}
