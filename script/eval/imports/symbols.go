package imports

import "reflect"

//go:generate yaegi extract errors math math/rand time sort
//go:generate yaegi extract bytes strings strconv unicode io fmt
//go:generate yaegi extract context image/color
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/types
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/object
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/player
//go:generate yaegi extract github.com/noxworld-dev/opennox-lib/script
//go:generate goimports -w .

var Symbols = make(map[string]map[string]reflect.Value)

// This block ensures that we compile the stdlib bindings.
// It may fail if the compiler version doesn't match the one used to generate bindings.
// To fix this, we remove build tags from stdlib files and remove any bindings introduced in newer Go versions.

var (
	_ _fmt_Formatter
	_ _io_Closer
	_ _math_rand_Source
	_ _sort_Interface
	_ _image_color_Color
)
