package imports

import (
	"reflect"
)

//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/types
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/object
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/wall
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/player
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/script
//go:generate goimports -w .

var Symbols = make(map[string]map[string]reflect.Value)
