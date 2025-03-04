package main

import (
	"bytes"
	"image"
	"image/png"
	"log"

	wordle "github.com/DTVegaArchChapter/GameProgramming/wordle/game"
	"github.com/hajimehoshi/ebiten/v2"

	_ "embed"
)

//go:embed game/assets/icon.png
var iconData []byte

func main() {
	// Decode the embedded PNG data
	icon, err := png.Decode(bytes.NewReader(iconData))
	if err != nil {
		log.Fatal(err)
	}

	game := wordle.NewGame()
	ebiten.SetWindowSize(wordle.ScreenWidth, wordle.ScreenHeight)
	ebiten.SetWindowTitle("Türkçe Wordle")
	ebiten.SetWindowIcon([]image.Image{icon})

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
