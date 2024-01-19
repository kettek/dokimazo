package game

import (
	"rotate-test/internal/res"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	RVec2
	*res.Image
	z float64
}

func NewSpriteFromImage(image *res.Image) *Sprite {
	return &Sprite{
		Image: image,
		z:     1,
	}
}

func (s *Sprite) Z() float64 {
	return s.z
}

func (s *Sprite) SetZ(z float64) {
	s.z = z
}

func (s *Sprite) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	// Rotate about center.
	op.GeoM.Translate(-s.HalfWidth(), -s.HalfHeight())
	op.GeoM.Rotate(s.Angle())
	op.GeoM.Translate(s.HalfWidth(), s.HalfHeight())
	// Translate to position.
	op.GeoM.Translate(s.X(), s.Y())
	op.GeoM.Concat(geom)
	screen.DrawImage(s.Image.Image, op)
}

func (s *Sprite) Position() Vec2 {
	return Vec2{s.X(), s.Y()}
}

func (s *Sprite) Size() Vec2 {
	return Vec2{float64(s.Bounds().Dx()), float64(s.Bounds().Dy())}
}
