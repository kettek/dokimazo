package game

import (
	"rotate-test/internal/res"
)

type Static struct {
	*Sprite
}

func NewStatic() *Static {
	return &Static{
		Sprite: NewSpriteFromSheet(res.MustLoadSheetWithSize("tile.png", 20, 20)),
	}
}

func (s *Static) Draw(drawOpts DrawOpts) {
	s.Sprite.Draw(drawOpts)
}

func (s *Static) Update() {
}
