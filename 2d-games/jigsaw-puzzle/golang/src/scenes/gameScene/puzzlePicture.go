package gameScene

import (
	"fmt"
	"image"
	"math"
	"math/rand/v2"
	"sort"

	"github.com/DTVegaArchChapter/GameProgramming/jigsaw-puzzle/common"
	"github.com/hajimehoshi/ebiten/v2"
)

type Direction int

const (
	Left = Direction(iota)
	Right
	Up
	Down
)

type PuzzlePicture struct {
	image         *ebiten.Image
	Pieces        *[]*Piece
	Cols, Rows    int
	snappedPieces [][]int
}

func CreatePuzzlePicture(image *ebiten.Image) *PuzzlePicture {
	w, h := image.Bounds().Dx(), image.Bounds().Dy()
	scale := 1.0
	if w > 600 {
		scale = 600.0 / float64(w)
		w = int(float64(w) * scale)
		h = int(float64(h) * scale)
	}

	puzzleImage := ebiten.NewImage(w, h)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	puzzleImage.DrawImage(image, op)

	return &PuzzlePicture{image: puzzleImage, snappedPieces: [][]int{}}
}

func (p *PuzzlePicture) CreatePuzzlePieces(pieceCount int) {
	cols, rows := findClosestDivisors(pieceCount)
	pictureW, pictureH := p.image.Bounds().Dx(), p.image.Bounds().Dy()
	pieceW, pieceH := int(pictureW/cols), int(pictureH/rows)
	randomX1Min, randomX1Max, randomX2Min, randomX2Max, randomYMin, randomYMax := 0., 320.-float64(pieceW), 960., float64(common.ScreenWidth)-float64(pieceW), 68., float64(common.ScreenHeight)-float64(pieceH)-60

	p.Cols, p.Rows = cols, rows

	pieces := make([]*Piece, pieceCount)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			n := i*cols + j

			top, right, bottom, left := randomTabPosition(), randomTabPosition(), randomTabPosition(), randomTabPosition()

			if i == 0 {
				top = 0
			}

			if j == 0 {
				left = 0
			}

			if i == rows-1 {
				bottom = 0
			} else {
				if rand.N(2) == 0 {
					bottom = -bottom
				}
			}

			if j == cols-1 {
				right = 0
			} else {
				if rand.N(2) == 0 {
					right = -right
				}
			}

			if j > 0 {
				left = -(*pieces[i*cols+j-1]).GetRight()
			}

			if i > 0 {
				top = -(*pieces[(i-1)*cols+j]).GetBottom()
			}

			// Generate a random position for the piece within the bounds of the puzzle image
			var randomX float64
			if i%2 == 0 {
				randomX = randomX1Min + rand.Float64()*(randomX1Max-randomX1Min)

			} else {
				randomX = randomX2Min + rand.Float64()*(randomX2Max-randomX2Min)
			}

			randomY := randomYMin + rand.Float64()*(randomYMax-randomYMin)

			puzzlePiece := CreatePuzzlePiece(n, randomX, randomY, j*pieceW, i*pieceH, pieceW, pieceH, top, right, bottom, left, p)
			pieces[n] = &puzzlePiece
		}
	}

	p.Pieces = &pieces
}

func (p *PuzzlePicture) SubImage(rectangle image.Rectangle) *ebiten.Image {
	return p.image.SubImage(rectangle).(*ebiten.Image)
}

func (p *PuzzlePicture) Draw(screen *ebiten.Image) {
	for _, piece := range *p.Pieces {
		(*piece).Draw(screen)
	}
}

func (p *PuzzlePicture) SetPieceBeingDragged(mx, my float64) {
	pieces := *p.Pieces
	p.sortPiecesByZIndex()

	for i := len(pieces) - 1; i >= 0; i-- {
		piece := *pieces[i]
		if piece.Contains(mx, my) {
			i1 := p.getSnappedPieceSliceIndex(piece.GetN())
			if i1 == -1 {
				piece.SetBeingDragged(true)
				piece.ChangeOffset(mx, my)

				p.increaseZIndexOfPiece(i)
			} else {
				zIndex := p.increaseZIndexOfPiece(i)
				for _, n := range p.snappedPieces[i1] {
					p1Ptr, i2 := p.FindPiece(n)
					p1 := *p1Ptr
					p1.SetBeingDragged(true)
					p1.ChangeOffset(mx, my)

					(*pieces[i2]).SetZIndex(zIndex)
				}
			}

			p.sortPiecesByZIndex()

			break
		}
	}
}

func (p *PuzzlePicture) MoveDraggedPiece(mx, my float64) {
	pieces := *p.Pieces
	for _, piece := range pieces {
		if (*piece).GetBeingDragged() {
			(*piece).ChangePosition(mx, my)
		}
	}
}

