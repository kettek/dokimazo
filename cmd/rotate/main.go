package main

import (
	"github.com/kettek/dokimazo/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.New()

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("dokimazo")

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
