package object

import "encoding/json"

var WeaponClassNames = []string{
	"FLAG",
	"QUIVER",
	"BOW",
	"CROSSBOW",
	"ARROW",
	"BOLT",
	"CHAKRAM",
	"SHURIKEN",
	"SWORD",
	"LONG_SWORD",
	"GREAT_SWORD",
	"MACE",
	"AXE",
	"OGRE_AXE",
	"HAMMER",
	"STAFF",
	"STAFF_SULPHOROUS_FLARE",
	"STAFF_SULPHOROUS_SHOWER",
	"STAFF_LIGHTNING",
	"STAFF_FIREBALL",
	"STAFF_TRIPLE_FIREBALL",
	"STAFF_FORCE_OF_NATURE",
	"STAFF_DEATH_RAY",
	"STAFF_OBLIVION_HALBERD",
	"STAFF_OBLIVION_HEART",
	"STAFF_OBLIVION_WIERDLING",
	"STAFF_OBLIVION_ORB",
}

func (c SubClass) AsWeapon() WeaponClass {
	return WeaponClass(c)
}

func ParseWeaponClass(s string) (WeaponClass, error) {
	v, err := parseEnum("weapon class", s, WeaponClassNames)
	return WeaponClass(v), err
}

func ParseWeaponClassSet(s string) (WeaponClass, error) {
	v, err := parseEnumSet("weapon class", s, WeaponClassNames)
	return WeaponClass(v), err
}

var _ enum[WeaponClass] = WeaponClass(0)

type WeaponClass uint32

const (
	WeaponFlag = WeaponClass(1 << iota)
	WeaponQuiver
	WeaponBow
	WeaponCrossbow
	WeaponArrow
	WeaponBolt
	WeaponChakram
	WeaponShuriken
	WeaponSword
	WeaponLongSword
	WeaponGreatSword
	WeaponMace
	WeaponAxe
	WeaponOgreAxe
	WeaponHammer
	WeaponStaff
	WeaponStaffSulphorousFlare
	WeaponStaffSulphorousShower
	WeaponStaffLightning
	WeaponStaffFireball
	WeaponStaffTripleFireball
	WeaponStaffForceOfNature
	WeaponStaffDeathRay
	WeaponStaffOblivionHalberd
	WeaponStaffOblivionHeart
	WeaponStaffOblivionWierdling
	WeaponStaffOblivionOrb
)

func (c WeaponClass) Has(c2 WeaponClass) bool {
	return c&c2 != 0
}

func (c WeaponClass) HasAny(c2 WeaponClass) bool {
	return c&c2 != 0
}

func (c WeaponClass) Split() []WeaponClass {
	return splitBits(c)
}

func (c WeaponClass) String() string {
	return stringBits(uint32(c), WeaponClassNames)
}

func (c WeaponClass) MarshalJSON() ([]byte, error) {
	var arr []string
	for _, s := range c.Split() {
		arr = append(arr, s.String())
	}
	return json.Marshal(arr)
}
