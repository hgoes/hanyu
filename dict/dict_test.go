package dict

import (
	"strings"
	"testing"

	"github.com/hgoes/hanyu/pinyin"
)

func Test(t *testing.T) {
	l, m := Main.Lookup([]rune("不可胜数"))
	if l != 4 {
		t.Fatal("wrong length:", l)
	}
	if l := len(m); l != 1 {
		t.Fatal("wrong length of meanings:", l)
	}
	if l := len(m[0].Pinyin); l != 4 {
		t.Fatal("wrong length of pinyins:", l)
	}
	p := make([]pinyin.Pinyin, 4)
	for i, pin := range m[0].Pinyin {
		if pin.Literal != "" {
			t.Fatal("unexpected literal:", pin.Literal)
		}
		p[i] = pin.Pinyin
	}
	if str := pinyin.RenderMany(p); str != "bùkěshèngshǔ" {
		t.Fatal("unexpected pinyin:", str)
	}
	if l := len(m[0].Meanings); l != 2 {
		t.Fatal("wrong length of meanings:", l)
	}
	if mean := m[0].Meanings[0]; mean != "countless" {
		t.Fatal("wrong meaning:", mean)
	}
}

func TestLookup(t *testing.T) {
	tests := []struct {
		Input   string
		Meaning string
	}{
		{
			Input:   "A",
			Meaning: "(slang) (Tw) to steal/",
		},
		{
			Input:   "A B",
			Meaning: "(slang) (Tw) to steal/",
		},
		{
			Input: "Ain’t",
		},
		{
			Input:   "不复杂",
			Meaning: "no; not so/(bound form) not; un-/",
		},
		{
			Input:   "做",
			Meaning: "to make; to produce/to write; to compose/to do; to engage in; to hold (a party etc)/(of a person) to be (an intermediary, a good student etc); to become (husband and wife, friends etc)/(of a thing) to serve as; to be used for/to assume (an air or manner)/",
		},
		{
			Input:   "有时候",
			Meaning: "sometimes/",
		},
		{
			Input:   "在那儿",
			Meaning: "(adverbial expression indicating that the attention of the subject of the verb is focused on what they are doing, not distracted by anything else)/just ...ing (and nothing else)/",
		},
	}
	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			_, m := Main.Lookup([]rune(test.Input))
			if len(m) == 0 {
				if test.Meaning != "" {
					t.Error("expected a match")
				}
				return
			}
			var buf strings.Builder
			for _, mean := range m {
				for _, meaning := range mean.Meanings {
					buf.WriteString(meaning)
					buf.WriteString("/")
				}
			}
			if test.Meaning == "" {
				t.Fatalf("expected no match, got %q", buf.String())
			}
			if actual := buf.String(); actual != test.Meaning {
				t.Errorf("wrong meaning: %q", actual)
			}
		})
	}
}

func TestHSK(t *testing.T) {
	tests := []struct {
		Input string
		HSK   byte
	}{
		{
			Input: "要",
			HSK:   2,
		},
		{
			Input: "半途而废",
			HSK:   6,
		},
		{
			Input: "一呼百应",
		},
	}
	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			rs := []rune(test.Input)
			l, mean := Main.Lookup(rs)
			if l != len(rs) {
				t.Fatal("not found")
			}
			if len(mean) < 1 {
				t.Fatal("no meaning")
			}
			if lvl := mean[0].HSKLevel; lvl != test.HSK {
				t.Error("wrong HSK level:", lvl)
			}
		})
	}
}

func TestPinyin(t *testing.T) {
	tests := []struct {
		Input  string
		Output string
	}{
		{
			Input:  "想听你听过的音乐",
			Output: "xiǎng [tīng|yǐn] [nǐ|nǐ] [tīng|yǐn] [guò|guō|guo] [de|dī|dí|dì] yīnyuè",
		},
	}
	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			cur := []rune(test.Input)
			var buf strings.Builder
			i := 0
			for len(cur) > 0 {
				if i != 0 {
					buf.WriteRune(' ')
				}
				i++
				l, meaning := Main.Lookup(cur)
				if len(meaning) != 0 {
					if len(meaning) > 1 {
						buf.WriteRune('[')
					}
					for j, m := range meaning {
						if j != 0 {
							buf.WriteRune('|')
						}
						for _, p := range m.Pinyin {

							if p.Literal != "" {
								buf.WriteString(p.Literal)
							} else {
								buf.WriteString(p.Pinyin.String())
							}
						}
					}
					if len(meaning) > 1 {
						buf.WriteRune(']')
					}
					cur = cur[l:]
					continue
				}
				buf.WriteRune(cur[0])
				cur = cur[1:]
			}
			result := buf.String()
			if result != test.Output {
				t.Error("wrong output:", result)
			}
		})
	}
}

func TestVariant(t *testing.T) {
	tests := []struct {
		Input       string
		Simplified  string
		Traditional string
	}{
		{
			Input:       "不拉幾",
			Traditional: "不拉幾",
			Simplified:  "不拉几",
		},
	}
	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			inp := []rune(test.Input)
			l, m := Main.Lookup(inp)
			if l != len(inp) {
				t.Fatal("not found")
			}
			if len(m) != 1 {
				t.Fatal("no meaning")
			}
			if x := m[0].Simplified; x != test.Simplified {
				t.Error("wrong simplified:", x)
			}
			if x := m[0].Traditional; x != test.Traditional {
				t.Error("wrong traditional:", x)
			}
		})
	}
}
