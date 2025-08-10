package gameScene

import (
	"image"
	"image/color"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/DTVegaArchChapter/GameProgramming/jigsaw-puzzle/common"
)

type Piece interface {
	GetN() int
	GetZIndex() int
	SetZIndex(i int)
	Contains(x, y float64) bool
	Draw(screen *ebiten.Image)
	GetCenterBoxCoordinates() (float64, float64)
	SetCenterBoxCoordinates(x, y float64)
	SetBeingDragged(b bool)
	GetBeingDragged() bool
	GetW() int
	GetH() int
	ChangePosition(mx, my float64)
	ChangeOffset(mx, my float64)
	GetBottom() float64
	GetRight() float64
	GetPathVertices() []common.Point
	SetPosition(x, y float64)
}

type PuzzlePiece struct {
	N                                                                            int
	X, Y                                                                         float64
	W, H, PicX, PicY, VTabWidth, VTabHeight, HTabWidth, HTabHeight, VNeck, HNeck int
	OffsetX                                                                      float64
	OffsetY                                                                      float64
	BeingDragged                                                                 bool
	Top, Right, Bottom, Left                                                     float64
	Path                                                                         *vector.Path
	TopPath                                                                      *vector.Path
	BottomPath                                                                   *vector.Path
	LeftPath                                                                     *vector.Path
	RightPath                                                                    *vector.Path
	PathVertices                                                                 []common.Point
	Image                                                                        *ebiten.Image
	PuzzlePicture                                                                *PuzzlePicture
	ZIndex                                                                       int
}

func CreatePuzzlePiece(n int, x, y float64, picX, picY, w, h int, top, right, bottom, left float64, puzzlePicture *PuzzlePicture) Piece {
	p := &PuzzlePiece{
		N:             n,
		X:             x,
		Y:             y,
		PicX:          picX,
		PicY:          picY,
		Top:           top,
		Right:         right,
		Bottom:        bottom,
		Left:          left,
		PuzzlePicture: puzzlePicture,
		W:             w,
		H:             h,
		VTabWidth:     int(0.2 * float32(h)),
		VTabHeight:    int(0.2 * float32(h)),
		VNeck:         int(0.1 * float32(h)),
		HTabWidth:     int(0.2 * float32(w)),
		HTabHeight:    int(0.2 * float32(w)),
		HNeck:         int(0.1 * float32(w)),
		ZIndex:        0,
	}

	if left > 0 {
		p.X = p.X - float64(p.HTabHeight)
		p.PicX = p.PicX - p.HTabHeight
	}

	if top > 0 {
		p.Y = p.Y - float64(p.VTabHeight)
		p.PicY = p.PicY - p.VTabHeight
	}

	p.generatePath()
	return p
}

func (p *PuzzlePiece) Contains(x, y float64) bool {
	localPoint := common.Point{
		X: x - p.X,
		Y: y - p.Y,
	}
	return common.IsPointInPathVertices(localPoint, p.PathVertices)
}

func (p *PuzzlePiece) Draw(screen *ebiten.Image) {
	if p.BeingDragged {
		p.drawShadow(screen)
	}

	mx := ebiten.GeoM{}
	mx.Translate(p.X, p.Y)
	screen.DrawImage(p.Image, &ebiten.DrawImageOptions{GeoM: mx})
}

func (p *PuzzlePiece) GetN() int {
	return p.N
}

func (p *PuzzlePiece) GetZIndex() int {
	return p.ZIndex
}

func (p *PuzzlePiece) SetZIndex(i int) {
	p.ZIndex = i
}

func (p *PuzzlePiece) SetBeingDragged(b bool) {
	p.BeingDragged = b
}

func (p *PuzzlePiece) GetBeingDragged() bool {
	return p.BeingDragged
}

func (p *PuzzlePiece) GetW() int {
	return p.W
}

func (p *PuzzlePiece) GetH() int {
	return p.H
}

func (p *PuzzlePiece) ChangeOffset(mx, my float64) {
	p.OffsetX = p.X - mx
	p.OffsetY = p.Y - my
}

func (p *PuzzlePiece) ChangePosition(mx, my float64) {
	p.X = mx + p.OffsetX
	p.Y = my + p.OffsetY
}

