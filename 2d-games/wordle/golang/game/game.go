package wordle

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 400
	ScreenHeight = 600
)

type Game struct {
	board    *board
	keyboard *keyboard
	keys     []ebiten.Key
	runes    []rune
}

func NewGame() *Game {
	k := newKeyboard()
	g := &Game{
		keyboard: k,
		board:    newBoard(k),
	}

	return g
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	g.runes = ebiten.AppendInputChars(g.runes[:0])
	if len(g.runes) > 0 {
		g.board.inputRune = g.runes[0]
	} else {
		g.board.inputRune = 0
	}

	g.keys = inpututil.AppendJustPressedKeys(g.keys[:0])
	if len(g.keys) > 0 {
		g.board.inputKey = g.keys[0]
	} else {
		g.board.inputKey = 0
	}

	g.board.Update()
	if g.board.state != gameInProgress {
		return nil
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	g.board.Draw(screen)
	g.keyboard.draw(screen)

	if g.board.wordNotFound {
		drawText(screen, "Not in word list", 1.75)
	}

	if g.board.state != gameInProgress {
		drawText(screen, TurkishUpper.String(string(g.board.answer)), 2)
	}
}
