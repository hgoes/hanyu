//go:generate go run ../cmd/gen-pinyin-parser
package pinyin

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Pinyin uint16

type parserState int

type Sound uint16

type Tone byte

const (
	Neutral Tone = 0
	Flat    Tone = 1
	Rising  Tone = 2
	Low     Tone = 3
	Falling Tone = 4
)

func New(p Sound, t Tone) Pinyin {
	return Pinyin(uint16(t) + uint16(p)*5)
}

func (p Pinyin) Decode() (Sound, Tone) {
	return Sound(p / 5), Tone(p % 5)
}

func (p Pinyin) String() string {
	str, _ := p.render()
	return str
}

func tonePosition(rs []rune) int {
	var wasO bool
	firstVowel := -1
	for i, r := range rs {
		switch r {
		case 'a', 'e':
			return i
		case 'o':
			if firstVowel != -1 {
				return i
			}
			wasO = true
			firstVowel = i
		case 'u':
			if wasO {
				return i - 1
			}
			if firstVowel != -1 {
				return i
			}
			firstVowel = i
		case 'i', 'ü':
			if firstVowel != -1 {
				return i
			}
			firstVowel = i
		default:
			if firstVowel != -1 {
				return firstVowel
			}
		}
	}
	if firstVowel != -1 {
		return firstVowel
	}
	return 0
}

func (p Pinyin) render() (string, bool) {
	c, t := p.Decode()
	str := c.String()
	if t == Neutral {
		return str, c.special()
	}
	runes := []rune(str)
	tonePos := tonePosition(runes)
	var repl rune
	switch runes[tonePos] {
	case 'a':
		switch t {
		case Flat:
			repl = 'ā'
		case Rising:
			repl = 'á'
		case Low:
			repl = 'ǎ'
		case Falling:
			repl = 'à'
		}
	case 'e':
		switch t {
		case Flat:
			repl = 'ē'
		case Rising:
			repl = 'é'
		case Low:
			repl = 'ě'
		case Falling:
			repl = 'è'
		}
	case 'i':
		switch t {
		case Flat:
			repl = 'ī'
		case Rising:
			repl = 'í'
		case Low:
			repl = 'ǐ'
		case Falling:
			repl = 'ì'
		}
	case 'o':
		switch t {
		case Flat:
			repl = 'ō'
		case Rising:
			repl = 'ó'
		case Low:
			repl = 'ǒ'
		case Falling:
			repl = 'ò'
		}
	case 'u':
		switch t {
		case Flat:
			repl = 'ū'
		case Rising:
			repl = 'ú'
		case Low:
			repl = 'ǔ'
		case Falling:
			repl = 'ù'
		}
	case 'ü':
		switch t {
		case Flat:
			repl = 'ǖ'
		case Rising:
			repl = 'ǘ'
		case Low:
			repl = 'ǚ'
		case Falling:
			repl = 'ǜ'
		}
	case 'n':
		switch t {
		case Flat:
			repl = 'n'
		case Rising:
			repl = 'ń'
		case Low:
			repl = 'ň'
		case Falling:
			repl = 'ǹ'
		}
	}
	if repl != 0 {
		runes[tonePos] = repl
		return string(runes), c.special()
	}
	return "?", false
}

type parser struct {
	State parserState
	Tone  Tone
	Done  bool
}

func (p *parser) setTone(t Tone) bool {
	if p.Tone != 0 {
		return false
	}
	p.Tone = t
	return true
}

