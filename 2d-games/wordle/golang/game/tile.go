package wordle

import (
	"container/list"
	"image/color"
	"math"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	flipDuration   = 0.5 * 60
	bounceDuration = 0.2 * 60
)

type tile struct {
	x      float64
	y      float64
	size   int
	col    int
	row    int
	r      rune
	status CharacterStatus

	text            *TextRenderer
	fontColor       color.Color
	borderColor     color.Color
	backgroundColor color.Color

	flipped bool

	geoM  *ebiten.GeoM
	tween *list.List
	board *board
}

func newTile(col, row int, board *board) *tile {
	size := 60

	return &tile{
		x:               44 + float64(col*(size+3)+3),
		y:               float64(row*(size+3) + 3),
		size:            size,
		col:             col,
		row:             row,
		board:           board,
		text:            NewTextRenderer(RobotoBoldFontName, color.Black, 30),
		fontColor:       color.Black,
		borderColor:     lightGrayColor,
		backgroundColor: color.White,
		tween:           list.New(),
	}
}

func (t *tile) Update() {
	t.updateBackgroundColor()

	t.fontColor = color.Black
	if t.backgroundColor != color.White {
		t.fontColor = color.White
	}

	if t.backgroundColor != color.White {
		t.borderColor = t.backgroundColor
	} else if t.r > 0 || t.board.calcPos(t.col, t.row) == t.board.pos {
		t.borderColor = grayColor
	} else {
		t.borderColor = lightGrayColor
	}

	if e := t.tween.Front(); e != nil {
		tween := e.Value.(*Tween)
		tween.Update(1)
		if tween.isCompleted {
			t.tween.Remove(e)
		}
	}
}

func (t *tile) Draw(screen *ebiten.Image) {
	outerRect := ebiten.NewImage(t.size, t.size)
	outerRect.Fill(t.borderColor)

	innerRect := ebiten.NewImage(t.size-4, t.size-4)
	innerRect.Fill(t.backgroundColor)

	r := t.r
	if r != 0 {
		t.text.SetColor(t.fontColor)
		t.text.Draw(innerRect, string(unicode.TurkishCase.ToUpper(r)), innerRect.Bounds().Dx()/2, innerRect.Bounds().Dy()/2)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(2, 2)
	outerRect.DrawImage(innerRect, op)

	op.GeoM.Reset()

	if t.geoM != nil {
		op.GeoM.Concat(*t.geoM)
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

	if r > 0 {
		t.tween.PushBack(NewTween(0, 1, 0.125*60, LinearTweenFunc, func(v float64) {
			scale := 1 + math.Sin(v*math.Pi)*0.06
			middle := float64(t.size) / 2

			geoM := &ebiten.GeoM{}
			geoM.Translate(-middle, -middle)
			geoM.Scale(scale, scale)
			geoM.Translate(middle, middle)

			t.geoM = geoM
		}, func() {
			t.geoM = nil
		}))
	}
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

func (t *tile) shake() {
	if t.tween.Len() > 0 {
		return
	}

	x := t.x
	t.tween.PushBack(NewTween(0, 1, 0.5*60, LinearTweenFunc, func(v float64) {
		t.x = x + math.Sin(v*3*2*math.Pi)*5
	}, func() {
		t.x = x
	}))
}

func (t *tile) flip() {
	t.appendDelayTween(func(c int) float64 { return float64(c) * flipDuration })
	t.tween.PushBack(NewTween(0, 1, 0.5*60, LinearTweenFunc, func(v float64) {
		angle := v * math.Pi
		scale := math.Cos(angle)

		geoM := &ebiten.GeoM{}
		geoM.Translate(0, -float64(t.size/2))
		geoM.Scale(1, math.Abs(scale))
		geoM.Translate(0, float64(t.size/2))

		t.geoM = geoM

		t.flipped = math.Cos(v*math.Pi) < 0
	}, func() {
		t.geoM = nil
	}))
}

func (t *tile) celebrateWin() {
	y := t.y
	t.appendDelayTween(func(c int) float64 { return (float64(t.board.cols-c-1) * flipDuration) + float64(c)*bounceDuration })
	t.tween.PushBack(NewTween(0, 1, 0.5*60, LinearTweenFunc,
		func(v float64) {
			t.y = y - math.Sin(v*math.Pi)*10
		}, func() {
			t.y = y
			t.board.tileWinAnimationFinishedCounter++
		}))
}

func (t *tile) updateBackgroundColor() {
	if t.flipped {
		st := t.status

		t.backgroundColor = st.getColor()
	} else {
		t.backgroundColor = color.White
	}
}

func (t *tile) appendDelayTween(delayFunc func(c int) float64) {
	colNum := t.col % t.board.cols
	delay := delayFunc(colNum)
	if delay == 0 {
		delay = 1
	}

	t.tween.PushBack(NewTween(0, 1, delay, LinearTweenFunc, func(v float64) {}, func() {}))
}
