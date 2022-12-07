package types

import (
	"image"
	"math"
)

func Ptf(x, y float32) Pointf {
	return Pointf{X: x, Y: y}
}

func Point2f(p image.Point) Pointf {
	return Pointf{
		X: float32(p.X),
		Y: float32(p.Y),
	}
}

type Pointf struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

func (p Pointf) Pos() Pointf {
	return p
}

func (p Pointf) Point() image.Point {
	return image.Point{
		X: int(p.X),
		Y: int(p.Y),
	}
}

func (p Pointf) Add(p2 Pointf) Pointf {
	p.X += p2.X
	p.Y += p2.Y
	return p
}

func (p Pointf) Sub(p2 Pointf) Pointf {
	p.X -= p2.X
	p.Y -= p2.Y
	return p
}

func (p Pointf) Mul(v float32) Pointf {
	p.X *= v
	p.Y *= v
	return p
}

func (p Pointf) Div(v float32) Pointf {
	p.X /= v
	p.Y /= v
	return p
}

func (p Pointf) Len() float64 {
	x, y := float64(p.X), float64(p.Y)
	return math.Hypot(x, y)
}

func (p Pointf) Normalize() Pointf {
	return p.Div(float32(p.Len()))
}

type Rectf struct {
	Left   float32
	Top    float32
	Right  float32
	Bottom float32
}

func (r *Rectf) IsEmpty() bool {
	return r.Right <= r.Left || r.Bottom <= r.Top
}

func RectFromPointsf(p1, p2 Pointf) Rectf {
	var r Rectf
	if p1.X >= p2.X {
		r.Left = p2.X
		r.Right = p1.X
	} else {
		r.Left = p1.X
		r.Right = p2.X
	}
	if p1.Y >= p2.Y {
		r.Top = p2.Y
		r.Bottom = p1.Y
	} else {
		r.Top = p1.Y
		r.Bottom = p2.Y
	}
	return r
}
