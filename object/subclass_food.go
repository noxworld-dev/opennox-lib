package object

import "encoding/json"

var FoodClassNames = []string{
	"SIMPLE", "APPLE", "JUG",
	"POTION", "HEALTH_POTION", "MANA_POTION", "CURE_POISON_POTION",
	"MUSHROOM",
	"HASTE_POTION", "INVISIBILITY_POTION", "SHIELD_POTION",
	"FIRE_PROTECT_POTION", "SHOCK_PROTECT_POTION", "POISON_PROTECT_POTION",
	"INVULNERABILITY_POTION", "INFRAVISION_POTION", "VAMPIRISM_POTION",
}

func (c SubClass) AsFood() FoodClass {
	return FoodClass(c)
}

func ParseFoodClass(s string) (FoodClass, error) {
	v, err := parseEnum("food class", s, FoodClassNames)
	return FoodClass(v), err
}

func ParseFoodClassSet(s string) (FoodClass, error) {
	v, err := parseEnumSet("food class", s, FoodClassNames)
	return FoodClass(v), err
}

var _ enum[FoodClass] = FoodClass(0)

type FoodClass uint32

const (
	FoodSimple = FoodClass(1 << iota)
	FoodApple
	FoodJug
	FoodPotion
	FoodHealthPotion
	FoodManaPotion
	FoodCurePoisonPotion
	FoodMushroom
	FoodHastePotion
	FoodInvisibilityPotion
	FoodShieldPotion
	FoodFireProtectPotion
	FoodShockProtectPotion
	FoodPoisonProtectPotion
	FoodInvulnerabilityPotion
	FoodInfravisionPotion
	FoodVampirismPotion
)

func (c FoodClass) Has(c2 FoodClass) bool {
	return c&c2 != 0
}

func (c FoodClass) HasAny(c2 FoodClass) bool {
	return c&c2 != 0
}

func (c FoodClass) Split() []FoodClass {
	return splitBits(c)
}

func (c FoodClass) String() string {
	return stringBits(uint32(c), FoodClassNames)
}

func (c FoodClass) MarshalJSON() ([]byte, error) {
	var arr []string
	for _, s := range c.Split() {
		arr = append(arr, s.String())
	}
	return json.Marshal(arr)
}
