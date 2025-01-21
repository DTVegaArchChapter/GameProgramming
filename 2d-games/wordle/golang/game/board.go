package wordle

import (
	"strings"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
)

type gameState int

const (
	gameInProgress gameState = iota
	gameLost
	gameWon
)

type board struct {
	rows         int
	cols         int
	alphabet     string
	dict         *dictionary
	answer       []rune
	tiles        []*tile
	pos          int
	wordNotFound bool
	keyboard     *keyboard
	state        gameState
	inputRune    rune
	inputKey     ebiten.Key
}

func newBoard(keyboard *keyboard) *board {
	b := &board{
		rows:      6,
		cols:      5,
		alphabet:  "ABCÇDEFGĞHIİJKLMNOÖPRSŞTUÜVYZ",
		dict:      NewDictionary(),
		keyboard:  keyboard,
		inputRune: 0,
		inputKey:  0,
	}

	b.init()

	return b
}

func (b *board) Update() {
	if b.inputKey == ebiten.KeyBackspace {
		b.deleteCurrentChar()
	} else if b.inputKey == ebiten.KeyEnter || b.inputKey == ebiten.KeyNumpadEnter {
		b.checkCurrentWord()
	} else if b.inputRune > 0 {
		b.addChar(b.inputRune)
	}

	for _, t := range b.tiles {
		t.Update()
	}
}

func (b *board) Draw(screen *ebiten.Image) {
	for _, t := range b.tiles {
		t.Draw(screen)
	}
}

func (b *board) deleteCurrentChar() bool {
	if b.pos <= 0 {
		return false
	}

	t := b.tiles[b.pos]
	if b.isPosInLastChar() && !t.isEmpty() {
		if t.isCharStatusNone() {
			t.clearRune()
			b.wordNotFound = false

			return true
		}
	} else {
		tPrev := b.tiles[b.pos-1]
		if tPrev.isCharStatusNone() {
			tPrev.clearRune()
			b.wordNotFound = false
			b.pos--

			return true
		}
	}

	return false
}

func (b *board) checkCurrentWord() {
	if !b.isPosInLastChar() {
		return
	}

	answer := make([]rune, b.cols)
	for i, c := b.pos-b.cols+1, 0; i < b.pos+1; i++ {
		t := b.tiles[i]
		answer[c] = t.r
		c++
	}

	if !b.dict.WordExists(string(answer)) {
		b.wordNotFound = true

		for i := b.pos - b.cols + 1; i < b.pos+1; i++ {
			b.tiles[i].startShake()
		}

		return
	}

	won := true

	checkResult := CheckAnswerRunes(answer, b.answer)

	for i, c := b.pos-b.cols+1, 0; i < b.pos+1; i++ {
		t := b.tiles[i]
		s := checkResult[c]

		t.setStatus(s)
		b.keyboard.setKeyStatus(t.r, t.status)

		if s != CharacterStatusCorrectLocation {
			won = false
		}

		t.flip()

		c++
	}

	if won {
		b.state = gameWon
	} else if b.pos >= len(b.tiles)-1 {
		b.state = gameLost
	}

	if b.pos < len(b.tiles)-1 {
		b.pos++
	}
}

func (b *board) addChar(r rune) bool {
	if b.pos >= len(b.tiles) {
		return false
	}

	t := b.tiles[b.pos]

	if !(b.isPosInLastChar() && !t.isEmpty()) && strings.ContainsRune(b.alphabet, unicode.TurkishCase.ToUpper(r)) {
		t.setRune(r)

		if !b.isPosInLastChar() {
			b.pos++
		}

		return true
	}

	return false
}

func (b *board) init() {
	tiles := make([]*tile, b.cols*b.rows)
	for col := 0; col < b.cols; col++ {
		for row := 0; row < b.rows; row++ {
			tiles[b.calcPos(col, row)] = newTile(col, row, b)
		}
	}

	for i := 0; i < len(tiles); i++ {
		tiles[i].clearValue()
	}

	b.state = gameInProgress
	b.pos = 0
	b.tiles = tiles
	b.answer = []rune(b.dict.GetRandomWord())
}

func (b *board) isPosInLastChar() bool {
	return b.pos%b.cols == b.cols-1
}

func (b *board) calcPos(col, row int) int {
	return (row * b.cols) + col
}

func CheckAnswerRunes(answer, correct []rune) []CharacterStatus {
	result := make([]CharacterStatus, len(answer))
	l := []rune{}

	for i := 0; i < len(answer); i++ {
		a := answer[i]
		c := correct[i]
		if unicode.TurkishCase.ToUpper(a) == unicode.TurkishCase.ToUpper(c) {
			result[i] = CharacterStatusCorrectLocation
		} else {
			l = append(l, c)
		}
	}

	for i := 0; i < len(answer); i++ {
		r := result[i]

		if r == CharacterStatusCorrectLocation {
			continue
		}

		a := answer[i]

		if contains(l, a) {
			result[i] = CharacterStatusWrongLocation

			l = removeRune(l, a)
		} else {
			result[i] = CharacterStatusNotPresent
		}
	}

	return result
}

func removeRune(l []rune, r rune) []rune {
	for i := 0; i < len(l); i++ {
		if unicode.TurkishCase.ToUpper(l[i]) == unicode.TurkishCase.ToUpper(r) {
			return append(l[:i], l[i+1:]...)
		}
	}

	return l
}
