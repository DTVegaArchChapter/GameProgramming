package wordle

import (
	"image/color"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
)

type keyboard struct {
	rows    *[]*keyboardRow
	keysMap *map[rune]*keyboardKey
	text    *TextRenderer
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
		text:   NewTextRenderer(RobotoBoldFontName, color.White, 18),
		rows:   createKeyboardKeyRows("ERTYUIOPĞÜ", "ASDFGHJKLŞİ", "ZCVBNMÖÇ"),
	}

	keysMap := make(map[rune]*keyboardKey)

	for _, r := range *keyboard.rows {
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

		if key.status == CharacterStatusWrongLocation && s != CharacterStatusCorrectLocation {
			return
		}

		key.status = s
	}
}

func (k *keyboard) draw(screen *ebiten.Image) {
	y := 438
	for _, r := range *k.rows {
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

			k.text.Draw(keyboardButton, string(unicode.TurkishCase.ToUpper(key.char)), keyboardButton.Bounds().Dx()/2, keyboardButton.Bounds().Dy()/2)
			screen.DrawImage(keyboardButton, op)

			x += k.boxW + k.boxGap
		}

		y += k.boxH + k.boxGap
	}
}

func createKeyboardKeyRows(rows ...string) *[]*keyboardRow {
	result := make([]*keyboardRow, len(rows))

	for i, r := range rows {
		result[i] = &keyboardRow{
			keys: createKeyboardKeyArray(r),
		}
	}

	return &result
}

func createKeyboardKeyArray(keys string) *[]*keyboardKey {
	runeArray := []rune(keys)
	result := make([]*keyboardKey, len(runeArray))

	for i, r := range runeArray {
		result[i] = &keyboardKey{
			char:   r,
			status: CharacterStatusNone,
		}
	}

	return &result
}
