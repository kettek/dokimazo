package game

import (
	"github.com/kettek/dokimazo/internal/res"
	input "github.com/quasilyte/ebitengine-input"
)

type Player struct {
	input *input.Handler
	*SpriteStack
	Velocity Vec2
	Inventory
	Emotions
	chunk *Chunk
}

func NewPlayer() *Player {
	return &Player{
		SpriteStack: NewSpriteStackFromSheet(res.MustLoadSheet("koinon.png")),
	}
}

func (p *Player) Chunk() *Chunk {
	return p.chunk
}

func (p *Player) SetChunk(chunk *Chunk) {
	p.chunk = chunk
}

func (p *Player) Draw(drawOpts DrawOpts) {
	p.SpriteStack.Draw(drawOpts)
}

func (p *Player) Update() (requests []Request) {
	r := 0.0
	if p.input.ActionIsPressed(InputTurnLeft) {
		r -= 0.05
	} else if p.input.ActionIsPressed(InputTurnRight) {
		r += 0.05
	}
	if r != 0.0 {
		requests = append(requests, RequestRotate{Rotation: r})
	}

	dir := Vec2{}
	if p.input.ActionIsPressed(InputMoveForward) {
		dir = p.Forward()
		p.Velocity.Add(dir)
	} else if p.input.ActionIsPressed(InputMoveBackward) {
		dir = p.Forward()
		dir.Mul(Vec2{-1, -1})
		p.Velocity.Add(dir)
	}
	if p.Velocity.X() != 0.0 || p.Velocity.Y() != 0.0 {
		v := p.Vec2.Clone()
		v.Add(p.Velocity)
		requests = append(requests, RequestMove{From: p.Vec2, To: v})
	}
	p.Velocity.Mul(Vec2{0.5, 0.5})
	if p.Velocity.X() < 0.01 && p.Velocity.X() > -0.01 {
		p.Velocity.Assign(Vec2{0, p.Velocity.Y()})
	}
	if p.Velocity.Y() < 0.01 && p.Velocity.Y() > -0.01 {
		p.Velocity.Assign(Vec2{p.Velocity.X(), 0})
	}

	return
}

func (p *Player) HandleRequest(request Request, success bool) {
	if !success {
		return
	}
	switch request := request.(type) {
	case RequestRotate:
		p.Rotate(request.Rotation)
	case RequestMove:
		p.Assign(request.To)
	}
}
