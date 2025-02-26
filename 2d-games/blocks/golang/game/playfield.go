package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PlayField struct {
	x          int
	y          int
	width      int
	height     int
	tileSize   int
	blocks     [][]color.Color
	emptyColor color.Color
}

func newPlayField(x, y, width, height, tileSize int) *PlayField {
	emptyColor := color.RGBA{0, 0, 0, 220}
	blocks := make([][]color.Color, height)
	for i := 0; i < height; i++ {
		blocks[i] = make([]color.Color, width)

		for j := range blocks[i] {
			blocks[i][j] = emptyColor
		}
	}

	return &PlayField{
		x:          x,
		y:          y,
		width:      width,
		height:     height,
		tileSize:   tileSize,
		blocks:     blocks,
		emptyColor: emptyColor,
	}
}

func (p *PlayField) GetSize() (int, int) {
	return p.width * p.tileSize, p.height * p.tileSize
}

func (p *PlayField) Draw(screen *ebiten.Image) {
	for i := range p.blocks {
		for j := range p.blocks[i] {
			p.FillBlock(screen, float32(j), float32(i), p.blocks[i][j])
		}
	}
}

func (p *PlayField) SetBlock(x, y int, color color.Color) {
	p.blocks[y][x] = color
}

func (p *PlayField) FillBlock(screen *ebiten.Image, x, y float32, color color.Color) {
	if y < 0 {
		return
	}

	vector.DrawFilledRect(screen, float32(p.x)+x*float32(p.tileSize), float32(p.y)+y*float32(p.tileSize), float32(p.tileSize), float32(p.tileSize), color, false)
}

func (p *PlayField) IsBlocked(x, y int) bool {
	if x < 0 || x >= p.width {
		return true
	}

	if y >= p.height {
		return true
	}

	if y < 0 {
		return false
	}

	return p.blocks[y][x] != p.emptyColor
}

func (p *PlayField) ClearLines() int {
	l := 0
	for r := 0; r < len(p.blocks); r++ {
		clearLine := true
		for _, b := range p.blocks[r] {
			if b == p.emptyColor {
				clearLine = false
				break
			}
		}

		if clearLine {
			l++
			for c := range p.blocks[r] {
				p.blocks[r][c] = p.emptyColor
			}

			for n := r; n >= 0; n-- {
				for c := range p.blocks[n] {
					if n == 0 {
						p.blocks[n][c] = p.emptyColor
					} else {
						p.blocks[n][c] = p.blocks[n-1][c]
					}
				}
			}
		}
	}

	return l
}
