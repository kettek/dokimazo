package game

import (
	"math"
	"rotate-test/internal/res"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	camera  Camera
	visuals Visuals
	things  Things
}

func New() *Game {
	g := &Game{
		camera: *NewCamera(),
	}

	var statics []*Static

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			s := NewStatic()
			s.Sprite = NewSpriteFromImage(res.MustLoadImage("dirt.png"))
			s.Assign(Vec2{float64(i) * 16, float64(j) * 16})
			s.SetZ(-10000)
			statics = append(statics, s)
		}
	}

	for _, s := range statics {
		g.visuals.Add(s)
	}

	{
		c := NewSpriteStackFromImageSheet(res.NewImageSheet(res.MustLoadImage("palisade.png"), 16, 16))
		c.Rotate(math.Pi / 2)
		c.SetZ(10000)

		for i := 0; i < 10; i++ {
			if i > 0 && i < 9 {
				continue
			}
			for j := 0; j < 10; j++ {
				nc := c.Clone()
				nc.Assign(Vec2{float64(i) * nc.HalfWidth() * 2, float64(j) * nc.HalfWidth() * 2})
				g.visuals.Add(nc)
			}
		}
	}

	m := NewMover()
	m.Assign(Vec2{200, 200})
	g.visuals.Add(m)
	g.things.Add(m)

	g.camera.Target = m

	return g
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.camera.Layout(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	g.camera.Update()
	for _, t := range g.things {
		t.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.camera.Draw(screen, g.visuals)
}