func (p *PuzzlePicture) HandleDraggedPieceSnapping() {
	pieces := *p.Pieces
	for _, currPiecePtr := range pieces {
		currPiece := *currPiecePtr
		if currPiece.GetBeingDragged() {
			draggedPiece := currPiece

			last := len(pieces) - 1
			left := p.GetLeftN(draggedPiece)
			right := p.GetRightN(draggedPiece)
			top := p.GetTopN(draggedPiece)
			down := p.GetBottomN(draggedPiece)

			dragX, dragY := draggedPiece.GetCenterBoxCoordinates()

			if left >= 0 && left <= last {
				piecePtr, _ := p.FindPiece(left)
				piece := *piecePtr
				if !piece.GetBeingDragged() {
					pieceX, pieceY := piece.GetCenterBoxCoordinates()

					if math.Abs((pieceX+float64(piece.GetW()))-dragX) < 5 && math.Abs(pieceY-dragY) < 5 {
						p.shiftPieces(piece, draggedPiece)
						p.mergePieces(draggedPiece.GetN(), piece.GetN())
					}
				}
			}

			if right >= 0 && right <= last {
				piecePtr, _ := p.FindPiece(right)
				piece := *piecePtr
				if !piece.GetBeingDragged() {
					pieceX, pieceY := piece.GetCenterBoxCoordinates()

					if math.Abs(pieceX-(dragX+float64(draggedPiece.GetW()))) < 5 && math.Abs(pieceY-dragY) < 5 {
						p.shiftPieces(piece, draggedPiece)
						p.mergePieces(draggedPiece.GetN(), piece.GetN())
					}
				}
			}

			if top >= 0 && top <= last {
				piecePtr, _ := p.FindPiece(top)
				piece := *piecePtr
				if !piece.GetBeingDragged() {
					pieceX, pieceY := piece.GetCenterBoxCoordinates()

					if math.Abs(pieceX-dragX) < 5 && math.Abs((pieceY+float64(piece.GetH()))-dragY) < 5 {
						p.shiftPieces(piece, draggedPiece)
						p.mergePieces(draggedPiece.GetN(), piece.GetN())
					}
				}
			}

			if down >= 0 && down <= last {
				piecePtr, _ := p.FindPiece(down)
				piece := *piecePtr
				if !piece.GetBeingDragged() {
					pieceX, pieceY := piece.GetCenterBoxCoordinates()

					if math.Abs(pieceX-dragX) < 5 && math.Abs(pieceY-(dragY+float64(draggedPiece.GetH()))) < 5 {
						p.shiftPieces(piece, draggedPiece)
						p.mergePieces(draggedPiece.GetN(), piece.GetN())
					}
				}
			}
		}
	}
}

func (p *PuzzlePicture) GetCompletePercentage() int {
	total := 0
	for _, a := range p.snappedPieces {
		total += len(a)
	}

	return int(math.Floor(float64(total) / (float64(p.Cols) * float64(p.Rows)) * 100))
}

func (p *PuzzlePicture) shiftPieces(piece, draggedPiece Piece) {
	sX, sY := p.calculateShift(draggedPiece, piece)
	si := p.getSnappedPieceSliceIndex(piece.GetN())
	if si > -1 {
		for _, n := range p.snappedPieces[si] {
			p1Ptr, _ := p.FindPiece(n)
			p1 := *p1Ptr

			x, y := p1.GetCenterBoxCoordinates()
			p1.SetCenterBoxCoordinates(x+sX, y+sY)
		}
	} else {
		x, y := piece.GetCenterBoxCoordinates()
		piece.SetCenterBoxCoordinates(x+sX, y+sY)
	}
}

func (p *PuzzlePicture) calculateShift(draggedPiece, piece Piece) (float64, float64) {
	x, y := piece.GetCenterBoxCoordinates()
	dragX, dragY := draggedPiece.GetCenterBoxCoordinates()
	dirs := p.GetDirections(draggedPiece.GetN(), piece.GetN())
	for _, d := range dirs {
		if d == Left {
			return (dragX - float64(piece.GetW())) - x, dragY - y
		} else if d == Right {
			return (dragX + float64(piece.GetW())) - x, dragY - y
		} else if d == Up {
			return dragX - x, (dragY - float64(piece.GetH())) - y
		} else if d == Down {
			return dragX - x, (dragY + float64(piece.GetH())) - y
		}
	}

	return 0, 0
}

func (p *PuzzlePicture) GetLeftN(piece Piece) int {
	n := piece.GetN()
	if n%p.Cols == 0 {
		return -1
	}
	return n - 1
}

func (p *PuzzlePicture) GetRightN(piece Piece) int {
	n := piece.GetN()
	if n%p.Cols == p.Cols-1 {
		return -1
	}
	return n + 1
}

func (p *PuzzlePicture) GetTopN(piece Piece) int {
	n := piece.GetN()
	if n < p.Cols {
		return -1
	}

	t := n - p.Cols
	if t < 0 {
		return -1
	}

	return t
}

