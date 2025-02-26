package game

import (
	"image"
	"image/color"
	"math"

	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type PieceType int

const (
	PieceTypeI PieceType = iota
	PieceTypeJ
	PieceTypeL
	PieceTypeO
	PieceTypeS
	PieceTypeT
	PieceTypeZ
	PieceTypeII
	PieceTypeIII
	PieceTypeDot
)

type pieceDefinition struct {
	blocks         []Point
	pivotIndex     int
	rotationAngles []int
	color          color.Color
}

var pieceDefinitions map[PieceType]pieceDefinition = map[PieceType]pieceDefinition{
	PieceTypeI: pieceDefinition{
		blocks: []Point{
			{0, 0},
			{1, 0}, // Pivot
			{2, 0},
			{3, 0},
		},
		pivotIndex:     1,
		rotationAngles: []int{90, -90},
		color:          color.RGBA{R: 175, G: 238, B: 238, A: 255},
	},
	PieceTypeJ: pieceDefinition{
		blocks: []Point{
			{0, 0},
			{0, 1},
			{1, 1}, // Pivot
			{2, 1},
		},
		pivotIndex:     2,
		rotationAngles: []int{90},
		color:          color.RGBA{R: 137, G: 207, B: 240, A: 255},
	},
	PieceTypeL: pieceDefinition{
		blocks: []Point{
			{0, 1},
			{1, 1}, // Pivot
			{2, 1},
			{2, 0},
		},
		pivotIndex:     1,
		rotationAngles: []int{90},
		color:          color.RGBA{R: 255, G: 179, B: 102, A: 255},
	},
	PieceTypeO: pieceDefinition{
		blocks: []Point{
			{0, 0},
			{0, 1},
			{1, 0},
			{1, 1},
		},
		pivotIndex:     -1,
		rotationAngles: []int{90},
		color:          color.RGBA{R: 253, G: 253, B: 150, A: 255},
	},
	PieceTypeS: pieceDefinition{
		blocks: []Point{
			{0, 1},
			{1, 1}, // Pivot
			{1, 0},
			{2, 0},
		},
		pivotIndex:     1,
		rotationAngles: []int{90, -90},
		color:          color.RGBA{R: 119, G: 221, B: 119, A: 255},
	},
	PieceTypeT: pieceDefinition{
		blocks: []Point{
			{0, 1},
			{1, 1}, // Pivot
			{2, 1},
			{1, 0},
		},
		pivotIndex:     1,
		rotationAngles: []int{90},
		color:          color.RGBA{R: 216, G: 191, B: 216, A: 255},
	},
	PieceTypeZ: pieceDefinition{
		blocks: []Point{
			{0, 0},
			{1, 0},
			{1, 1}, // Pivot
			{2, 1},
		},
		pivotIndex:     2,
		rotationAngles: []int{90, -90},
		color:          color.RGBA{R: 255, G: 153, B: 153, A: 255},
	},
	PieceTypeII: pieceDefinition{
		blocks: []Point{
			{0, 0}, // Pivot
			{1, 0},
		},
		pivotIndex:     0,
		rotationAngles: []int{90, -90},
		color:          color.RGBA{R: 200, G: 200, B: 200, A: 255},
	},
	PieceTypeIII: pieceDefinition{
		blocks: []Point{
			{0, 0},
			{1, 0}, // Pivot
			{2, 0},
		},
		pivotIndex:     1,
		rotationAngles: []int{90, -90},
		color:          color.RGBA{R: 181, G: 101, B: 29, A: 255},
	},
	PieceTypeDot: pieceDefinition{
		blocks: []Point{
			{0, 0},
		},
		pivotIndex:     -1,
		rotationAngles: []int{90},
		color:          color.RGBA{R: 191, G: 255, B: 164, A: 255},
	},
}

type Piece struct {
	playField      *PlayField
	blocks         *[]*Point
	pivotIndex     int
	color          color.Color
	rotationAngles []int
	rotationIndex  int
}

func (piece *Piece) Turn() {
	if piece.pivotIndex < 0 {
		return
	}

	angle := piece.rotationAngles[piece.rotationIndex]
	piece.rotationIndex = (piece.rotationIndex + 1) % (len(piece.rotationAngles))

	piece.turn(angle)

	if piece.collides() {
		piece.turn(-angle)
	}
}

func (piece *Piece) MoveLeft() {
	piece.translate(-1, 0)

	if piece.collides() {
		piece.MoveRight()
	}
}

func (piece *Piece) MoveRight() {
	piece.translate(1, 0)

	if piece.collides() {
		piece.MoveLeft()
	}
}

func (piece *Piece) MoveDown() bool {
	piece.translate(0, 1)

	if piece.collides() {
		piece.translate(0, -1)
		return false
	}

	return true
}

func (piece *Piece) AbsorbIntoPlayField() {
	for _, p := range *piece.blocks {
		piece.playField.SetBlock(p.X, p.Y, piece.color)
	}
}

func (piece *Piece) Draw(screen *ebiten.Image) {
	for _, p := range *piece.blocks {
		piece.playField.FillBlock(screen, float32(p.X), float32(p.Y), piece.color)
	}
}

func createNewPiece(playField *PlayField) *Piece {
	t := PieceType(rand.Intn(len(pieceDefinitions)))
	d := pieceDefinitions[t]
	blocks := make([]*Point, len(d.blocks))

	for i, b := range d.blocks {
		blocks[i] = &b
	}

	p := Piece{
		blocks:         &blocks,
		pivotIndex:     d.pivotIndex,
		color:          d.color,
		playField:      playField,
		rotationAngles: d.rotationAngles,
		rotationIndex:  0,
	}

	return &p
}

func (piece *Piece) getRectangle() image.Rectangle {
	minX, minY, maxX, maxY := 0, 0, 0, 0

	for _, p := range *piece.blocks {
		if p.X < minX {
			minX = p.X
		}

		if p.Y < minY {
			minY = p.Y
		}

		if p.X > maxX {
			maxX = p.X
		}

		if p.Y > maxY {
			maxY = p.Y
		}
	}

	return image.Rect(minX, minY, maxX+1, maxY+1)
}

func (piece *Piece) getPivot() Point {
	return *(*piece.blocks)[piece.pivotIndex]
}

func (piece *Piece) translate(x, y int) {
	for _, p := range *piece.blocks {
		p.X += x
		p.Y += y
	}
}

func (piece *Piece) rotate(angle int) {
	rad := float64(angle) * math.Pi / 180
	cos := int(math.Cos(rad))
	sin := int(math.Sin(rad))

	for _, p := range *piece.blocks {
		x := p.X
		y := p.Y

		p.X = x*cos - y*sin
		p.Y = x*sin + y*cos
	}
}

func (piece *Piece) turn(angle int) {
	pivot := piece.getPivot()
	piece.translate(-pivot.X, -pivot.Y)
	piece.rotate(angle)
	piece.translate(pivot.X, pivot.Y)
}

func (piece *Piece) collides() bool {
	for _, b := range *piece.blocks {
		if piece.playField.IsBlocked(b.X, b.Y) {
			return true
		}
	}

	return false
}
