//go:generate go run ../cmd/gen-dict
package dict

import (
	_ "embed"
	"encoding/binary"
	"sort"
	"unicode"

	"github.com/hgoes/hanyu/cedict"
	"github.com/hgoes/hanyu/pinyin"
)

type Meaning struct {
	Pinyin      []cedict.Pinyin
	Meanings    []string
	HSKLevel    byte
	Simplified  string
	Traditional string
}

//go:embed gen.bin
var dict []byte

var Main = Dict{
	bin: dict,
}

type Dict struct {
	bin []byte
}

func uint24(data []byte) uint32 {
	return uint32(data[0])<<16 |
		uint32(data[1])<<8 |
		uint32(data[2])
}

func decodeRune(dict []byte, x uint16) rune {
	r := uint24(dict[3+int(x)*3:])
	return rune(r)
}

func (d *Dict) Lookup(
	str []rune,
) (int, []Meaning) {
	var lastWordLen int
	var softMatch, lastWord Lookup
	l := d.Begin()
	for i, c := range str {
		if !l.Consume(c) {
			break
		}
		if unicode.Is(unicode.Latin, c) &&
			i+1 < len(str) &&
			unicode.Is(unicode.Latin, str[i+1]) {
			// avoid matching partial latin words
			continue
		}
		if l.IsWord() {
			lastWordLen = i + 1
			lastWord = l
			if i == 0 && isSoftMatch(c) {
				softMatch = l
			}
		}
	}
	if !softMatch.IsZero() {
		l = d.Begin()
		for i, c := range str[1:] {
			if !l.Consume(c) {
				break
			}
			if unicode.Is(unicode.Latin, c) &&
				i+1 < len(str) &&
				unicode.Is(unicode.Latin, str[i+1]) {
				// avoid matching partial latin words
				continue
			}
			if l.IsWord() && i >= lastWordLen-1 {
				return 1, softMatch.Meanings(str[:1])
			}
		}
	}
	return lastWordLen, lastWord.Meanings(str[:lastWordLen])
}

func isSoftMatch(r rune) bool {
	switch r {
	case '不', '在', '有', '没':
		return true
	}
	return false
}

func (d *Dict) Begin() Lookup {
	// read the rune index length
	runeIdxLen := uint24(d.bin)
	return Lookup{
		dict:     d.bin,
		meanings: -1,
		index:    6 + 3*int(runeIdxLen),
	}
}

type Lookup struct {
	dict     []byte
	meanings int
	index    int
}

func (l *Lookup) IsZero() bool {
	return l.dict == nil
}

func (cur *Lookup) Consume(c rune) bool {
	// get the index length
	l := binary.BigEndian.Uint16(cur.dict[cur.index:])
	idx := sort.Search(int(l), func(p int) bool {
		// read the rune
		r := decodeRune(cur.dict, binary.BigEndian.Uint16(cur.dict[cur.index+2+p*5:]))
		return r >= c
	})
	if idx >= int(l) {
		return false
	}
	r := decodeRune(cur.dict, binary.BigEndian.Uint16(cur.dict[cur.index+2+idx*5:]))
	if r != c {
		return false
	}
	newMeanings := uint24(cur.dict[cur.index+2+idx*5+2:])
	lenMeanings := cur.dict[newMeanings]
	newIndex := int(newMeanings) + 1 + int(lenMeanings)*3
	cur.meanings = int(newMeanings)
	cur.index = newIndex
	return true
}

func (cur *Lookup) IsWord() bool {
	if cur.meanings == -1 || cur.dict == nil {
		return false
	}
	lenMeanings := cur.dict[cur.meanings]
	return lenMeanings > 0
}

func (cur *Lookup) Meanings(word []rune) []Meaning {
	if cur.meanings == -1 || cur.dict == nil {
		return nil
	}
	lenMeanings := cur.dict[cur.meanings]
	if lenMeanings == 0 {
		return nil
	}
	// read the rune index length
	runeIdxLen := uint24(cur.dict)
	// read the meanings offset
	offset := uint24(cur.dict[3+3*runeIdxLen:]) + 6 + 3*runeIdxLen

	meanings := make([]Meaning, lenMeanings)
	for i := range meanings {
		idx := uint24(cur.dict[cur.meanings+1+i*3:])
		meaningOffset := int(uint24(cur.dict[int(offset)+int(idx)*3:]))
		pos := int(offset) + meaningOffset
		hsk := cur.dict[pos]
		meanings[i].HSKLevel = hsk
		pos += 1
		pinSz := cur.dict[pos]
		pos += 1
		var pinyins []cedict.Pinyin
		if pinSz != 0 {
			pinyins = make([]cedict.Pinyin, pinSz)
		}
		for j := range pinyins {
			b := cur.dict[pos]
			if b > 127 {
				// it's a regular pinyin
				raw := binary.BigEndian.Uint16(cur.dict[pos:]) & 0x7FFF
				pinyins[j].Pinyin = pinyin.Pinyin(raw)
				pos += 2
			} else {
				// it's a literal pinyin
				pos++
				str := string(cur.dict[pos : pos+int(b)])
				pinyins[j].Literal = str
				pos += int(b)
			}
		}
		meanings[i].Pinyin = pinyins
		var means []string
		meansSize := cur.dict[pos]
		pos++
		if meansSize != 0 {
			means = make([]string, meansSize)
		}
		for j := range means {
			l := binary.BigEndian.Uint16(cur.dict[pos:])
			pos += 2
			means[j] = string(cur.dict[pos : pos+int(l)])
			pos += int(l)
		}
		meanings[i].Meanings = means
		varSize := cur.dict[pos]
		pos++
		var simp, trad []rune
		if varSize != 0 {
			simp = make([]rune, len(word))
			copy(simp, word)
			trad = make([]rune, len(word))
			copy(trad, word)
		}
		for j := byte(0); j < varSize; j++ {
			cpos := cur.dict[pos]
			pos++
			trad[cpos] = decodeRune(cur.dict, binary.BigEndian.Uint16(cur.dict[pos:]))
			pos += 2
			simp[cpos] = decodeRune(cur.dict, binary.BigEndian.Uint16(cur.dict[pos:]))
			pos += 2
		}
		meanings[i].Simplified = string(simp)
		meanings[i].Traditional = string(trad)
	}
	return meanings
}
