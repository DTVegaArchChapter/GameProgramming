package game

import (
	"image"
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"
)

type GameScene struct {
	playField       *PlayField
	currentPiece    *Piece
	nextPiece       *Piece
	gameOver        bool
	gameOverImage   *ebiten.Image
	moveDownCounter *TicksCounter
	text            *TextRenderer
	nextPieceRect   image.Rectangle
	score           int
	lines           int
	level           int
}

func newGameScene() *GameScene {
	g := &GameScene{
		playField:       newPlayField(20, 20, 10, 20, 25),
		gameOver:        false,
		moveDownCounter: NewTicksCounter(ebiten.TPS()),
		text:            NewTextRenderer(RobotoBoldFontName, color.Black, 20, etxt.Center),
		score:           0,
		lines:           0,
		level:           0,
	}

	playFieldW, _ := g.playField.GetSize()
	nextPieceX, nextPieceY := playFieldW+g.playField.x*2, g.playField.y+25
	g.nextPieceRect = image.Rect(nextPieceX, nextPieceY, nextPieceX+g.playField.tileSize*6, nextPieceY+g.playField.tileSize*6)

	w, h := g.GetSize()
	g.gameOverImage = ebiten.NewImage(w, h)
	g.gameOverImage.Fill(color.RGBA{0, 0, 0, 192})
	g.text.SetColor(color.Opaque)
	g.text.Draw(g.gameOverImage, "GAME OVER\nPRESS ANY KEY TO RESTART", g.gameOverImage.Bounds().Dx()/2, g.gameOverImage.Bounds().Dy()/2)

	g.setNewPiece()

	return g
}

func (g *GameScene) setNewPiece() bool {
	if g.currentPiece == nil {
		g.currentPiece = createNewPiece(g.playField)
	} else {
		g.currentPiece = g.nextPiece
	}

	g.nextPiece = createNewPiece(g.playField)

	return !g.currentPiece.collides()
}

func (g *GameScene) GetSize() (screenWidth, screenHeight int) {
	w, h := g.playField.GetSize()
	return w + g.playField.x*3 + g.nextPieceRect.Dx(), h + g.playField.y*2
}

func (g *GameScene) Update() error {
	if g.gameOver {
		return nil
	}

	switch GetKeyPressed() {
	case KeyRotate:
		g.currentPiece.Turn()
	case KeyLeft:
		g.currentPiece.MoveLeft()
	case KeyRight:
		g.currentPiece.MoveRight()
	case KeyDown:
		g.currentPiece.MoveDown()
	}

	canMoveDown := true
	if g.moveDownCounter.Update() {
		canMoveDown = g.currentPiece.MoveDown()
	}

	if !canMoveDown {
		g.currentPiece.AbsorbIntoPlayField()
		if l := g.playField.ClearLines(); l > 0 {
			f := func(n int) int {
				result := 1
				for i := 2; i <= n; i++ {
					result *= i
				}
				return result
			}
			g.score += 50 * f(l) * (g.level + 1)
			g.lines += l
			g.level = g.lines / 10

			tps := float64(ebiten.TPS())
			g.moveDownCounter.SetTicks(int(math.Max(1.0, 4.0*tps/float64(g.level+4))))
		}
		if !g.setNewPiece() {
			g.gameOver = true
		}
	}

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 225, G: 225, B: 225, A: 255})
	g.playField.Draw(screen)

	if g.currentPiece != nil {
		g.currentPiece.Draw(screen)
	}

	g.text.SetAlign(etxt.Top | etxt.Left)
	g.text.SetColor(color.RGBA{R: 170, G: 50, B: 50, A: 255})
	g.text.Draw(screen, "NEXT", g.nextPieceRect.Min.X, g.nextPieceRect.Min.Y-25)
	vector.DrawFilledRect(screen, float32(g.nextPieceRect.Min.X), float32(g.nextPieceRect.Min.Y), float32(g.nextPieceRect.Dx()), float32(g.nextPieceRect.Dy()), g.playField.emptyColor, false)

	if g.nextPiece != nil {
		rect := g.nextPiece.getRectangle()
		x := float32(6-rect.Dx()) / 2
		y := float32(6-rect.Dy()) / 2

		for _, p := range *g.nextPiece.blocks {
			vector.DrawFilledRect(screen, float32(g.nextPieceRect.Min.X)+((float32(p.X)+x)*float32(g.playField.tileSize)), float32(g.nextPieceRect.Min.Y)+((float32(p.Y)+y)*float32(g.playField.tileSize)), float32(g.playField.tileSize), float32(g.playField.tileSize), g.nextPiece.color, false)
		}
	}

	g.text.SetColor(color.RGBA{R: 170, G: 50, B: 50, A: 255})
	g.text.Draw(screen, "SCORE", g.nextPieceRect.Min.X, g.nextPieceRect.Min.Y+g.nextPieceRect.Dy()+15)
	vector.DrawFilledRect(screen, float32(g.nextPieceRect.Min.X), float32(g.nextPieceRect.Min.Y+g.nextPieceRect.Dy()+40), float32(g.nextPieceRect.Dx()), 35, g.playField.emptyColor, false)

	g.text.SetColor(color.White)
	g.text.SetAlign(etxt.Right)
	g.text.Draw(screen, strconv.Itoa(g.score), g.nextPieceRect.Min.X+g.nextPieceRect.Dx()-5, g.nextPieceRect.Min.Y+g.nextPieceRect.Dy()+40+7)

	g.text.SetColor(color.RGBA{R: 170, G: 50, B: 50, A: 255})
	g.text.SetAlign(etxt.Top | etxt.Left)
	g.text.Draw(screen, "LEVEL", g.nextPieceRect.Min.X, g.nextPieceRect.Min.Y+g.nextPieceRect.Dy()+90)
	vector.DrawFilledRect(screen, float32(g.nextPieceRect.Min.X), float32(g.nextPieceRect.Min.Y+g.nextPieceRect.Dy()+115), float32(g.nextPieceRect.Dx()), 35, g.playField.emptyColor, false)

	g.text.SetColor(color.White)
	g.text.SetAlign(etxt.Right)
	g.text.Draw(screen, strconv.Itoa(g.level), g.nextPieceRect.Min.X+g.nextPieceRect.Dx()-5, g.nextPieceRect.Min.Y+g.nextPieceRect.Dy()+115+7)

	g.text.SetColor(color.RGBA{R: 170, G: 50, B: 50, A: 255})
	g.text.SetAlign(etxt.Top | etxt.Left)
	g.text.Draw(screen, "LINES", g.nextPieceRect.Min.X, g.nextPieceRect.Min.Y+g.nextPieceRect.Dy()+165)
	vector.DrawFilledRect(screen, float32(g.nextPieceRect.Min.X), float32(g.nextPieceRect.Min.Y+g.nextPieceRect.Dy()+190), float32(g.nextPieceRect.Dx()), 35, g.playField.emptyColor, false)

	g.text.SetColor(color.White)
	g.text.SetAlign(etxt.Right)
	g.text.Draw(screen, strconv.Itoa(g.lines), g.nextPieceRect.Min.X+g.nextPieceRect.Dx()-5, g.nextPieceRect.Min.Y+g.nextPieceRect.Dy()+190+7)

	if g.gameOver {
		screen.DrawImage(g.gameOverImage, nil)
	}
}
