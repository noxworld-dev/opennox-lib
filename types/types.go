package types

import (
	"image"
	"math"
)

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

func (p Pointf) Len() float64 {
	x, y := float64(p.X), float64(p.Y)
	return math.Sqrt(x*x + y*y)
}

type Rect struct {
	Left   int
	Top    int
	Right  int
	Bottom int
}

func (r *Rect) IsEmpty() bool {
	return r.Right <= r.Left || r.Bottom <= r.Top
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

func UtilRectXxx(r1, r2 Rect) (Rect, bool) { // nox_xxx_utilRect_49F930
	left := r2.Left
	if r1.Left >= left {
		left = r1.Left
	}
	right := r2.Right
	if r1.Right <= right {
		right = r1.Right
	}
	if left >= right {
		return Rect{}, false
	}
	top := r2.Top
	if r1.Top >= top {
		top = r1.Top
	}
	bottom := r2.Bottom
	if r1.Bottom <= bottom {
		bottom = r1.Bottom
	}
	if top >= bottom {
		return Rect{}, false
	}
	return Rect{
		Left:   left,
		Top:    top,
		Right:  right,
		Bottom: bottom,
	}, true
}
