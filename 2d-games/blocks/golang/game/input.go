package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type KeyCode int

const (
	KeyRotate KeyCode = iota
	KeyLeft
	KeyRight
	KeyDown
)

func GetKeyPressed() KeyCode {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		return KeyRotate
	} else if d := inpututil.KeyPressDuration(ebiten.KeyArrowLeft); isKeyPressDurationValid(d) {
		return KeyLeft
	} else if d := inpututil.KeyPressDuration(ebiten.KeyArrowRight); isKeyPressDurationValid(d) {
		return KeyRight
	} else if d := inpututil.KeyPressDuration(ebiten.KeyArrowDown); isKeyPressDurationValid(d) {
		return KeyDown
	}

	return -1
}

func isKeyPressDurationValid(d int) bool {
	return d > 0 && (d == 1 || (d >= 10 && d%2 == 0))
}
