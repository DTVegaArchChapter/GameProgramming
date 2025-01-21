package wordle

import "image/color"

type CharacterStatus int

const (
	CharacterStatusNone CharacterStatus = iota
	CharacterStatusNotPresent
	CharacterStatusWrongLocation
	CharacterStatusCorrectLocation
)

func (s *CharacterStatus) getColor() color.Color {
	var r color.Color = color.White

	switch *s {
	case CharacterStatusWrongLocation:
		r = yellowColor
	case CharacterStatusNotPresent:
		r = grayColor
	case CharacterStatusCorrectLocation:
		r = greenColor
	}

	return r
}

func (s CharacterStatus) String() string {
	switch s {
	case CharacterStatusWrongLocation:
		return "CharacterStatusWrongLocation"
	case CharacterStatusNotPresent:
		return "CharacterStatusNotPresent"
	case CharacterStatusCorrectLocation:
		return "CharacterStatusCorrectLocation"
	default:
		return "unknown"
	}
}
