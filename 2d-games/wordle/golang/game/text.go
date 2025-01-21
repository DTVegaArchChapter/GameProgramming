package wordle

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func drawText(rt *ebiten.Image, str string, scale float64) {
	drawTextWithShadow(rt, str, ScreenWidth/2, ScreenHeight-195, scale, redColor)
}

func drawTextWithShadow(rt *ebiten.Image, str string, x, y, scale float64, clr color.Color) {
	fontBaseSize := float64(8)

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x)+1, float64(y)+1)
	op.ColorScale.ScaleWithColor(shadowColor)
	op.LineSpacing = fontBaseSize * float64(scale)
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignStart
	text.Draw(rt, str, &text.GoTextFace{
		Source: arcadeTextFaceSource,
		Size:   fontBaseSize * float64(scale),
	}, op)

	op.GeoM.Reset()
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.Reset()
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(rt, str, &text.GoTextFace{
		Source: arcadeTextFaceSource,
		Size:   fontBaseSize * float64(scale),
	}, op)
}
