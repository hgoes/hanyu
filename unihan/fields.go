package unihan

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hgoes/hanyu/pinyin"
)

func getField(tp FieldType) Field {
	switch tp {
	case AccountingNumeric:
		var n AccountingNumericF
		return &n
	case GradeLevel:
		var g GradeLevelF
		return &g
	case HanyuPinyin:
		var h HanyuPinyinF
		return &h
	case Mandarin:
		return &MandarinF{}
	case PrimaryNumeric:
		var n PrimaryNumericF
		return &n
	case OtherNumeric:
		var n OtherNumericF
		return &n
	case Definition:
		var d DefinitionF
		return &d
	case SimplifiedVariant:
		var v SimplifiedVariantF
		return &v
	case TraditionalVariant:
		var v TraditionalVariantF
		return &v
	default:
		return &Generic{tp: tp}
	}
}

type Generic struct {
	tp      FieldType
	Content string
}

func (g *Generic) Type() FieldType {
	return g.tp
}

func (g *Generic) Parse(c string) error {
	g.Content = c
	return nil
}

type AccountingNumericF int64

func (f *AccountingNumericF) Type() FieldType {
	return AccountingNumeric
}

func (f *AccountingNumericF) Parse(c string) (err error) {
	*(*int64)(f), err = strconv.ParseInt(c, 10, 64)
	return
}

func (f *AccountingNumericF) String() string {
	return strconv.FormatInt(int64(*f), 10)
}

type GradeLevelF int8

func (_ *GradeLevelF) Type() FieldType {
	return GradeLevel
}

func (f *GradeLevelF) Parse(c string) error {
	grade, err := strconv.ParseInt(c, 10, 8)
	if err != nil {
		return err
	}
	*f = GradeLevelF(grade)
	return nil
}

func (f *GradeLevelF) String() string {
	return strconv.FormatInt(int64(*f), 10)
}

type HanyuPinyinF []HanyuPinyinReading

type HanyuPinyinReading struct {
	Locations []string
	Pinyin    []pinyin.Pinyin
	PinyinRaw []string
}

func (_ *HanyuPinyinF) Type() FieldType {
	return HanyuPinyin
}

func (f *HanyuPinyinF) Parse(c string) error {
	var readings []HanyuPinyinReading

	for {
		if c == "" {
			break
		}
		locsRaw, rest, ok := strings.Cut(c, ":")
		locsRaw = strings.TrimSpace(locsRaw)
		if !ok {
			return fmt.Errorf("invalid reading: %q", c)
		}
		locs := strings.Split(locsRaw, ",")
		pinyinRaw, rest, _ := strings.Cut(rest, " ")
		pinyinsRaw := strings.Split(pinyinRaw, ",")
		pinyins := make([]pinyin.Pinyin, len(pinyinsRaw))
		for i, p := range pinyinsRaw {
			ok, result, rest := pinyin.Parse([]rune(p))
			if !ok || len(rest) != 0 {
				pinyins = nil
				break
			}
			pinyins[i] = result
		}
		readings = append(readings, HanyuPinyinReading{
			Locations: locs,
			Pinyin:    pinyins,
			PinyinRaw: pinyinsRaw,
		})
		c = rest
	}
	*f = HanyuPinyinF(readings)
	return nil
}

type MandarinF struct {
	ReadingCN pinyin.Pinyin
	ReadingTW pinyin.Pinyin
}

func (_ *MandarinF) Type() FieldType {
	return Mandarin
}

func (f *MandarinF) Parse(c string) error {
	ok, pin, rest := pinyin.Parse([]rune(c))
	if !ok {
		return fmt.Errorf("invalid pinyin: %q", c)
	}
	f.ReadingCN = pin
	if len(rest) == 0 {
		f.ReadingTW = pin
		return nil
	}
	ok, pin, rest = pinyin.Parse(rest[1:])
	if !ok || len(rest) != 0 {
		return fmt.Errorf("invalid pinyin: %q", c)
	}
	f.ReadingTW = pin
	return nil
}

type PrimaryNumericF int64

func (_ *PrimaryNumericF) Type() FieldType {
	return PrimaryNumeric
}

func (f *PrimaryNumericF) Parse(c string) (err error) {
	*(*int64)(f), err = strconv.ParseInt(c, 10, 64)
	return
}

func (f *PrimaryNumericF) String() string {
	return strconv.FormatInt(int64(*f), 10)
}

type OtherNumericF int64

func (_ *OtherNumericF) Type() FieldType {
	return OtherNumeric
}

func (f *OtherNumericF) Parse(c string) (err error) {
	*(*int64)(f), err = strconv.ParseInt(c, 10, 64)
	return
}

func (f *OtherNumericF) String() string {
	return strconv.FormatInt(int64(*f), 10)
}

type DefinitionF string

func (_ *DefinitionF) Type() FieldType {
	return Definition
}

func (f *DefinitionF) Parse(c string) error {
	*f = DefinitionF(c)
	return nil
}

func (f *DefinitionF) String() string {
	return string(*f)
}

type SimplifiedVariantF []rune

func (_ *SimplifiedVariantF) Type() FieldType {
	return SimplifiedVariant
}

func (f *SimplifiedVariantF) Parse(c string) error {
	runes, err := parseRunes(c)
	if err != nil {
		return err
	}
	*f = runes
	return nil
}

func (f *SimplifiedVariantF) String() string {
	return string(*f)
}

type TraditionalVariantF []rune

func (_ *TraditionalVariantF) Type() FieldType {
	return TraditionalVariant
}

func (f *TraditionalVariantF) Parse(c string) error {
	runes, err := parseRunes(c)
	if err != nil {
		return err
	}
	*f = runes
	return nil
}

func (f *TraditionalVariantF) String() string {
	return string(*f)
}

func parseRunes(c string) ([]rune, error) {
	chars := strings.Split(c, " ")
	result := make([]rune, len(chars))
	for i, char := range chars {
		numStr, ok := strings.CutPrefix(char, "U+")
		if !ok {
			return nil, fmt.Errorf("invalid character %q", char)
		}
		num, err := strconv.ParseInt(numStr, 16, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid character %q", c)
		}
		result[i] = rune(num)
	}
	return result, nil
}
