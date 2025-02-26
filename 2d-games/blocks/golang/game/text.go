package game

import (
	"embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
	"github.com/tinne26/etxt/cache"
	"github.com/tinne26/etxt/font"
	"golang.org/x/image/font/sfnt"
)

const (
	RobotoRegularFontName = "Roboto"
	RobotoBoldFontName    = "Roboto Bold"
)

type TextRenderer struct {
	renderer *etxt.Renderer
	font     *sfnt.Font
	color    color.Color
	align    etxt.Align
	size     float64
}

var (
	//go:embed assets/fonts/*
	fonts        embed.FS
	textRenderer *etxt.Renderer
	fontLibrary  *font.Library
)

func init() {
	fontLibrary = font.NewLibrary()
	loaded, skipped, err := fontLibrary.ParseAllFromFS(fonts, "assets/fonts")
	if err != nil {
		log.Fatalf("Error while loading fonts: %s", err.Error())
	}

	fontLibrary.EachFont(func(s string, f *etxt.Font) error {
		log.Printf("%s font is loaded\n", s)

		return nil
	})

	log.Printf("Loaded fonts: %d, Skipped fonts: %d\n", loaded, skipped)

	textRenderer = etxt.NewRenderer()
	glyphsCache := cache.NewDefaultCache(16 * 1024 * 1024) // 16MiB cache
	textRenderer.SetCacheHandler(glyphsCache.NewHandler())
	textRenderer.SetColor(color.RGBA{239, 91, 91, 255})
	textRenderer.SetAlign(etxt.Center)
	textRenderer.SetSize(32)
}

func NewTextRenderer(fontName string, color color.Color, size float64, align etxt.Align) *TextRenderer {
	font := fontLibrary.GetFont(fontName)
	if font == nil {
		log.Fatalf("font '%s' not found\n", fontName)
	}

	r := &TextRenderer{
		renderer: textRenderer,
		font:     font,
		color:    color,
		align:    align,
		size:     size,
	}

	return r
}

func (t *TextRenderer) Draw(target *ebiten.Image, text string, x, y int) {
	t.renderer.SetFont(t.font)
	t.renderer.SetColor(t.color)
	t.renderer.SetSize(t.size)
	t.renderer.SetAlign(t.align)
	t.renderer.Draw(target, text, x, y)
}

func (t *TextRenderer) SetColor(color color.Color) {
	t.color = color
}

func (t *TextRenderer) SetAlign(align etxt.Align) {
	t.align = align
}
