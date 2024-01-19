package main

import (
	"rotate-test/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.New()

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
