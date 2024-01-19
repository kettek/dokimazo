package game

import "github.com/hajimehoshi/ebiten/v2"

type Visual interface {
	Draw(screen *ebiten.Image, geom ebiten.GeoM)
	Position() Vec2
	Size() Vec2
	Angle() float64
	Z() float64
	SetZ(float64)
}

type Visuals []Visual

func (v *Visuals) Add(visual Visual) {
	*v = append(*v, visual)
}

func (v *Visuals) Remove(visual Visual) {
	for i, visual2 := range *v {
		if visual2 == visual {
			*v = append((*v)[:i], (*v)[i+1:]...)
			return
		}
	}
}
