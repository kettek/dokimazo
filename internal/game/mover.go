package game

import (
	"rotate-test/internal/res"

	"github.com/hajimehoshi/ebiten/v2"
)

type Mover struct {
	chunk *Chunk
	*SpriteStack
}

func NewMover() *Mover {
	return &Mover{
		SpriteStack: NewSpriteStackFromSheet(res.MustLoadSheet("koinon.png")),
	}
}

func (m *Mover) SetChunk(chunk *Chunk) {
	m.chunk = chunk
}

func (m *Mover) Chunk() *Chunk {
	return m.chunk
}

func (m *Mover) Draw(drawOpts DrawOpts) {
	m.SpriteStack.Draw(drawOpts)
}

func (m *Mover) Update() (actions []Request) {
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
	return
}
