package wordle

import (
	"image/color"

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
	text     *TextRenderer
}

func NewGame() *Game {
	k := newKeyboard()
	g := &Game{
		keyboard: k,
		board:    newBoard(k),
		text:     NewTextRenderer(RobotoBoldFontName, redColor, 18),
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
		g.setMessage(screen, "Not in Word List!!", redColor)
	}

	if g.board.state == gameLost {
		g.setMessage(screen, g.board.GetCorrectAnswer(), redColor)
	} else if g.board.state == gameWon && g.board.IsWinAnimationFinished() {
		g.setMessage(screen, "You Won!!", greenColor)
	}
}

func (g *Game) setMessage(screen *ebiten.Image, messageText string, color color.Color) {
	g.text.SetColor(color)
	g.text.Draw(screen, messageText, screen.Bounds().Dx()/2, int(g.board.maxY)+30)
}
