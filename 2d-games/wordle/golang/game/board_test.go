package wordle_test

import (
	"fmt"
	"testing"

	"github.com/DTVegaArchChapter/GameProgramming/wordle/game"
)

func TestCheckAnswerRunes(t *testing.T) {
	testCases := []struct {
		answer   []rune
		correct  []rune
		expected []wordle.CharacterStatus
	}{
		{
			answer:  []rune{'A', 'A', 'B', 'C', 'D'},
			correct: []rune{'A', 'F', 'G', 'B', 'B'},
			expected: []wordle.CharacterStatus{
				wordle.CharacterStatusCorrectLocation,
				wordle.CharacterStatusNotPresent,
				wordle.CharacterStatusWrongLocation,
				wordle.CharacterStatusNotPresent,
				wordle.CharacterStatusNotPresent,
			},
		},
		{
			answer:  []rune{'A', 'A', 'C', 'B', 'B'},
			correct: []rune{'B', 'C', 'A', 'B', 'B'},
			expected: []wordle.CharacterStatus{
				wordle.CharacterStatusWrongLocation,
				wordle.CharacterStatusNotPresent,
				wordle.CharacterStatusWrongLocation,
				wordle.CharacterStatusCorrectLocation,
				wordle.CharacterStatusCorrectLocation,
			},
		},
		{
			answer:  []rune{'A', 'B', 'C', 'D', 'E'},
			correct: []rune{'F', 'G', 'H', 'I', 'İ'},
			expected: []wordle.CharacterStatus{
				wordle.CharacterStatusNotPresent,
				wordle.CharacterStatusNotPresent,
				wordle.CharacterStatusNotPresent,
				wordle.CharacterStatusNotPresent,
				wordle.CharacterStatusNotPresent,
			},
		},
		{
			answer:  []rune{'A', 'B', 'C', 'D', 'E'},
			correct: []rune{'E', 'D', 'C', 'B', 'A'},
			expected: []wordle.CharacterStatus{
				wordle.CharacterStatusWrongLocation,
				wordle.CharacterStatusWrongLocation,
				wordle.CharacterStatusCorrectLocation,
				wordle.CharacterStatusWrongLocation,
				wordle.CharacterStatusWrongLocation,
			},
		},
		{
			answer:  []rune{'A', 'A', 'A', 'A', 'A'},
			correct: []rune{'A', 'B', 'C', 'A', 'A'},
			expected: []wordle.CharacterStatus{
				wordle.CharacterStatusCorrectLocation,
				wordle.CharacterStatusNotPresent,
				wordle.CharacterStatusNotPresent,
				wordle.CharacterStatusCorrectLocation,
				wordle.CharacterStatusCorrectLocation,
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		tc := testCases[i]

		t.Run(fmt.Sprintf("Answer: %c, Correct: %c", tc.answer, tc.correct), func(t *testing.T) {
			actual := wordle.CheckAnswerRunes(tc.answer, tc.correct)
			if !assertAreEqual(tc.expected, actual) {
				t.Errorf("\nExpected: %v\n Actual: %v", tc.expected, actual)
			}
		})
	}
}

func assertAreEqual(v1, v2 []wordle.CharacterStatus) bool {
	if len(v1) != len(v2) {
		return false
	}

	for i := 0; i < len(v1); i++ {
		if v1[i] != v2[i] {
			return false
		}
	}

	return true
}
