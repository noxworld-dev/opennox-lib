package wall

import (
	"image"

	"github.com/noxworld-dev/opennox-lib/common"
	"github.com/noxworld-dev/opennox-lib/types"
)

const GridStep = common.GridStep

func PosToGrid(pos types.Pointf) image.Point {
	return image.Point{
		X: int(pos.X / GridStep),
		Y: int(pos.Y / GridStep),
	}
}

func GridToPos(pos image.Point) types.Pointf {
	return types.Pointf{
		X: float32(pos.X) * GridStep,
		Y: float32(pos.Y) * GridStep,
	}
}
