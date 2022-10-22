package ns

type Spell = string

// AwardSpell awards spell level to object.
//
// This will raise the spell level of the object.
// If the object can not cast this spell then it will have no effect.
func AwardSpell(id Obj, spell Spell) bool {
	// header only
	return false
}

// GroupAwardSpell awards spell level to objects in a group.
//
// This will raise the spell level of the objects in the group.
// If an object can not cast this spell then it will have no effect on that object.
func GroupAwardSpell(group ObjGroup, spell Spell) {
	// header only
}

type EnchantType = string

// HasEnchant gets whether object has an enchant.
func HasEnchant(obj Obj, enchant EnchantType) bool {
	// header only
	return false
}

// Enchant grants object an enchantment of a specified duration.
func Enchant(obj Obj, enchant EnchantType, sec float32) {
	// header only
}

// GroupEnchant grants objects in a group an enchantment of a specified duration.
func GroupEnchant(group ObjGroup, enchant EnchantType, sec float32) {
	// header only
}

// EnchantOff removes enchant from an object.
func EnchantOff(obj Obj, enchant EnchantType) {
	// header only
}

type EffectType = string

// Effect triggers an effect from point (x1,y1) to (x2,y2).
// Some effects only have one point, in which case (x2,y2) is ignored.
func Effect(effect EffectType, x1 float32, y1 float32, x2 float32, y2 float32) {
	// header only
}

// CastSpellObjectObject casts a spell from source object to target object.
//
//  Example:
//    CastSpellObjectObject(spell.DEATH_RAY, Object("CruelDude"), GetHost())
func CastSpellObjectObject(spell Spell, source Obj, target Obj) {
	// header only
}

// CastSpellObjectLocation casts a spell from source object to target location (x,y).
func CastSpellObjectLocation(spell Spell, source Obj, x float32, y float32) {
	// header only
}

// CastSpellLocationObject casts a spell from source location (x,y) to target object.
func CastSpellLocationObject(spell Spell, x float32, y float32, target Obj) {
	// header only
}

// CastSpellLocationLocation casts a spell from source location (x1,y1) to target location (x2,y2).
func CastSpellLocationLocation(spell Spell, x1 float32, y1 float32, x2 float32, y2 float32) {
	// header only
}

// TrapSpells sets spells on a bomber.
func TrapSpells(obj Obj, spell1 Spell, spell2 Spell, spell3 Spell) {
	// header only
}
