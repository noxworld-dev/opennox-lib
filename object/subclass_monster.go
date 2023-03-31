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
	MonsterSmall             = MonsterClass(1 << iota) // 0x1
	MonsterMedium                                      // 0x2
	MonsterLarge                                       // 0x4
	MonsterShopkeeper                                  // 0x8
	MonsterNPC                                         // 0x10
	MonsterFemaleNPC                                   // 0x20
	MonsterUndead                                      // 0x40
	MonsterMonitor                                     // 0x80
	MonsterMigrate                                     // 0x100
	MonsterImmunePoison                                // 0x200
	MonsterImmuneFire                                  // 0x400
	MonsterImmuneElectricity                           // 0x800
	MonsterImmuneFear                                  // 0x1000
	MonsterBomber                                      // 0x2000
	MonsterNoTarget                                    // 0x4000
	MonsterNoSpellTarget                               // 0x8000
	MonsterHasSoul                                     // 0x10000
	MonsterWarcryStun                                  // 0x20000
	MonsterLookAround                                  // 0x40000
	MonsterWoundedNPC                                  // 0x80000
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
