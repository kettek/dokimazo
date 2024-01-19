package main

import (
	"rotate-test/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.New()

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
