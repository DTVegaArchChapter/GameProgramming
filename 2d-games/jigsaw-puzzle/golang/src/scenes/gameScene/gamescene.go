package gameScene

import (
	"fmt"
	"image/color"
	"time"

	"github.com/DTVegaArchChapter/GameProgramming/jigsaw-puzzle/common"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"
)

type GameScene struct {
	PuzzlePicture     *PuzzlePicture
	headerHeight      float32
	footerHeight      float32
	pictureName       string
	moves             int
	headerButtons     []*common.Button
	footerButtons     []*common.Button
	startTime         int64
	endTime           int64
	isPuzzleCompleted bool
	showGhost         bool
	showImage         bool

	// FPS diagram
	showFPSDiagram bool
	fpsHistory     []float64

	text *common.TextRenderer
}

func NewGameScene(gameImage *common.GameImage) *GameScene {
	text := common.NewTextRenderer(common.RobotoBoldFontName, common.BodyTextColor, 40, etxt.Center)
	buttonOptions := []common.ButtonOptFunc{
		common.ButtonOption.WithColor(common.HeaderButtonColor),
		common.ButtonOption.WithHoverColor(common.HeaderButtonHoverColor),
		common.ButtonOption.WithFontColor(common.BodyTextColor),
		common.ButtonOption.WithFontSize(20),
	}

	s := &GameScene{
		text:              text,
		headerHeight:      64,
		footerHeight:      60,
		pictureName:       gameImage.GetName(),
		startTime:         time.Now().Unix(),
		endTime:           0,
		moves:             0,
		showGhost:         false,
		showImage:         false,
		isPuzzleCompleted: false,
		headerButtons: []*common.Button{
			common.NewButton(
				1060, 12,
				55, 40,
				"Restart",
				buttonOptions...,
			),
			common.NewButton(
				1165, 12,
				55, 40,
				"Home",
				buttonOptions...,
			),
		},
		footerButtons: []*common.Button{
			common.NewButton(
				20, float32(common.ScreenHeight)-50,
				50, 40,
				"Image",
				buttonOptions...,
			),
			common.NewButton(
				120, float32(common.ScreenHeight)-50,
				50, 40,
				"Ghost",
				buttonOptions...,
			),
		},
	}
	s.PuzzlePicture = CreatePuzzlePicture(gameImage.GetImage())
	s.PuzzlePicture.CreatePuzzlePieces(48)

	return s
}

func (g *GameScene) Update(context *common.SceneContext) error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()

		g.PuzzlePicture.SetPieceBeingDragged(float64(mx), float64(my))
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.PuzzlePicture.HandleDraggedPieceSnapping()
		if g.PuzzlePicture.DropPuzzlePieces() {
			g.incrementMoves()
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		g.PuzzlePicture.MoveDraggedPiece(float64(mx), float64(my))
	}

	for _, button := range g.headerButtons {
		button.Update()
		if button.Clicked {
			switch button.Label {
			case "Home":
				context.SceneManager.SetScene("Home")
			case "Restart":
				context.SceneManager.SetScene("Game")
			}
		}
	}

	for _, button := range g.footerButtons {
		button.Update()

		if button.Clicked {
			switch button.Label {
			case "Image":
				g.showImage = !g.showImage
			case "Ghost":
				g.showGhost = !g.showGhost
			}
		}
	}

	if !g.isPuzzleCompleted {
		g.isPuzzleCompleted = g.PuzzlePicture.IsPuzzleCompleted()
		g.endTime = time.Now().Unix()
	}

	// Toggle FPS diagram with 'D'
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.showFPSDiagram = !g.showFPSDiagram
	}

	// Collect FPS samples when visible
	if g.showFPSDiagram {
		g.fpsHistory = append(g.fpsHistory, ebiten.ActualTPS())
		const maxSamples = 180
		if len(g.fpsHistory) > maxSamples {
			g.fpsHistory = g.fpsHistory[len(g.fpsHistory)-maxSamples:]
		}
	}

	return nil
}

func (g *GameScene) incrementMoves() {
	if !g.isPuzzleCompleted {
		g.moves++
	}
}

func (g *GameScene) Draw(screen *ebiten.Image, context *common.SceneContext) {
	screen.Fill(common.BackgroundColor)

	if g.showGhost {
		opt := &ebiten.DrawImageOptions{}

		opt.ColorScale.Scale(0.5, 0.5, 0.5, 1) // make ghost image semi-transparent

		// translate the ghost image to center
		opt.GeoM.Translate(
			(float64(common.ScreenWidth)-float64(g.PuzzlePicture.image.Bounds().Dx()))/2,
			(float64(common.ScreenHeight)-float64(g.PuzzlePicture.image.Bounds().Dy()))/2,
		)
		screen.DrawImage(g.PuzzlePicture.image, opt)
	}

	g.PuzzlePicture.Draw(screen)

	g.drawHeader(screen)
	g.drawFooter(screen)

	if g.showFPSDiagram {
		g.drawFPSDiagram(screen)
	}
}

