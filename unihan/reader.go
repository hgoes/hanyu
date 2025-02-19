package unihan

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Reader struct {
	rd                  io.Closer
	dictionaryIndices   *zip.File
	dictionaryLikeData  *zip.File
	irgSources          *zip.File
	numericValues       *zip.File
	otherMappins        *zip.File
	radicalStrokeCounts *zip.File
	readings            *zip.File
	variants            *zip.File
}

func Open(filename string) (*Reader, error) {
	rd, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	result := &Reader{
		rd: rd,
	}
	for _, file := range rd.File {
		switch file.Name {
		case "Unihan_DictionaryIndices.txt":
			result.dictionaryIndices = file
		case "Unihan_DictionaryLikeData.txt":
			result.dictionaryLikeData = file
		case "Unihan_IRGSources.txt":
			result.irgSources = file
		case "Unihan_NumericValues.txt":
			result.numericValues = file
		case "Unihan_OtherMappings.txt":
			result.otherMappins = file
		case "Unihan_RadicalStrokeCounts.txt":
			result.radicalStrokeCounts = file
		case "Unihan_Readings.txt":
			result.readings = file
		case "Unihan_Variants.txt":
			result.variants = file
		default:
			panic(file.Name)
		}
	}
	return result, nil
}

func (r *Reader) Close() error {
	return r.rd.Close()
}

func (r *Reader) All() *Entries {
	return &Entries{
		files: []*zip.File{
			r.dictionaryIndices,
			r.dictionaryLikeData,
			r.irgSources,
			r.numericValues,
			r.otherMappins,
			r.radicalStrokeCounts,
			r.readings,
			r.variants,
		},
	}
}

func (r *Reader) Get(tps ...FieldType) *Entries {
	var srcs []*zip.File
	for _, tp := range tps {
		f := r.fileForType(tp)
		present := false
		for _, src := range srcs {
			if src == f {
				present = true
				break
			}
		}
		if !present {
			srcs = append(srcs, f)
		}
	}
	return &Entries{
		files:      srcs,
		typeFilter: tps,
	}
}

func (r *Reader) fileForType(tp FieldType) *zip.File {
	switch tp {
	case CheungBauerIndex,
		CihaiT,
		Cowles,
		DaeJaweon,
		FennIndex,
		GSR,
		HanYu,
		IRGDaeJaweon,
		IRGDaiKanwaZiten,
		IRGHanyuDaZidian,
		IRGKangXi,
		KangXi,
		Karlgren,
		Lau,
		Matthews,
		MeyerWempe,
		Morohashi,
		Nelson,
		SBGY:
		return r.dictionaryIndices
	case AlternateTotalStrokes,
		Cangjie,
		CheungBauer,
		Fenn,
		FourCornerCode,
		Frequency,
		GradeLevel,
		HDZRadBreak,
		HKGlyph,
		Phonetic,
		Strange,
		UnihanCore2020:
		return r.dictionaryLikeData
	case CompatibilityVariant,
		IICore,
		IRG_GSource,
		IRG_HSource,
		IRG_JSource,
		IRG_KPSource,
		IRG_KSource,
		IRG_MSource,
		IRG_SSource,
		IRG_TSource,
		IRG_UKSource,
		IRG_USource,
		IRG_VSource,
		RSUnicode,
		TotalStrokes:
		return r.irgSources
	case AccountingNumeric,
		OtherNumeric,
		PrimaryNumeric:
		return r.numericValues
	case BigFive,
		CCCII,
		CNS1986,
		CNS1992,
		EACC,
		GB0,
		GB1,
		GB3,
		GB5,
		GB7,
		GB8,
		HKSCS,
		IBMJapan,
		Ja,
		JinmeiyoKanji,
		Jis0,
		Jis1,
		JIS0213,
		JoyoKanji,
		KoreanEducationHanja,
		KoreanName,
		KPS0,
		KPS1,
		KSC0,
		KSC1,
		MainlandTelegraph,
		PseudoGB1,
		TaiwanTelegraph,
		TGH,
		Xerox:
		return r.otherMappins
	case RSAdobe_Japan1_6,
		RSKangXi:
		return r.radicalStrokeCounts
	case Cantonese,
		Definition,
		Hangul,
		HanyuPinlu,
		HanyuPinyin,
		JapaneseKun,
		JapaneseOn,
		Korean,
		Mandarin,
		Tang,
		TGHZ2013,
		Vietnamese,
		XHC1983:
		return r.readings
	case SemanticVariant,
		SimplifiedVariant,
		SpecializedSemanticVariant,
		SpoofingVariant,
		TraditionalVariant,
		ZVariant:
		return r.variants
	default:
		panic(tp)
	}
}

type Entries struct {
	Filter     func(rune) bool
	typeFilter []FieldType
	files      []*zip.File
	cur        io.Closer
	reader     *bufio.Scanner
}

func (e *Entries) Close() {
	if e.cur == nil {
		return
	}
	e.cur.Close()
}

func (e *Entries) Next() (rune, Field, error) {
	for {
		if e.cur == nil {
			if len(e.files) == 0 {
				return 0, nil, io.EOF
			}
			rd, err := e.files[0].Open()
			if err != nil {
				return 0, nil, err
			}
			e.cur = rd
			e.reader = bufio.NewScanner(rd)
			e.files = e.files[1:]
		}
		if !e.reader.Scan() {
			e.cur.Close()
			if err := e.reader.Err(); err != nil {
				return 0, nil, err
			}
			e.cur = nil
			continue
		}
		ln := e.reader.Text()
		if len(ln) == 0 || strings.HasPrefix(ln, "#") {
			continue
		}
		charEnc, rest, ok := strings.Cut(ln, "\t")
		if !ok {
			return 0, nil, fmt.Errorf("invalid line (no tab): %q", ln)
		}
		charNum, ok := strings.CutPrefix(charEnc, "U+")
		if !ok {
			return 0, nil, fmt.Errorf("invalid char: %q", charEnc)
		}
		charDec, err := strconv.ParseInt(charNum, 16, 64)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid char: %q", charEnc)
		}
		r := rune(charDec)
		if e.Filter != nil && !e.Filter(r) {
			continue
		}
		fieldNameEnc, rest, ok := strings.Cut(rest, "\t")
		if !ok {
			return 0, nil, fmt.Errorf("invalid line (no second tab): %q, %q", r, ln)
		}
		fieldName, ok := strings.CutPrefix(fieldNameEnc, "k")
		if !ok {
			return 0, nil, fmt.Errorf("invalid field name: %q", fieldNameEnc)
		}
		tp := getFieldType(fieldName)
		if tp == 0 {
			return 0, nil, fmt.Errorf("unknown field name: %q", fieldName)
		}
		if e.typeFilter != nil {
			found := false
			for _, f := range e.typeFilter {
				if f == tp {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		field := getField(tp)
		err = field.Parse(rest)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to parse field %q for %q: %w",
				fieldName, r, err)
		}
		return r, field, nil
	}
}
