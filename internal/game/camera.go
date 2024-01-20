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

func (c *Camera) Draw(screen *ebiten.Image, visuals Visuals) {
	t := Vec2{c.X(), c.Y()}
	o := Vec2{c.W / 2, c.H / 2}
	if c.Target != nil {
		p := c.Target.Position()
		s := c.Target.Size()
		t = Vec2{p.X() + s.X()/2, p.Y() + s.Y()/2}
	}

	// Sort visuals by their position with respect to the camera rotation.
	sort.Slice(visuals, func(i, j int) bool {
		p1, p2 := visuals[i].Position(), visuals[j].Position()
		p1.RotateAround(t, -c.Angle())
		p2.RotateAround(t, -c.Angle())
		return p1.Y()+visuals[i].Z() < p2.Y()+visuals[j].Z()
	})

	g := ebiten.GeoM{}

	g.Translate(-t.X(), -t.Y())
	g.Rotate(-c.Angle())
	g.Scale(c.Z, c.Z)
	g.Translate(o.X(), o.Y())

	c.image.Clear()
	drawOpts := DrawOpts{
		Image: screen,
		GeoM:  g,
		Z:     c.Z,
		Angle: c.Angle(),
	}
	for _, v := range visuals {
		v.Draw(drawOpts)
	}
	screen.DrawImage(c.image, nil)
}
