package numbers

import "testing"

func Test(t *testing.T) {
	tests := []struct {
		Input string
		Value int64
	}{
		{
			Input: "五",
			Value: 5,
		},
		{
			Input: "一百六十八",
			Value: 168,
		},
		{
			Input: "六十",
			Value: 60,
		},
		{
			Input: "二十",
			Value: 20,
		},
		{
			Input: "兩百",
			Value: 200,
		},
		{
			Input: "二千",
			Value: 2000,
		},
		{
			Input: "四十五",
			Value: 45,
		},
		{
			Input: "兩千三百六十二",
			Value: 2362,
		},
		{
			Input: "十四",
			Value: 14,
		},
		{
			Input: "一萬兩千",
			Value: 12000,
		},
		{
			Input: "一百一十四",
			Value: 114,
		},
		{
			Input: "一千一百五十八",
			Value: 1158,
		},
		{
			Input: "十二兆三千四百五十六亿七千八百九十万二千三百四十五",
			Value: 12345678902345,
		},
		{
			Input: "二百零五",
			Value: 205,
		},
		{
			Input: "十萬零四",
			Value: 100004,
		},
		{
			Input: "一千零五萬二十六",
			Value: 10050026,
		},
		{
			Input: "六零一二七",
			Value: 60127,
		},
	}
	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			var p Parser
			for _, c := range []rune(test.Input) {
				if !p.Consume(c) {
					t.Fatal("failed to parse")
				}
				t.Log(p.Value())
			}
			if actual := p.Value(); actual != test.Value {
				t.Error("wrong value:", actual)
			}
		})
	}
}
