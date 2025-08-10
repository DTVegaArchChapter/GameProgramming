package common

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawShadowForPath(screen *ebiten.Image, x, y float64, path *vector.Path) {
	maxOffset := 8
	baseAlpha := 50
	for i := 1; i <= maxOffset; i++ {
		alpha := uint8(baseAlpha / i) // fade out as it goes further
		offset := float32(i)

		shadowColor := color.RGBA{0, 0, 0, alpha}
		vertices, indices := path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{Width: 3})

		for iv := range vertices {
			v := &vertices[iv]
			v.DstX += float32(x) + offset
			v.DstY += float32(y) + offset
		}

		img := ebiten.NewImage(3, 3)
		img.Fill(shadowColor)
		op := &ebiten.DrawTrianglesOptions{
			FillRule:  ebiten.FillRuleEvenOdd,
			AntiAlias: true,
		}
		screen.DrawTriangles(vertices, indices, img, op)
	}
}
