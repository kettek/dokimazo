package game

import "github.com/hajimehoshi/ebiten/v2"

type DrawOpts struct {
	Image *ebiten.Image
	GeoM  ebiten.GeoM
	Z     float64
	Angle float64
}
