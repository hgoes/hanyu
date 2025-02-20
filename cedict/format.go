package cedict

import (
	"github.com/hgoes/hanyu/pinyin"
)

// Entry is a word definition. The word is given in both traditional
// and simplified writing. Since a word can have multiple meanings
// depending on context, a list of meanings is attached. The
// pronunciation of the word is given in the form of pinyin, where
// each pinyin corresponds to one character of the word.
type Entry struct {
	Traditional string
	Simplified  string
	Pinyin      []Pinyin
	Meaning     []string
}

// Pinyin can be either an encoded pinyin, or for certain special
// cases, a literal.
type Pinyin struct {
	Pinyin  pinyin.Pinyin
	Literal string
}

// Metadata can be attached to the dictionary for processing purposes.
type Metadata struct {
	Key   string
	Value string
}

// Comment represents a CEDICT comment line. Usually should be
// ignored.
type Comment string

// Line represents a line from a CEDICT dictionary and can either be
// an [Entry], [Metadata] or a [Comment].
type Line interface {
	isLine()
}

func (e Entry) isLine() {}

func (m Metadata) isLine() {}

func (c Comment) isLine() {}
