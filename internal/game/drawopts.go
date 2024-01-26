package game

import "github.com/hajimehoshi/ebiten/v2"

type DrawOpts struct {
	Image      *ebiten.Image
	GeoM       ebiten.GeoM
	Z          float64
	Angle      float64
	ExtraAngle float64
}

type DrawTargets struct {
	Ground *ebiten.Image
	Shadow *ebiten.Image
	World  *ebiten.Image
	Drops  *ebiten.Image
	Sky    *ebiten.Image
}
