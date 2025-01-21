package wordle_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/DTVegaArchChapter/GameProgramming/wordle/game"
)

func TestNewDictionary(t *testing.T) {
	d := wordle.NewDictionary()

	if len(d.Words) == 0 {
		t.Error("dictionary word count should not be zero")
	}

	for i := 0; i < len(d.Words); i++ {
		if d.Words[i] != wordle.TurkishUpper.String(d.Words[i]) {
			t.Errorf("dictionary word %s is not upper-case", d.Words[i])
		}
	}

	if !sort.StringsAreSorted(d.Words) {
		t.Error("dictionary words are not sorted")
	}
}

func TestGetRandomWord(t *testing.T) {
	d := wordle.NewDictionary()
	w := d.GetRandomWord()

	if w == "" {
		t.Error("dictionary random word cannot be empty")
	}
}

func TestWordExists(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"apple", false},
		{"APPLE", false},
		{"ayakkabı", false},
		{"ayakkabı", false},
		{"Arşiv", true},
		{"arşiv", true},
		{"ARŞİV", true},
		{"ArŞiv", true},
	}

	d := wordle.NewDictionary()

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("dictionary exists %s", tc.input), func(t *testing.T) {
			actual := d.WordExists(tc.input)
			if actual != tc.expected {
				t.Errorf("wordExists('%s') expected=%v actual=%v", tc.input, tc.expected, actual)
			}
		})
	}
}