func (p *parser) Advance(r rune) bool {
	if p.Done {
		return false
	}
	r = unicode.ToLower(r)
	switch r {
	case '\'':
		p.Done = true
		return true
	case 'ā':
		if !p.setTone(Flat) {
			return false
		}
		r = 'a'
	case 'á':
		if !p.setTone(Rising) {
			return false
		}
		r = 'a'
	case 'ǎ':
		if !p.setTone(Low) {
			return false
		}
		r = 'a'
	case 'à':
		if !p.setTone(Falling) {
			return false
		}
		r = 'a'
	case 'ē':
		if !p.setTone(Flat) {
			return false
		}
		r = 'e'
	case 'é':
		if !p.setTone(Rising) {
			return false
		}
		r = 'e'
	case 'ě':
		if !p.setTone(Low) {
			return false
		}
		r = 'e'
	case 'è':
		if !p.setTone(Falling) {
			return false
		}
		r = 'e'
	case 'ī':
		if !p.setTone(Flat) {
			return false
		}
		r = 'i'
	case 'í':
		if !p.setTone(Rising) {
			return false
		}
		r = 'i'
	case 'ǐ':
		if !p.setTone(Low) {
			return false
		}
		r = 'i'
	case 'ì':
		if !p.setTone(Falling) {
			return false
		}
		r = 'i'
	case 'ō':
		if !p.setTone(Flat) {
			return false
		}
		r = 'o'
	case 'ó':
		if !p.setTone(Rising) {
			return false
		}
		r = 'o'
	case 'ǒ':
		if !p.setTone(Low) {
			return false
		}
		r = 'o'
	case 'ò':
		if !p.setTone(Falling) {
			return false
		}
		r = 'o'
	case 'ū':
		if !p.setTone(Flat) {
			return false
		}
		r = 'u'
	case 'ú':
		if !p.setTone(Rising) {
			return false
		}
		r = 'u'
	case 'ǔ':
		if !p.setTone(Low) {
			return false
		}
		r = 'u'
	case 'ù':
		if !p.setTone(Falling) {
			return false
		}
		r = 'u'
	case 'ǖ':
		if !p.setTone(Flat) {
			return false
		}
		r = 'ü'
	case 'ǘ':
		if !p.setTone(Rising) {
			return false
		}
		r = 'ü'
	case 'ǚ':
		if !p.setTone(Low) {
			return false
		}
		r = 'ü'
	case 'ǜ':
		if !p.setTone(Falling) {
			return false
		}
		r = 'ü'
	case '1':
		p.Done = true
		return p.setTone(Flat)
	case '2':
		p.Done = true
		return p.setTone(Rising)
	case '3':
		p.Done = true
		return p.setTone(Low)
	case '4':
		p.Done = true
		return p.setTone(Falling)
	case '5':
		p.Done = true
		return p.setTone(Neutral)
	case 'ń':
		if !p.setTone(Rising) {
			return false
		}
		r = 'n'
	case 'ň':
		if !p.setTone(Low) {
			return false
		}
		r = 'n'
	case 'ǹ':
		if !p.setTone(Falling) {
			return false
		}
		r = 'n'
	case 'ḿ':
		if !p.setTone(Rising) {
			return false
		}
		r = 'm'
	}
	ok, nxt := p.State.next(r)
	if !ok {
		return false
	}
	p.State = nxt
	return true
}

func (p *parser) Result() (bool, Pinyin) {
	ok, sound := p.State.Sound()
	if !ok {
		return false, 0
	}
	return true, New(sound, p.Tone)
}

func Parse(str []rune) (bool, Pinyin, []rune) {
	var p parser
	var lastIdx int
	var lastResult Pinyin
	for i, r := range str {
		ok := p.Advance(r)
		if !ok {
			break
		}
		ok, result := p.Result()
		if ok {
			lastIdx = i + 1
			lastResult = result
		}
	}
	if lastIdx == 0 {
		return false, 0, nil
	}
	return true, lastResult, str[lastIdx:]
}

func ParseMany(str []rune) (bool, []Pinyin, []rune) {
	var result []Pinyin
	for {
		if len(str) == 0 {
			return true, result, nil
		}
		ok, p, rest := Parse(str)
		if !ok {
			return false, result, rest
		}
		result = append(result, p)
		str = rest
	}
}

func RenderMany(ps []Pinyin) string {
	var buf strings.Builder
	RenderManyWriter(&buf, ps)
	return buf.String()
}

func RenderManyWriter(w io.Writer, ps []Pinyin) (int, error) {
	c := 0
	for i, p := range ps {
		str, special := p.render()
		if i != 0 && special {
			n, err := fmt.Fprint(w, "'")
			c += n
			if err != nil {
				return c, err
			}
		}
		n, err := fmt.Fprint(w, str)
		c += n
		if err != nil {
			return c, err
		}
	}
	return c, nil
}