func (p *PuzzlePiece) GetBottom() float64 {
	return p.Bottom
}

func (p *PuzzlePiece) GetRight() float64 {
	return p.Right
}

func (p *PuzzlePiece) GetPathVertices() []common.Point {
	return p.PathVertices
}

func (p *PuzzlePiece) getBoxSize() (int, int) {
	tabWidth := 0
	if p.Left > 0 {
		tabWidth += p.HTabWidth
	}

	if p.Right > 0 {
		tabWidth += p.HTabWidth
	}

	tabHeight := 0
	if p.Top > 0 {
		tabHeight += p.VTabHeight
	}

	if p.Bottom > 0 {
		tabHeight += p.VTabHeight
	}

	return p.W + tabWidth, p.H + tabHeight
}

func (p *PuzzlePiece) GetCenterBoxCoordinates() (float64, float64) {
	x, y := p.X, p.Y
	if p.Left > 0 {
		x += float64(p.HTabHeight)
	}

	if p.Top > 0 {
		y += float64(p.VTabHeight)
	}

	return x, y
}

func (p *PuzzlePiece) SetCenterBoxCoordinates(x, y float64) {
	if p.Left > 0 {
		x -= float64(p.HTabHeight)
	}

	if p.Top > 0 {
		y -= float64(p.VTabHeight)
	}

	p.X = x
	p.Y = y
}

