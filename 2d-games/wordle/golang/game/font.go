package wordle

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	arcadeTextFaceSource       *text.GoTextFaceSource
	mplusRegularTextFaceSource *text.GoTextFaceSource
	fontSize                   int = 24
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}

	arcadeTextFaceSource = s

	s, err = text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusRegularTextFaceSource = s
}
