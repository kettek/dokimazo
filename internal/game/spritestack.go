package game

import (
	"image/color"
	"math"

	"github.com/kettek/dokimazo/internal/res"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteStack struct {
	RVec2
	*res.Sheet
	LayerDistance float64
	z             float64
}

func NewSpriteStackFromSheet(sheet *res.Sheet) *SpriteStack {
	return &SpriteStack{
		Sheet:         sheet,
		LayerDistance: 1,
		z:             1,
	}
}

func (s *SpriteStack) Clone() *SpriteStack {
	return &SpriteStack{
		RVec2:         s.RVec2,
		Sheet:         s.Sheet,
		z:             s.z,
		LayerDistance: s.LayerDistance,
	}
}

func (s *SpriteStack) Z() float64 {
	return s.z
}

func (s *SpriteStack) SetZ(z float64) {
	s.z = z
}

func (s *SpriteStack) Draw(drawOpts DrawOpts) {
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear

	//s.DrawShadow(drawOpts)

	// Rotate about center.
	op.GeoM.Translate(-s.HalfWidth(), -s.HalfHeight())
	op.GeoM.Rotate(s.Angle())
	op.GeoM.Translate(s.HalfWidth(), s.HalfHeight())
	// Translate to position.
	op.GeoM.Translate(s.X(), s.Y())
	op.GeoM.Concat(drawOpts.GeoM)
	for col := 0; col < s.Sheet.Cols(); col++ {
		op.ColorScale.Reset()
		r := float64(col) / float64(s.Sheet.Cols())
		c := uint8(150.0 + 105*r)
		op.ColorScale.ScaleWithColor(color.NRGBA{c, c, c, 255})
		drawOpts.Image.DrawImage(s.Sheet.At(col, 0), op)
		op.GeoM.Translate(0, -s.LayerDistance*drawOpts.Z)
	}
}

func (s *SpriteStack) DrawShadow(drawOpts DrawOpts) {
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear

	// FIXME: Move this shadow elsewhere. It should be rendered to its own "shadow" layer... potentially we should separate the world into "Y" layers, and then have each Y layer have its own shadow. This would be a bit expensive, but would look nice.
	// Let's try a lil shadow.
	op.GeoM.Translate(-s.HalfWidth(), -s.HalfHeight())
	op.GeoM.Rotate(s.Angle())
	op.GeoM.Translate(s.HalfWidth(), s.HalfHeight())
	// Translate to position.
	op.GeoM.Translate(s.X(), s.Y())
	op.GeoM.Concat(drawOpts.GeoM)

	for col := 0; col < s.Sheet.Cols(); col++ {
		op.ColorScale.Reset()
		r := float64(col) / float64(s.Sheet.Cols())
		c := 230 - uint8(100.0+105*r)
		op.ColorScale.ScaleWithColor(color.NRGBA{0, 0, 0, c})
		drawOpts.Image.DrawImage(s.Sheet.At(col, 0), op)
		// NOTE: Instead of passing in R and Z as drawOpts, we could just decompose the matrix into its components and use those. The math is a bit over my head, but should be possible.
		x := math.Cos(-drawOpts.Angle+0.5) * drawOpts.Z * s.LayerDistance
		y := math.Sin(-drawOpts.Angle+0.5) * drawOpts.Z * s.LayerDistance
		op.GeoM.Translate(x, y)
	}
}

func (s *SpriteStack) Position() Vec2 {
	return Vec2{s.X(), s.Y()}
}

func (s *SpriteStack) Size() Vec2 {
	return Vec2{float64(s.Bounds().Dx()), float64(s.Bounds().Dy())}
}
