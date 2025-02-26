package main

import (
	"log"

	"github.com/DTVegaArchChapter/GameProgramming/blocks/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := game.NewGame()
	w, h := game.GetSize()
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("Blocks")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
