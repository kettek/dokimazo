package game

import (
	"math"

	"github.com/kettek/dokimazo/internal/res"
)

type Drops []*Drop

func (d *Drops) Add(drop *Drop) {
	*d = append(*d, drop)
}

func (d *Drops) Remove(drop *Drop) {
	for i, drop2 := range *d {
		if drop2 == drop {
			*d = append((*d)[:i], (*d)[i+1:]...)
			return
		}
	}
}

type Drop struct {
	x, y     float64 // starting x/y.... ugh
	RID      res.RID
	State    interface{} // For now...
	drop     *res.Drop
	lifetime int
	*SpriteStack
}

func (d *Drop) Update() []Request {
	d.lifetime++
	d.x = math.Cos(float64(d.lifetime) / 40.0)
	d.y = math.Sin(float64(d.lifetime) / 10.0)
	// nada
	return nil
}

func (d *Drop) Draw(drawOpts DrawOpts) {
	drawOpts.GeoM.Translate(d.x*drawOpts.Z, (d.y-8.0)*drawOpts.Z)
	d.SpriteStack.Draw(drawOpts)
}

func (d *Drop) DrawShadow(drawOpts DrawOpts) {
	drawOpts.GeoM.Translate(d.x*drawOpts.Z, 0)
	d.SpriteStack.DrawShadow(drawOpts)
}
