package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

type Game struct {
	text      *TextRenderer
	gameScene *GameScene
}

func NewGame() *Game {
	g := &Game{
		text:      NewTextRenderer(RobotoBoldFontName, color.White, 18, etxt.Center),
		gameScene: newGameScene(),
	}

	return g
}

func (g *Game) GetSize() (screenWidth, screenHeight int) {
	return g.gameScene.GetSize()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.GetSize()
}

func (g *Game) Update() error {
	g.gameScene.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.gameScene.Draw(screen)
}
