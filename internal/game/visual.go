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
