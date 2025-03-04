package wordle

import (
	_ "embed"
	"math/rand"
	"sort"
	"strings"
)

//go:embed dict.txt
var dictTxt string

//go:embed target.txt
var targetTxt string

type dictionary struct {
	Words   []string
	targets []string
}

func NewDictionary() *dictionary {
	d := &dictionary{}

	d.init()

	return d
}

func (d *dictionary) init() {
	d.Words = loadWords(dictTxt)
	d.targets = loadWords(targetTxt)
}

func (d *dictionary) GetRandomWord() string {
	return d.targets[rand.Intn(len(d.targets))]
}

func (d *dictionary) WordExists(w string) bool {
	wUpper := TurkishUpper.String(w)

	i := sort.SearchStrings(d.Words, wUpper)
	return i < len(d.Words) && d.Words[i] == wUpper
}

func loadWords(target string) []string {
	w := strings.Split(strings.ReplaceAll(target, "\r", ""), "\n")
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
