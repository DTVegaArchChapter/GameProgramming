package homeScene

import (
	"image/color"
	"image/jpeg"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sqweek/dialog"
	"github.com/tinne26/etxt"

	"github.com/DTVegaArchChapter/GameProgramming/jigsaw-puzzle/common"
)

// Subtitle - soft orange for friendly feel
var SubtitleColor = color.RGBA{R: 255, G: 165, B: 0, A: 255} // #FFA500

// Button background - vivid sky blue
var ButtonBgColor = color.RGBA{R: 0, G: 153, B: 255, A: 255} // #0099FF

// Button hover - lighter cyan
var ButtonHoverColor = color.RGBA{R: 51, G: 204, B: 255, A: 255} // #33CCFF

// Button text - pure white
var ButtonTextColor = color.RGBA{R: 255, G: 255, B: 255, A: 255} // #FFFFFF

// Button shadow - deep navy
var ButtonShadowColor = color.RGBA{R: 0, G: 0, B: 0, A: 200} // #000000C8

type imageWithName struct {
	image *ebiten.Image
	name  string
}

type puzzleImage struct {
	image   *ebiten.Image
	name    string
	x, y    float64
	scale   float64
	hovered bool
}

type HomeScene struct {
	images       []*puzzleImage
	uploadButton *common.Button
	gameImage    *common.GameImage
	text         *common.TextRenderer
}

func NewHomeScene(gameImage *common.GameImage) *HomeScene {
	images, err := loadImages("./pictures")
	if err != nil {
		panic(err)
	}

	var puzzleImages []*puzzleImage
	var x, y float64
	top := 150.
	left := 20.
	cols := 5
	for i, img := range images {
		x = left + 250*float64(i%cols)
		y = top + 170*float64(i/cols)

		puzzleImages = append(puzzleImages, &puzzleImage{image: img.image, name: img.name, x: x, y: y, scale: 0.375})
	}

	text := common.NewTextRenderer(common.RobotoBoldFontName, common.BodyTextColor, 40, etxt.Center)

	return &HomeScene{
		images:    puzzleImages,
		gameImage: gameImage,
		uploadButton: common.NewButton(
			(float32(common.ScreenWidth)-250)/2, 600,
			200, 50,
			"Upload Image",
			common.ButtonOption.WithFontSize(24),
			common.ButtonOption.WithFontColor(ButtonTextColor),
			common.ButtonOption.WithHoverColor(ButtonHoverColor),
			common.ButtonOption.WithColor(ButtonBgColor),
			common.ButtonOption.WithShadowColor(ButtonShadowColor),
		),
		text: text,
	}
}

func (h *HomeScene) Update(context *common.SceneContext) error {
	mx, my := ebiten.CursorPosition()

	h.uploadButton.Update()

	if h.uploadButton.Clicked {
		name, img, err := loadImageFromDesktop()
		if err != nil && err != dialog.Cancelled {
			return err
		}

		if img != nil {
			h.gameImage.SetImage(name, img)
			context.SceneManager.SetScene("Game")
			return nil
		}
	}

	for _, img := range h.images {
		if isPointInImage(float64(mx), float64(my), img) {
			img.hovered = true
		} else {
			img.hovered = false
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for _, img := range h.images {
			if isPointInImage(float64(mx), float64(my), img) {
				h.gameImage.SetImage(img.name, img.image)
				context.SceneManager.SetScene("Game")
			}
		}
	}

	return nil
}

func (h *HomeScene) Draw(screen *ebiten.Image, context *common.SceneContext) {
	screen.Fill(common.BackgroundColor)

	h.text.SetColor(common.TitleColor)
	h.text.SetSize(36)
	h.text.DrawHorizontalCenter(screen, "Welcome to Jigsaw Puzzle!", 50)

	h.text.SetColor(SubtitleColor)
	h.text.SetSize(24)
	h.text.DrawHorizontalCenter(screen, "Select a puzzle image", 115)

	h.text.SetColor(common.BodyTextColor)
	h.text.SetSize(24)
	h.text.DrawHorizontalCenter(screen, "OR", 540)

	isHovered := false
	for _, img := range h.images {
		scale := img.scale
		if img.hovered {
			scale *= 1.0125
			isHovered = true
		}

		geoM := ebiten.GeoM{}
		geoM.Scale(scale, scale)
		geoM.Translate(img.x, img.y)
		opt := &ebiten.DrawImageOptions{GeoM: geoM}
		screen.DrawImage(img.image, opt)
	}

	if isHovered {
		ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}

	h.uploadButton.Draw(screen)
}

func isPointInImage(x, y float64, puzzleImage *puzzleImage) bool {
	geoM := ebiten.GeoM{}
	geoM.Scale(puzzleImage.scale, puzzleImage.scale)
	geoM.Translate(puzzleImage.x, puzzleImage.y)

	if !geoM.IsInvertible() {
		return false
	}

	geoM.Invert()

	imgX, imgY := geoM.Apply(x, y)

	img := puzzleImage.image
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	return imgX >= 0 && imgX < float64(w) && imgY >= 0 && imgY < float64(h)
}

func loadImageFromDesktop() (string, *ebiten.Image, error) {
	path, err := dialog.File().Filter("Image files", "jpg", "jpeg").Load()
	if err != nil {
		return "", nil, err
	}

	img, err := loadJpegImageFromPath(path)

	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), img, err
}

func loadImages(path string) ([]*imageWithName, error) {
	var images []*imageWithName

	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if ext := filepath.Ext(p); ext == ".jpg" || ext == ".jpeg" {
				img, err := loadJpegImageFromPath(p)
				if err != nil {
					return err
				}

				// get the file name without the path and extension
				name := strings.TrimSuffix(filepath.Base(p), ext)

				images = append(images, &imageWithName{
					image: img,
					name:  name,
				})
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return images, nil
}

func loadJpegImageFromPath(path string) (*ebiten.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
}
