package cedict

import (
	"github.com/hgoes/hanyu/pinyin"
)

type Entry struct {
	Traditional string
	Simplified  string
	Pinyin      []Pinyin
	Meaning     []string
}

type Pinyin struct {
	Pinyin  pinyin.Pinyin
	Literal string
}

type Metadata struct {
	Key   string
	Value string
}

type Comment string

type Line interface {
	isLine()
}

func (e Entry) isLine() {}

func (m Metadata) isLine() {}

func (c Comment) isLine() {}
