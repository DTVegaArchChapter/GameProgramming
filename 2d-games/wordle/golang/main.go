package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/DTVegaArchChapter/GameProgramming/wordle/game"
)

func main() {
	game := wordle.NewGame()
	ebiten.SetWindowSize(wordle.ScreenWidth, wordle.ScreenHeight)
	ebiten.SetWindowTitle("Türkçe Wordle")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
