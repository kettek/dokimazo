package game

import (
	"rotate-test/internal/res"
)

type Static struct {
	*Sprite
}

func NewStatic() *Static {
	return &Static{
		Sprite: NewSpriteFromImage(res.MustLoadImage("tile.png")),
	}
}

func (s *Static) Draw(drawOpts DrawOpts) {
	s.Sprite.Draw(drawOpts)
}

func (s *Static) Update() {
}
