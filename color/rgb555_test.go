package noxcolor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRGB555(t *testing.T) {
	for i := 0; i < 0x7FFF; i++ {
		c := RGB555(i)
		cl := c.ColorNRGBA()
		c2 := RGB555Color(cl.R, cl.G, cl.B)
		require.Equal(t, c, c2, "0x%x", int(c))
	}
}
