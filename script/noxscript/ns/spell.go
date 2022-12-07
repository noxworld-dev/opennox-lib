package ns

import (
	"github.com/noxworld-dev/opennox-lib/script"
	"github.com/noxworld-dev/opennox-lib/script/noxscript/ns/effect"
	"github.com/noxworld-dev/opennox-lib/script/noxscript/ns/spell"
)

// Effect triggers an effect from point p1 to p2.
// Some effects only have one point, in which case p2 is ignored.
func Effect(effect effect.Effect, p1, p2 script.Positioner) {
	if impl == nil {
		return
	}
	impl.Effect(effect, p1, p2)
}

// CastSpell casts a spell from source to target.
//
//  Example:
//    CastSpellObjectObject(spell.DEATH_RAY, Object("CruelDude"), GetHost())
//    CastSpellObjectObject(spell.DEATH_RAY, types.Ptf(10, 5), Waypoint("Target"))
func CastSpell(spell spell.Spell, source, target script.Positioner) {
	if impl == nil {
		return
	}
	impl.CastSpell(spell, source, target)
}
