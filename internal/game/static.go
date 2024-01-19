package game

import (
	"rotate-test/internal/res"

	"github.com/hajimehoshi/ebiten/v2"
)

type Static struct {
	*Sprite
}

func NewStatic() *Static {
	return &Static{
		Sprite: NewSpriteFromImage(res.MustLoadImage("tile.png")),
	}
}

func (s *Static) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	s.Sprite.Draw(screen, geom)
}

func (s *Static) Update() {
}