func (p *PuzzlePiece) generatePath() {
	boxSizeW, boxSizeH := p.getBoxSize()

	path := vector.Path{}
	topPath := vector.Path{}
	bottomPath := vector.Path{}
	leftPath := vector.Path{}
	rightPath := vector.Path{}

	w, h := float32(p.W), float32(p.H)
	hTabWidth, vTabWidth := float32(p.HTabWidth), float32(p.VTabWidth)
	hTabHeight, vTabHeight := float32(p.HTabHeight), float32(p.VTabHeight)
	hNeck, vNeck := float32(p.HNeck), float32(p.VNeck)

	var x, y float32
	if p.Top <= 0 {
		y = 0
	} else {
		y = vTabHeight
	}

	if p.Left <= 0 {
		x = 0
	} else {
		x = hTabHeight
	}

	// Top
	path.MoveTo(x, y)
	topPath.MoveTo(x, y)

	if p.Top != 0 {
		path.LineTo(x+w*float32(math.Abs(p.Top))-vNeck, y)
		topPath.LineTo(x+w*float32(math.Abs(p.Top))-vNeck, y)

		path.CubicTo(
			x+w*float32(math.Abs(p.Top))-vNeck, y-vTabHeight*0.2*sign(p.Top),
			x+w*float32(math.Abs(p.Top))-vTabWidth, y-vTabHeight*sign(p.Top),
			x+w*float32(math.Abs(p.Top)), y-vTabHeight*sign(p.Top),
		)
		path.CubicTo(
			x+w*float32(math.Abs(p.Top))+vTabWidth, y-vTabHeight*sign(p.Top),
			x+w*float32(math.Abs(p.Top))+vNeck, y-vTabHeight*0.2*sign(p.Top),
			x+w*float32(math.Abs(p.Top))+vNeck, y,
		)

		topPath.CubicTo(
			x+w*float32(math.Abs(p.Top))-vNeck, y-vTabHeight*0.2*sign(p.Top),
			x+w*float32(math.Abs(p.Top))-vTabWidth, y-vTabHeight*sign(p.Top),
			x+w*float32(math.Abs(p.Top)), y-vTabHeight*sign(p.Top),
		)
		topPath.CubicTo(
			x+w*float32(math.Abs(p.Top))+vTabWidth, y-vTabHeight*sign(p.Top),
			x+w*float32(math.Abs(p.Top))+vNeck, y-vTabHeight*0.2*sign(p.Top),
			x+w*float32(math.Abs(p.Top))+vNeck, y,
		)
	}

	// Right
	path.LineTo(x+w, y)
	topPath.LineTo(x+w, y)
	rightPath.MoveTo(x+w, y)

	if p.Right != 0 {
		path.LineTo(x+w, y+h*float32(math.Abs(p.Right))-hNeck)
		rightPath.LineTo(x+w, y+h*float32(math.Abs(p.Right))-hNeck)

		path.CubicTo(
			x+w+hTabHeight*0.2*sign(p.Right), y+h*float32(math.Abs(p.Right))-hNeck,
			x+w+hTabHeight*sign(p.Right), y+h*float32(math.Abs(p.Right))-hTabWidth,
			x+w+hTabHeight*sign(p.Right), y+h*float32(math.Abs(p.Right)),
		)
		path.CubicTo(
			x+w+hTabHeight*sign(p.Right), y+h*float32(math.Abs(p.Right))+hTabWidth,
			x+w+hTabHeight*0.2*sign(p.Right), y+h*float32(math.Abs(p.Right))+hNeck,
			x+w, y+h*float32(math.Abs(p.Right))+hNeck,
		)

		rightPath.CubicTo(
			x+w+hTabHeight*0.2*sign(p.Right), y+h*float32(math.Abs(p.Right))-hNeck,
			x+w+hTabHeight*sign(p.Right), y+h*float32(math.Abs(p.Right))-hTabWidth,
			x+w+hTabHeight*sign(p.Right), y+h*float32(math.Abs(p.Right)),
		)
		rightPath.CubicTo(
			x+w+hTabHeight*sign(p.Right), y+h*float32(math.Abs(p.Right))+hTabWidth,
			x+w+hTabHeight*0.2*sign(p.Right), y+h*float32(math.Abs(p.Right))+hNeck,
			x+w, y+h*float32(math.Abs(p.Right))+hNeck,
		)
	}

	// Bottom
	path.LineTo(x+w, y+h)
	rightPath.LineTo(x+w, y+h)
	bottomPath.MoveTo(x+w, y+h)

	if p.Bottom != 0 {
		path.LineTo(x+w*float32(math.Abs(p.Bottom))-vNeck, y+h)
		bottomPath.LineTo(x+w*float32(math.Abs(p.Bottom))-vNeck, y+h)

		path.CubicTo(
			x+w*float32(math.Abs(p.Bottom))-vNeck, y+h+vTabHeight*0.2*sign(p.Bottom),
			x+w*float32(math.Abs(p.Bottom))-vTabWidth, y+h+vTabHeight*sign(p.Bottom),
			x+w*float32(math.Abs(p.Bottom)), y+h+vTabHeight*sign(p.Bottom),
		)
		path.CubicTo(
			x+w*float32(math.Abs(p.Bottom))+vTabWidth, y+h+vTabHeight*sign(p.Bottom),
			x+w*float32(math.Abs(p.Bottom))+vNeck, y+h+vTabHeight*0.2*sign(p.Bottom),
			x+w*float32(math.Abs(p.Bottom))+vNeck, y+h,
		)

		bottomPath.CubicTo(
			x+w*float32(math.Abs(p.Bottom))-vNeck, y+h+vTabHeight*0.2*sign(p.Bottom),
			x+w*float32(math.Abs(p.Bottom))-vTabWidth, y+h+vTabHeight*sign(p.Bottom),
			x+w*float32(math.Abs(p.Bottom)), y+h+vTabHeight*sign(p.Bottom),
		)
		bottomPath.CubicTo(
			x+w*float32(math.Abs(p.Bottom))+vTabWidth, y+h+vTabHeight*sign(p.Bottom),
			x+w*float32(math.Abs(p.Bottom))+vNeck, y+h+vTabHeight*0.2*sign(p.Bottom),
			x+w*float32(math.Abs(p.Bottom))+vNeck, y+h,
		)
	}

	// Left
	path.LineTo(x, y+h)
	bottomPath.LineTo(x, y+h)
	leftPath.MoveTo(x, y+h)

	if p.Left != 0 {
		path.LineTo(x, y+h*float32(math.Abs(p.Left))-hNeck)
		leftPath.LineTo(x, y+h*float32(math.Abs(p.Left))-hNeck)

		path.CubicTo(
			x-hTabHeight*0.2*sign(p.Left), y+h*float32(math.Abs(p.Left))-hNeck,
			x-hTabHeight*sign(p.Left), y+h*float32(math.Abs(p.Left))-hTabWidth,
			x-hTabHeight*sign(p.Left), y+h*float32(math.Abs(p.Left)),
		)
		path.CubicTo(
			x-hTabHeight*sign(p.Left), y+h*float32(math.Abs(p.Left))+hTabWidth,
			x-hTabHeight*0.2*sign(p.Left), y+h*float32(math.Abs(p.Left))+hNeck,
			x, y+h*float32(math.Abs(p.Left))+hNeck,
		)

		leftPath.CubicTo(
			x-hTabHeight*0.2*sign(p.Left), y+h*float32(math.Abs(p.Left))-hNeck,
			x-hTabHeight*sign(p.Left), y+h*float32(math.Abs(p.Left))-hTabWidth,
			x-hTabHeight*sign(p.Left), y+h*float32(math.Abs(p.Left)),
		)
		leftPath.CubicTo(
			x-hTabHeight*sign(p.Left), y+h*float32(math.Abs(p.Left))+hTabWidth,
			x-hTabHeight*0.2*sign(p.Left), y+h*float32(math.Abs(p.Left))+hNeck,
			x, y+h*float32(math.Abs(p.Left))+hNeck,
		)
	}

	path.Close()
	leftPath.LineTo(x, y)

	p.TopPath = &topPath
	p.BottomPath = &bottomPath
	p.LeftPath = &leftPath
	p.RightPath = &rightPath
	p.Path = &path
	p.PathVertices = common.GetPathVertices(&path)

	mask := newMask(boxSizeW, boxSizeH, path)

	// Create intermediate image
	intermediate := ebiten.NewImage(boxSizeW, boxSizeH)

	// Draw mask to intermediate
	intermediate.DrawImage(mask, nil)

	subImage := p.PuzzlePicture.SubImage(image.Rect(p.PicX, p.PicY, p.PicX+boxSizeW, p.PicY+boxSizeH))

	// Apply composite operation
	op := &ebiten.DrawImageOptions{}
	op.Blend = ebiten.BlendSourceIn
	intermediate.DrawImage(subImage, op)

	p.Image = intermediate
}

