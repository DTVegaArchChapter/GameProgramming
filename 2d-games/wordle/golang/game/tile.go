package wordle

import (
	"image/color"
	"math"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type tile struct {
	x      float64
	y      float64
	size   int
	col    int
	row    int
	r      rune
	status CharacterStatus

	isFlipping   bool
	flipProgress float64

	tween *Tween
	board *board
}

func newTile(col, row int, board *board) *tile {
	size := 60

	return &tile{
		x:            44 + float64(col*(size+3)+3),
		y:            float64(row*(size+3) + 3),
		size:         size,
		col:          col,
		row:          row,
		board:        board,
		isFlipping:   false,
		flipProgress: 1,
	}
}

func (t *tile) Update() {
	if t.tween != nil {
		t.tween.Update(1)
	}
}

func (t *tile) Draw(screen *ebiten.Image) {
	outerRect := ebiten.NewImage(t.size, t.size)
	pos := t.board.calcPos(t.col, t.row)
	borderColor := lightGrayColor
	if pos == t.board.pos {
		borderColor = grayColor
	}

	outerRect.Fill(borderColor)

	innerRect := ebiten.NewImage(t.size-2, t.size-2)
	tileColor := getTileColor(t)

	innerRect.Fill(tileColor)

	r := t.r
	if r != 0 {
		fontColor := color.Black
		if tileColor != color.White {
			fontColor = color.White
		}

		textOp := &text.DrawOptions{}
		textOp.GeoM.Translate(float64(innerRect.Bounds().Max.X)/2, float64(innerRect.Bounds().Max.Y)/2)
		textOp.ColorScale.ScaleWithColor(fontColor)
		textOp.PrimaryAlign = text.AlignCenter
		textOp.SecondaryAlign = text.AlignCenter
		text.Draw(innerRect, string(unicode.TurkishCase.ToUpper(r)), &text.GoTextFace{
			Source: mplusRegularTextFaceSource,
			Size:   float64(fontSize),
		}, textOp)
	}

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(1, 1)
	outerRect.DrawImage(innerRect, op)

	op.GeoM.Reset()

	if t.isFlipping {
		angle := t.flipProgress * math.Pi
		scale := math.Cos(angle)
		skew := math.Sin(angle) * 0.075
		op.GeoM.Translate(-float64(t.size/2), 0)
		op.GeoM.Skew(skew, 0)
		op.GeoM.Scale(math.Abs(scale), 1)
		op.GeoM.Translate(float64(t.size/2), 0)
	}

	op.GeoM.Translate(t.x, t.y)

	screen.DrawImage(outerRect, op)
}

func (t *tile) clearValue() {
	t.clearRune()
	t.clearStatus()
}

func (t *tile) clearRune() {
	t.setRune(0)
}

func (t *tile) setRune(r rune) {
	t.r = r
}

func (t *tile) clearStatus() {
	t.setStatus(CharacterStatusNone)
}

func (t *tile) setStatus(s CharacterStatus) {
	t.status = s
}

func (t *tile) isCharStatusNone() bool {
	return t.status == CharacterStatusNone
}

func (t *tile) isEmpty() bool {
	return t.r == 0
}

func (t *tile) startShake() {
	x := t.x
	t.tween = NewTween(0, 1, 0.5*60, LinearTweenFunc, func(v float64) {
		t.x = x + math.Sin(v*3*2*math.Pi)*5
	}, func() {
		t.x = x
		t.tween = nil
	})
}

func (r *tile) flip() {
	r.isFlipping = true
}

func getTileColor(t *tile) color.Color {
	var defaultColor color.Color = color.White

	if t.isFlipping {
		t.flipProgress -= 0.075
		if t.flipProgress <= 0 {
			t.isFlipping = false
			t.flipProgress = 0
		}
	}

	if math.Cos(t.flipProgress*math.Pi) < 0 {
		return defaultColor
	} else {
		st := t.status

		return st.getColor()
	}
}
