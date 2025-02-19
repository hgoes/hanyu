package pinyin

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		Input  string
		Output string
	}{
		{
			Input:  "a",
			Output: "a",
		},
		{
			Input:  "a3",
			Output: "ǎ",
		},
		{
			Input:  "tuan4",
			Output: "tuàn",
		},
		{
			Input:  "tiaō",
			Output: "tiāo",
		},
		{
			Input:  "n2",
			Output: "ń",
		},
		{
			Input:  "lu:3",
			Output: "lǚ",
		},
		{
			Input:  "pie1",
			Output: "piē",
		},
		{
			Input:  "mei3",
			Output: "měi",
		},
		{
			Input:  "r5",
			Output: "r",
		},
	}
	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			ok, result, rest := Parse([]rune(test.Input))
			if !ok || len(rest) != 0 {
				t.Fatal("failed to parse")
			}
			if result.String() != test.Output {
				t.Errorf("wrong result: %+v", result)
			}
		})
	}
}

func TestParseMany(t *testing.T) {
	tests := []struct {
		Input  string
		Output string
	}{
		{
			Input:  "nü'ér",
			Output: "nü'ér",
		},
	}
	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			ok, result, rest := ParseMany([]rune(test.Input))
			if !ok || len(rest) != 0 {
				t.Fatal("failed to parse")
			}
			if RenderMany(result) != test.Output {
				t.Errorf("wrong result: %+v", result)
			}
		})
	}
}
