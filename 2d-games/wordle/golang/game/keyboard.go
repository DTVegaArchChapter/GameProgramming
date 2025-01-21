package wordle

import (
	"image/color"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type keyboard struct {
	rows    *[3]*keyboardRow
	keysMap *map[rune]*keyboardKey
	boxW    int
	boxH    int
	boxGap  int
}

type keyboardRow struct {
	keys  *[]*keyboardKey
	width int
}

type keyboardKey struct {
	char   rune
	status CharacterStatus
}

func newKeyboard() *keyboard {
	keyboard := &keyboard{
		boxW:   30,
		boxH:   50,
		boxGap: 3,
		rows: &[3]*keyboardRow{
			&keyboardRow{
				keys: &[]*keyboardKey{
					&keyboardKey{
						char:   'E',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'R',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'T',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'Y',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'U',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'I',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'O',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'P',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'Ğ',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'Ü',
						status: CharacterStatusNone,
					},
				},
			},
			&keyboardRow{
				keys: &[]*keyboardKey{
					&keyboardKey{
						char:   'A',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'S',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'D',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'F',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'G',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'H',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'J',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'K',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'L',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'Ş',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'İ',
						status: CharacterStatusNone,
					},
				},
			},
			&keyboardRow{
				keys: &[]*keyboardKey{
					&keyboardKey{
						char:   'Z',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'C',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'V',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'B',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'N',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'M',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'Ö',
						status: CharacterStatusNone,
					},
					&keyboardKey{
						char:   'Ç',
						status: CharacterStatusNone,
					},
				},
			},
		},
	}

	keysMap := make(map[rune]*keyboardKey)

	for _, r := range keyboard.rows {
		r.width = len(*r.keys)*(keyboard.boxW+keyboard.boxGap) - keyboard.boxGap

		for _, k := range *r.keys {
			keysMap[k.char] = k
		}
	}

	keyboard.keysMap = &keysMap

	return keyboard
}

func (k *keyboard) setKeyStatus(r rune, s CharacterStatus) {
	if key, exists := (*k.keysMap)[unicode.TurkishCase.ToUpper(r)]; exists {
		if key.status == CharacterStatusCorrectLocation {
			return
		}

		key.status = s
	}
}

func (k *keyboard) draw(screen *ebiten.Image) {
	y := 438
	for _, r := range k.rows {
		x := (screen.Bounds().Dx() - r.width) / 2
		for _, key := range *r.keys {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))

			keyboardButton := ebiten.NewImage(k.boxW, k.boxH)
			btnColor := key.status.getColor()
			if btnColor == color.White {
				btnColor = lightGrayColor
			}
			keyboardButton.Fill(btnColor)

			textOp := &text.DrawOptions{}
			textOp.GeoM.Translate(float64(keyboardButton.Bounds().Dx())/2, float64(keyboardButton.Bounds().Dy())/2)
			textOp.ColorScale.ScaleWithColor(color.White)
			textOp.PrimaryAlign = text.AlignCenter
			textOp.SecondaryAlign = text.AlignCenter
			text.Draw(keyboardButton, string(unicode.TurkishCase.ToUpper(key.char)), &text.GoTextFace{
				Source: mplusRegularTextFaceSource,
				Size:   float64(14),
			}, textOp)
			screen.DrawImage(keyboardButton, op)

			x += k.boxW + k.boxGap
		}

		y += k.boxH + k.boxGap
	}
}