func (p *PuzzlePiece) drawShadow(screen *ebiten.Image) {
	i1 := p.PuzzlePicture.getSnappedPieceSliceIndex(p.N)
	if i1 < 0 {
		common.DrawShadowForPath(screen, p.X, p.Y, p.Path)
		return
	}

	snappedPieces := p.PuzzlePicture.snappedPieces[i1]

	topN := p.PuzzlePicture.GetTopN(p)
	if !slices.ContainsFunc(snappedPieces, func(n int) bool { return n == topN }) {
		common.DrawShadowForPath(screen, p.X, p.Y, p.TopPath)
	}

	bottomN := p.PuzzlePicture.GetBottomN(p)
	if !slices.ContainsFunc(snappedPieces, func(n int) bool { return n == bottomN }) {
		common.DrawShadowForPath(screen, p.X, p.Y, p.BottomPath)
	}

	leftN := p.PuzzlePicture.GetLeftN(p)
	if !slices.ContainsFunc(snappedPieces, func(n int) bool { return n == leftN }) {
		common.DrawShadowForPath(screen, p.X, p.Y, p.LeftPath)
	}

	rightN := p.PuzzlePicture.GetRightN(p)
	if !slices.ContainsFunc(snappedPieces, func(n int) bool { return n == rightN }) {
		common.DrawShadowForPath(screen, p.X, p.Y, p.RightPath)
	}
}

func (p *PuzzlePiece) SetPosition(x, y float64) {
	p.X = x
	p.Y = y
}

func sign(n float64) float32 {
	if n > 0 {
		return 1
	}

	if n < 0 {
		return -1
	}

	return 0
}

func newMask(width, height int, path vector.Path) *ebiten.Image {
	whiteImg := ebiten.NewImage(1, 1)
	whiteImg.Fill(color.White)

	// Generate vertices and indices from path
	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)

	// Create mask
	mask := ebiten.NewImage(width, height)
	mask.Fill(color.Transparent)

	// Draw the path to mask using triangles
	op := &ebiten.DrawTrianglesOptions{
		FillRule:  ebiten.FillRuleEvenOdd,
		AntiAlias: true,
	}
	mask.DrawTriangles(vertices, indices, whiteImg, op)
	return mask
}