func (g *GameScene) drawHeader(screen *ebiten.Image) {
	// Draw header background
	vector.DrawFilledRect(screen, 0, 0, float32(common.ScreenWidth), g.headerHeight, common.HeaderColor, true)

	// Draw title text
	g.text.SetColor(common.TitleColor)
	g.text.SetSize(36)
	g.text.SetAlign(etxt.Left)
	g.text.Draw(screen, fmt.Sprintf("Jigsaw Puzzle - %s", g.pictureName), 32, 32)

	// Draw buttons
	for _, button := range g.headerButtons {
		button.Draw(screen)
	}
}

func (g *GameScene) drawFooter(screen *ebiten.Image) {
	// Draw footer background
	vector.DrawFilledRect(screen, 0, float32(common.ScreenHeight)-g.footerHeight, float32(common.ScreenWidth), g.footerHeight, common.FooterColor, true)

	// Draw footer text
	g.text.SetColor(common.BodyTextColor)
	g.text.SetSize(24)
	g.text.SetAlign(etxt.Right)

	g.text.Draw(
		screen,
		fmt.Sprintf("Moves: %d - Time: %s", g.moves, g.getElapsedTime()),
		int(float32(common.ScreenWidth)-20),
		int(float32(common.ScreenHeight)-g.footerHeight/2),
	)

	// Draw footer buttons
	for _, button := range g.footerButtons {
		button.Draw(screen)
	}

	if g.isPuzzleCompleted {
		g.text.SetColor(common.TitleColor)
		g.text.SetSize(26)
		g.text.SetAlign(etxt.Center)
		g.text.DrawHorizontalCenter(screen, "Puzzle Completed!", int(float32(common.ScreenHeight)-g.footerHeight/2))
	} else {
		g.text.SetColor(common.TitleColor)
		g.text.SetSize(26)
		g.text.SetAlign(etxt.Center)
		g.text.DrawHorizontalCenter(screen, fmt.Sprintf("%%%d Completed", g.PuzzlePicture.GetCompletePercentage()), int(float32(common.ScreenHeight)-g.footerHeight/2))
	}

	if g.showImage {
		scale := 0.5
		margin := 10.
		opt := &ebiten.DrawImageOptions{}

		opt.GeoM.Scale(scale, scale)
		opt.GeoM.Translate(
			margin,
			float64(common.ScreenHeight)-float64(g.PuzzlePicture.image.Bounds().Dy())*scale-float64(g.footerHeight)-margin,
		)

		screen.DrawImage(g.PuzzlePicture.image, opt)
	}
}

func (g *GameScene) getElapsedTime() string {
	var elapsed int64
	if g.isPuzzleCompleted {
		elapsed = g.endTime - g.startTime
	} else {
		// If puzzle is not completed, calculate elapsed time from start time
		elapsed = time.Now().Unix() - g.startTime
	}

	days := elapsed / (24 * 3600)
	hours := (elapsed % (24 * 3600)) / 3600
	minutes := (elapsed % 3600) / 60
	seconds := elapsed % 60

	var timeStr string
	if days > 0 {
		timeStr = fmt.Sprintf("%d.%02d:%02d:%02d", days, hours, minutes, seconds)
	} else if hours > 0 {
		timeStr = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	} else {
		timeStr = fmt.Sprintf("%02d:%02d", minutes, seconds)
	}
	return timeStr
}

func (g *GameScene) drawFPSDiagram(screen *ebiten.Image) {
	if len(g.fpsHistory) == 0 {
		return
	}

	const (
		w          = 200.0
		h          = 70.0
		marginX    = 10.0
		marginY    = 10.0
		maxFPSBase = 120.0
	)

	// Background panel
	vector.DrawFilledRect(
		screen,
		float32(marginX), float32(marginY),
		float32(w), float32(h),
		color.RGBA{0, 0, 0, 160}, // semi-transparent black
		true,
	)

	// Axes baseline (optional thin line at bottom)
	vector.DrawFilledRect(
		screen,
		float32(marginX), float32(marginY+h-1),
		float32(w), 1,
		color.RGBA{200, 200, 200, 160}, // light gray
		true,
	)

	n := len(g.fpsHistory)
	if n < 2 {
		return
	}

	barW := w / float64(len(g.fpsHistory))
	x := marginX

	for _, fps := range g.fpsHistory {
		if fps > maxFPSBase {
			fps = maxFPSBase
		}
		ratio := fps / maxFPSBase
		barH := ratio * (h - 15) // leave room for text
		y := marginY + (h - 1) - barH

		barColor := color.RGBA{80, 220, 120, 200} // #50FF7F
		if fps < 50 {
			barColor = color.RGBA{220, 120, 80, 200} // #DC784C
		}

		vector.DrawFilledRect(
			screen,
			float32(x), float32(y),
			float32(barW-1), float32(barH),
			barColor,
			true,
		)
		x += barW
	}

	// FPS label
	g.text.SetColor(common.BodyTextColor)
	g.text.SetSize(16)
	g.text.SetAlign(etxt.Left)
	current := g.fpsHistory[len(g.fpsHistory)-1]
	g.text.Draw(screen, fmt.Sprintf("TPS: %.1f", current), int(marginX+6), int(marginY+14))
}
