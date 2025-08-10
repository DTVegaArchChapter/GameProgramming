package common

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"
)

type ButtonOptFunc func(*Button)

type ButtonOptionBuilder struct{}

var ButtonOption = ButtonOptionBuilder{}

func (ButtonOptionBuilder) WithFontSize(fontSize float64) ButtonOptFunc {
	return func(b *Button) {
		b.FontSize = fontSize
	}
}

func (ButtonOptionBuilder) WithFontColor(fontColor color.Color) ButtonOptFunc {
	return func(b *Button) {
		b.FontColor = fontColor
	}
}

func (ButtonOptionBuilder) WithFontName(fontName string) ButtonOptFunc {
	return func(b *Button) {
		b.FontName = fontName
	}
}

func (ButtonOptionBuilder) WithColor(color color.RGBA) ButtonOptFunc {
	return func(b *Button) {
		b.Color = color
	}
}

func (ButtonOptionBuilder) WithHoverColor(hoverColor color.RGBA) ButtonOptFunc {
	return func(b *Button) {
		b.HoverColor = hoverColor
	}
}

func (ButtonOptionBuilder) WithShadowColor(shadowColor color.RGBA) ButtonOptFunc {
	return func(b *Button) {
		b.ShadowColor = shadowColor
	}
}

type Button struct {
	X           float32
	Y           float32
	Width       float32
	Height      float32
	Label       string
	Color       color.RGBA
	HoverColor  color.RGBA
	ShadowColor color.RGBA
	Pressed     bool
	Hovered     bool
	Clicked     bool
	OnClick     func()
	FontSize    float64
	FontColor   color.Color
	FontName    string
	text        *TextRenderer
	path        *vector.Path
}

func NewButton(x, y, width, height float32, label string, opts ...ButtonOptFunc) *Button {
	fontSize := 16.0
	fontColor := color.Black
	fontName := RobotoBoldFontName

	r := height / 2
	path := vector.Path{}
	path.MoveTo(r, height)
	path.Arc(r, height/2, r, -math.Pi*1.5, -math.Pi*0.5, vector.Clockwise)
	path.LineTo(r, 0)
	path.LineTo(width+r, 0)
	path.Arc(width+r, height/2, r, -math.Pi*0.5, math.Pi*0.5, vector.Clockwise)
	path.LineTo(width+r, height)
	path.Close()

	textRenderer := NewTextRenderer(fontName, fontColor, fontSize, etxt.Center)
	btn := &Button{
		X:           x,
		Y:           y,
		Width:       width,
		Height:      height,
		Label:       label,
		Color:       color.RGBA{R: 54, G: 153, B: 255, A: 255}, // bluish
		HoverColor:  color.RGBA{R: 72, G: 176, B: 255, A: 255}, // lighter bluish
		ShadowColor: color.RGBA{R: 0, G: 0, B: 0, A: 50},       // semi-transparent black
		FontName:    fontName,
		FontSize:    fontSize,
		FontColor:   fontColor,
		text:        textRenderer,
		path:        &path,
	}

	for _, opt := range opts {
		opt(btn)
	}

	return btn
}

func (b *Button) Draw(screen *ebiten.Image) {
	bColor := b.Color
	if b.Hovered {
		bColor = b.HoverColor
	}

	if b.Pressed {
		// darker when pressed
		bColor = color.RGBA{R: uint8(float32(bColor.R) * 0.85), G: uint8(float32(bColor.G) * 0.85), B: uint8(float32(bColor.B) * 0.85), A: bColor.A}
	}

	b.drawShadow(screen)
	b.drawBackground(screen, bColor, 0, 0)
	b.drawText(screen)
}

func (b *Button) Update() {
	b.Clicked = false
	b.Hovered = b.isHovered()

	// detect click on release inside bounds
	pressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if pressed && b.Hovered {
		b.Pressed = true
	}

	if b.Pressed && !pressed {
		// mouse released; if still hovering, treat as click
		if b.Hovered {
			b.Clicked = true
			if b.OnClick != nil {
				b.OnClick()
			}
		}
		b.Pressed = false
	}
}

func (b *Button) isHovered() bool {
	mx, my := ebiten.CursorPosition()

	vertices := GetPathVertices(b.path)

	for i := range vertices {
		vertices[i].X += float64(b.X)
		vertices[i].Y += float64(b.Y)
	}

	// Check if mouse is within button bounds using path
	return IsPointInPathVertices(Point{X: float64(mx), Y: float64(my)}, vertices)
}

func (b *Button) drawShadow(screen *ebiten.Image) {
	maxOffset := 6
	baseAlpha := 60.0
	for i := 1; i <= maxOffset; i++ {
		alpha := uint8(baseAlpha / (1.0 + float64(i)/float64(maxOffset))) // fade out as it goes further
		offset := float32(i)

		shadowColor := color.RGBA{b.ShadowColor.R, b.ShadowColor.G, b.ShadowColor.B, alpha}

		// Draw the shadow with an offset
		StrokePathWithColor(screen, b.path, b.X+offset, b.Y+offset, 2, shadowColor)
	}
}

func (b *Button) drawBackground(screen *ebiten.Image, bColor color.Color, offsetX, offsetY float32) {
	FillPathWithColor(screen, b.path, b.X+offsetX, b.Y+offsetY, bColor)
}

func (b *Button) drawText(screen *ebiten.Image) {
	// Draw the label text in the center of the button
	textX, textY := b.getCenter()

	b.text.SetFont(b.FontName)
	b.text.SetSize(b.FontSize)
	b.text.SetColor(b.FontColor)
	b.text.SetAlign(etxt.Center)
	b.text.Draw(screen, b.Label, textX, textY)
}

func (b *Button) getCenter() (int, int) {
	r := int(b.Height / 2)
	return int(b.X) + int(b.Width/2) + r, int(b.Y) + int(b.Height/2)
}