func (p *PuzzlePicture) GetBottomN(piece Piece) int {
	n := piece.GetN()
	if n >= p.Cols*(p.Rows-1) {
		return -1
	}

	b := n + p.Cols
	if b >= p.Cols*p.Rows {
		return -1
	}

	return b
}

func (p *PuzzlePicture) GetDirections(source, target int) []Direction {
	var directions []Direction
	rs, cs := p.GetRowIndex(source), p.GetColIndex(source)
	rt, ct := p.GetRowIndex(target), p.GetColIndex(target)

	if cs < ct {
		directions = append(directions, Right)
	} else if cs > ct {
		directions = append(directions, Left)
	}

	if rs < rt {
		directions = append(directions, Down)
	} else if rs > rt {
		directions = append(directions, Up)
	}

	return directions
}

func (p *PuzzlePicture) GetColIndex(n int) int {
	return n % p.Cols
}

func (p *PuzzlePicture) GetRowIndex(n int) int {
	return n / p.Cols
}

func (p *PuzzlePicture) IsPuzzleCompleted() bool {
	if len(p.snappedPieces) != 1 {
		return false
	}

	for _, a := range p.snappedPieces {
		if len(a) != p.Cols*p.Rows {
			return false
		}
	}

	return true
}

func (p *PuzzlePicture) GetImage() *ebiten.Image {
	return p.image
}

func (p *PuzzlePicture) mergePieces(n1, n2 int) {
	i1 := p.getSnappedPieceSliceIndex(n1)
	i2 := p.getSnappedPieceSliceIndex(n2)

	if i1 == -1 && i2 == -1 {
		p.snappedPieces = append(p.snappedPieces, []int{n1, n2})
	} else if i1 > -1 && i2 > -1 {
		if i1 == i2 {
			return
		}

		a1 := p.snappedPieces[i1]
		a2 := p.snappedPieces[i2]

		p.snappedPieces[i2] = append(a2, a1...)

		p.snappedPieces = removeAt(p.snappedPieces, i1)
	} else if i1 > -1 {
		a1 := p.snappedPieces[i1]

		p.snappedPieces[i1] = append(a1, n2)
	} else if i2 > -1 {
		a2 := p.snappedPieces[i2]

		p.snappedPieces[i2] = append(a2, n1)
	}
}

func (p *PuzzlePicture) getSnappedPieceSliceIndex(n int) int {
	for i, a := range p.snappedPieces {
		for _, v := range a {
			if v == n {
				return i
			}
		}
	}

	return -1
}

func removeAt[T any](s []T, i int) []T {
	return append(s[:i], s[i+1:]...)
}

func (p *PuzzlePicture) FindPiece(n int) (*Piece, int) {
	for i, piece := range *p.Pieces {
		if (*piece).GetN() == n {
			return piece, i
		}
	}

	return nil, -1
}

func (p *PuzzlePicture) DropPuzzlePieces() bool {
	result := false

	for _, piece := range *p.Pieces {
		p := *piece
		if p.GetBeingDragged() {
			p.SetBeingDragged(false)
			result = true
		}
	}

	return result
}

func (p *PuzzlePicture) sortPiecesByZIndex() {
	pieces := *p.Pieces
	sort.Slice(pieces, func(i, j int) bool {
		p1 := (*pieces[i])
		p2 := (*pieces[j])
		z1 := p1.GetZIndex()
		z2 := p2.GetZIndex()

		if z1 == z2 {
			return p1.GetN() < p2.GetN()
		}
		return z1 < z2
	})
}

func (p *PuzzlePicture) increaseZIndexOfPiece(i int) int {
	pieces := *p.Pieces

	min := -1
	for _, p := range pieces {
		if min == -1 || (*p).GetZIndex() < min {
			min = (*p).GetZIndex()
		}
	}

	if min > 0 {
		for _, p := range pieces {
			(*p).SetZIndex((*p).GetZIndex() - min)
		}
	}

	max := 0
	for _, p := range pieces {
		if (*p).GetZIndex() > max {
			max = (*p).GetZIndex()
		}
	}

	(*pieces[i]).SetZIndex(max + 1)

	return max + 1
}

func randomTabPosition() float64 {
	return 0.3 + (0.7-0.3)*rand.Float64()
}

func findClosestDivisors(c int) (int, int) {
	d1, d2, minDiff := 0, 0, math.MaxInt

	for i := int(math.Sqrt(float64(c))); i > 0; i-- {
		if c%i == 0 {
			diff := int(math.Abs(float64(i - c/i)))
			if diff < minDiff {
				minDiff = diff
				d1 = int(math.Max(float64(i), float64(c/i)))
				d2 = int(math.Min(float64(i), float64(c/i)))
			}
		}
	}

	if minDiff == math.MaxInt {
		panic(fmt.Errorf("%c has no closest divisors", c))
	}

	if float64(d1)/float64(d2) > 3 {
		panic(fmt.Errorf("closest divisors' ratio of %c cannot be bigger than 3", c))
	}

	return d1, d2
}
