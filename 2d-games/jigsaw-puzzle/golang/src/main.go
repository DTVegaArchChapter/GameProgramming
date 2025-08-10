package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/DTVegaArchChapter/GameProgramming/jigsaw-puzzle/common"
	"github.com/DTVegaArchChapter/GameProgramming/jigsaw-puzzle/scenes/gameScene"
	"github.com/DTVegaArchChapter/GameProgramming/jigsaw-puzzle/scenes/homeScene"
)

type Game struct {
	SceneManager *common.SceneManager
}

func (g *Game) Update() error {
	return g.SceneManager.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.SceneManager.Draw(screen)
}

func (g *Game) Layout(_, _ int) (int, int) { return common.ScreenWidth, common.ScreenHeight }

func main() {
	ebiten.SetWindowSize(common.ScreenWidth, common.ScreenHeight)
	ebiten.SetWindowTitle("Jigsaw Puzzle")

	gameImage := common.NewGameImage()
	sceneManager := common.NewSceneManager()
	sceneManager.AddScene("Home", func() common.Scene { return homeScene.NewHomeScene(gameImage) })
	sceneManager.AddScene("Game", func() common.Scene { return gameScene.NewGameScene(gameImage) })

	// Set the initial scene to Home
	sceneManager.SetScene("Home")

	game := &Game{
		SceneManager: sceneManager,
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
