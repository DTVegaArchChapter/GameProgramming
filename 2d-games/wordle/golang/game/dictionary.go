package wordle

import (
	_ "embed"
	"math/rand"
	"sort"
	"strings"
)

//go:embed dict.txt
var dictTxt string

type dictionary struct {
	Words []string
}

func NewDictionary() *dictionary {
	d := &dictionary{}

	d.init()

	return d
}

func (d *dictionary) init() {
	d.Words = loadWords()
}

func (d *dictionary) GetRandomWord() string {
	return d.Words[rand.Intn(len(d.Words))]
}

func (d *dictionary) WordExists(w string) bool {
	wUpper := TurkishUpper.String(w)

	i := sort.SearchStrings(d.Words, wUpper)
	return i < len(d.Words) && d.Words[i] == wUpper
}

func loadWords() []string {
	w := strings.Split(strings.ReplaceAll(dictTxt, "\r", ""), "\n")
	toUpper(&w)
	sort.Strings(w)

	return w
}

func toUpper(arr *[]string) {
	l := len(*arr)

	for i := 0; i < l; i++ {
		(*arr)[i] = TurkishUpper.String((*arr)[i])
	}
}
