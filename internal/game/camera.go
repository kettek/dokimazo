package game

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	RVec2
	W, H     float64
	Z        float64
	image    *ebiten.Image
	Target   Visual
	sortNext bool
}

func NewCamera() *Camera {
	c := &Camera{
		W:        1,
		H:        1,
		Z:        1,
		sortNext: true,
	}
	c.updateImage()
	return c
}

func (c *Camera) Layout(outsideWidth, outsideHeight int) (int, int) {
	if c.W != float64(outsideWidth) || c.H != float64(outsideHeight) {
		c.W, c.H = float64(outsideWidth), float64(outsideHeight)
		c.updateImage()
	}
	return outsideWidth, outsideHeight
}

func (c *Camera) updateImage() {
	if c.image != nil {
		c.image.Dispose()
	}
	c.image = ebiten.NewImage(int(c.W), int(c.H))
}

func (c *Camera) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		c.Rotate(0.02)
		c.sortNext = true
	} else if ebiten.IsKeyPressed(ebiten.KeyE) {
		c.Rotate(-0.02)
		c.sortNext = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		c.Z -= 0.01
	} else if ebiten.IsKeyPressed(ebiten.KeyX) {
		c.Z += 0.01
	}
	return nil
}

func (c *Camera) Draw(screen *ebiten.Image, visuals []Visual) {
	tx, ty := c.X(), c.Y()
	ox, oy := c.W/2, c.H/2
	if c.Target != nil {
		p := c.Target.Position()
		s := c.Target.Size()
		tx, ty = p.X()+s.X()/2, p.Y()+s.Y()/2
	}

	// Sort visuals by their position with respect to the camera rotation.
	if c.sortNext {
		sort.Slice(visuals, func(i, j int) bool {
			p1, p2 := visuals[i].Position(), visuals[j].Position()
			p1z, p2z := visuals[i].Z(), visuals[j].Z()
			p1.Rotate(-c.Angle())
			p2.Rotate(-c.Angle())
			p1.Add(Vec2{p1z, p1z})
			p2.Add(Vec2{p2z, p2z})
			return p1.Y() < p2.Y()
		})
		c.sortNext = false
	}

	g := ebiten.GeoM{}

	g.Translate(-tx, -ty)
	g.Rotate(-c.Angle())
	g.Translate(ox, oy)

	c.image.Clear()
	for _, v := range visuals {
		v.Draw(c.image, g)
	}
	screen.DrawImage(c.image, nil)
}
