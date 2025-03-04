package wordle

import (
	"embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
	"golang.org/x/image/font/sfnt"
)

type TextRenderer struct {
	renderer *etxt.Renderer
	font     *sfnt.Font
	color    color.Color
	sizePx   int
}

var (
	//go:embed assets/fonts/*
	fonts        embed.FS
	textRenderer *etxt.Renderer
	fontLibrary  *etxt.FontLibrary
)

func init() {
	fontLibrary = etxt.NewFontLibrary()
	loaded, skipped, err := fontLibrary.ParseEmbedDirFonts("assets/fonts", fonts)
	if err != nil {
		log.Fatalf("Error while loading fonts: %s", err.Error())
	}

	fontLibrary.EachFont(func(s string, f *etxt.Font) error {
		log.Printf("%s font is loaded\n", s)

		return nil
	})

	log.Printf("Loaded fonts: %d, Skipped fonts: %d\n", loaded, skipped)

	textRenderer = etxt.NewStdRenderer()
	glyphsCache := etxt.NewDefaultCache(10 * 1024 * 1024) // 10MB
	textRenderer.SetCacheHandler(glyphsCache.NewHandler())
	textRenderer.SetAlign(etxt.YCenter, etxt.XCenter)
	textRenderer.SetSizePx(32)
}

func NewTextRenderer(fontName string, color color.Color, sizePx int) *TextRenderer {
	font := fontLibrary.GetFont(fontName)
	if font == nil {
		log.Fatalf("font '%s' not found\n", fontName)
	}

	return &TextRenderer{
		renderer: textRenderer,
		font:     font,
		color:    color,
		sizePx:   sizePx,
	}
}

func (t *TextRenderer) Draw(target *ebiten.Image, text string, x, y int) {
	t.renderer.SetTarget(target)
	t.renderer.SetFont(t.font)
	t.renderer.SetColor(t.color)
	t.renderer.SetSizePx(t.sizePx)
	t.renderer.Draw(text, x, y)
}

func (t *TextRenderer) SetColor(color color.Color) {
	t.color = color
}
