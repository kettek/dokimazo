package game

import (
	"rotate-test/internal/res"

	"github.com/hajimehoshi/ebiten/v2"
)

type Mover struct {
	*SpriteStack
}

func NewMover() *Mover {
	return &Mover{
		//Sprite: NewSpriteFromImage(res.MustLoadImage("thing.png")),
		SpriteStack: NewSpriteStackFromImageSheet(res.NewImageSheet(res.MustLoadImage("humus.png"), 16, 16)),
	}
}

func (m *Mover) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	m.SpriteStack.Draw(screen, geom)
}

func (m *Mover) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		m.Rotate(-0.05)
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		m.Rotate(0.05)
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		m.Add(m.Forward())
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		m.Sub(m.Forward())
	}
}
