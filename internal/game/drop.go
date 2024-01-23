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
	Visual
}

func (d *Drop) Update() []Request {
	d.lifetime++
	d.x = math.Cos(float64(d.lifetime) / 50.0)
	d.y = math.Sin(float64(d.lifetime) / 80.0)
	// nada
	return nil
}

func (d *Drop) Draw(drawOpts DrawOpts) {
	drawOpts.GeoM.Translate(d.x, d.y-8.0)
	d.Visual.Draw(drawOpts)
}

func (d *Drop) DrawShadow(drawOpts DrawOpts) {
	drawOpts.GeoM.Translate(d.x, 0)
	if v, ok := d.Visual.(VisualShadow); ok {
		v.DrawShadow(drawOpts)
	}
}
