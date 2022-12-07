package imports

import "reflect"

//go:generate yaegi extract errors math math/rand time sort
//go:generate yaegi extract bytes strings strconv unicode io fmt
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/types
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/object
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/player
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/script
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/script/noxscript/ns
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/script/noxscript/ns/audio
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/script/noxscript/ns/class
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/script/noxscript/ns/damage
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/script/noxscript/ns/effect
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/script/noxscript/ns/enchant
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/script/noxscript/ns/spell
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/script/noxscript/ns/subclass
//go:generate goimports -w .

var Symbols = make(map[string]map[string]reflect.Value)
