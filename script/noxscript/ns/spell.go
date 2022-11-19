package ns

import (
	"github.com/noxworld-dev/opennox-lib/script/noxscript/ns/effect"
	"github.com/noxworld-dev/opennox-lib/script/noxscript/ns/enchant"
	"github.com/noxworld-dev/opennox-lib/script/noxscript/ns/spell"
)

// AwardSpell awards spell level to object.
//
// This will raise the spell level of the object.
// If the object can not cast this spell then it will have no effect.
func AwardSpell(id Obj, spell spell.Spell) bool {
	// header only
	return false
}

// GroupAwardSpell awards spell level to objects in a group.
//
// This will raise the spell level of the objects in the group.
// If an object can not cast this spell then it will have no effect on that object.
func GroupAwardSpell(group ObjGroup, spell spell.Spell) {
	// header only
}

// HasEnchant gets whether object has an enchant.
func HasEnchant(obj Obj, enchant enchant.Enchant) bool {
	// header only
	return false
}

// Enchant grants object an enchantment of a specified duration.
func Enchant(obj Obj, enchant enchant.Enchant, sec float32) {
	// header only
}

// GroupEnchant grants objects in a group an enchantment of a specified duration.
func GroupEnchant(group ObjGroup, enchant enchant.Enchant, sec float32) {
	// header only
}

// EnchantOff removes enchant from an object.
func EnchantOff(obj Obj, enchant enchant.Enchant) {
	// header only
}

// Effect triggers an effect from point (x1,y1) to (x2,y2).
// Some effects only have one point, in which case (x2,y2) is ignored.
func Effect(effect effect.Effect, x1 float32, y1 float32, x2 float32, y2 float32) {
	// header only
}

// CastSpellObjectObject casts a spell from source object to target object.
//
//  Example:
//    CastSpellObjectObject(spell.DEATH_RAY, Object("CruelDude"), GetHost())
func CastSpellObjectObject(spell spell.Spell, source Obj, target Obj) {
	// header only
}

// CastSpellObjectLocation casts a spell from source object to target location (x,y).
func CastSpellObjectLocation(spell spell.Spell, source Obj, x float32, y float32) {
	// header only
}

// CastSpellLocationObject casts a spell from source location (x,y) to target object.
func CastSpellLocationObject(spell spell.Spell, x float32, y float32, target Obj) {
	// header only
}

// CastSpellLocationLocation casts a spell from source location (x1,y1) to target location (x2,y2).
func CastSpellLocationLocation(spell spell.Spell, x1 float32, y1 float32, x2 float32, y2 float32) {
	// header only
}

// TrapSpells sets spells on a bomber.
func TrapSpells(obj Obj, spell1 spell.Spell, spell2 spell.Spell, spell3 spell.Spell) {
	// header only
}
