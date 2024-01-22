package game

import (
	"rotate-test/internal/res"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	RVec2
	*res.Sheet
	z float64
}

func NewSpriteFromSheet(sheet *res.Sheet) *Sprite {
	return &Sprite{
		Sheet: sheet,
		z:     1,
	}
}

func (s *Sprite) Z() float64 {
	return s.z
}

func (s *Sprite) SetZ(z float64) {
	s.z = z
}

func (s *Sprite) Draw(drawOpts DrawOpts) {
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	// Rotate about center.
	op.GeoM.Translate(-s.HalfWidth(), -s.HalfHeight())
	op.GeoM.Rotate(s.Angle())
	op.GeoM.Translate(s.HalfWidth(), s.HalfHeight())
	// Translate to position.
	op.GeoM.Translate(s.X(), s.Y())
	op.GeoM.Concat(drawOpts.GeoM)
	drawOpts.Image.DrawImage(s.Sheet.At(0, 0), op)
}

func (s *Sprite) Position() Vec2 {
	return Vec2{s.X(), s.Y()}
}

func (s *Sprite) Size() Vec2 {
	return Vec2{float64(s.Bounds().Dx()), float64(s.Bounds().Dy())}
}
