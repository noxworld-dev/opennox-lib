package object

import (
	"encoding/json"
)

var MaterialNames = []string{
	"FLESH", "CLOTH", "ANIMAL_HIDE", "WOOD", "METAL", "STONE", "EARTH", "LIQUID",
	"GLASS", "PAPER", "SNOW", "MUD", "MAGIC", "DIAMOND", "NONE",
}

var _ enum[Material] = Material(0)

func ParseMaterial(s string) (Material, error) {
	v, err := parseEnum("material", s, MaterialNames)
	return Material(v), err
}

func ParseMaterialSet(s string) (Material, error) {
	v, err := parseEnumSet("material", s, MaterialNames)
	return Material(v), err
}

type Material uint32

const (
	MaterialFlesh = Material(1 << iota)
	MaterialCloth
	MaterialAnimalHide
	MaterialWood
	MaterialMetal
	MaterialStone
	MaterialEarth
	MaterialLiquid
	MaterialGlass
	MaterialPaper
	MaterialSnow
	MaterialMud
	MaterialMagic
	MaterialDiamond
	MaterialNone
)

func (c Material) Has(c2 Material) bool {
	return c&c2 != 0
}

func (c Material) HasAny(c2 Material) bool {
	return c&c2 != 0
}

func (c Material) Split() []Material {
	return splitBits(c)
}

func (c Material) String() string {
	return stringBits(uint32(c), MaterialNames)
}

func (c Material) MarshalJSON() ([]byte, error) {
	var arr []string
	for _, s := range c.Split() {
		arr = append(arr, s.String())
	}
	return json.Marshal(arr)
}
