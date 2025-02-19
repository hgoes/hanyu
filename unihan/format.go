package unihan

type FieldType byte

type Field interface {
	Type() FieldType
	Parse(c string) error
}

func getFieldType(s string) FieldType {
	switch s {
	case "AccountingNumeric":
		return AccountingNumeric
	case "AlternateTotalStrokes":
		return AlternateTotalStrokes
	case "BigFive":
		return BigFive
	case "Cangjie":
		return Cangjie
	case "Cantonese":
		return Cantonese
	case "CCCII":
		return CCCII
	case "CheungBauer":
		return CheungBauer
	case "CheungBauerIndex":
		return CheungBauerIndex
	case "CihaiT":
		return CihaiT
	case "CNS1986":
		return CNS1986
	case "CNS1992":
		return CNS1992
	case "CompatibilityVariant":
		return CompatibilityVariant
	case "Cowles":
		return Cowles
	case "DaeJaweon":
		return DaeJaweon
	case "Definition":
		return Definition
	case "EACC":
		return EACC
	case "Fenn":
		return Fenn
	case "FennIndex":
		return FennIndex
	case "FourCornerCode":
		return FourCornerCode
	case "Frequency":
		return Frequency
	case "GB0":
		return GB0
	case "GB1":
		return GB1
	case "GB3":
		return GB3
	case "GB5":
		return GB5
	case "GB7":
		return GB7
	case "GB8":
		return GB8
	case "GradeLevel":
		return GradeLevel
	case "GSR":
		return GSR
	case "Hangul":
		return Hangul
	case "HanYu":
		return HanYu
	case "HanyuPinlu":
		return HanyuPinlu
	case "HanyuPinyin":
		return HanyuPinyin
	case "HDZRadBreak":
		return HDZRadBreak
	case "HKGlyph":
		return HKGlyph
	case "HKSCS":
		return HKSCS
	case "IBMJapan":
		return IBMJapan
	case "IICore":
		return IICore
	case "IRG_GSource":
		return IRG_GSource
	case "IRG_HSource":
		return IRG_HSource
	case "IRG_JSource":
		return IRG_JSource
	case "IRG_KPSource":
		return IRG_KPSource
	case "IRG_KSource":
		return IRG_KSource
	case "IRG_MSource":
		return IRG_MSource
	case "IRG_SSource":
		return IRG_SSource
	case "IRG_TSource":
		return IRG_TSource
	case "IRG_UKSource":
		return IRG_UKSource
	case "IRG_USource":
		return IRG_USource
	case "IRG_VSource":
		return IRG_VSource
	case "IRGDaeJaweon":
		return IRGDaeJaweon
	case "IRGDaiKanwaZiten":
		return IRGDaiKanwaZiten
	case "IRGHanyuDaZidian":
		return IRGHanyuDaZidian
	case "IRGKangXi":
		return IRGKangXi
	case "Ja":
		return Ja
	case "JapaneseKun":
		return JapaneseKun
	case "JapaneseOn":
		return JapaneseOn
	case "JinmeiyoKanji":
		return JinmeiyoKanji
	case "Jis0":
		return Jis0
	case "Jis1":
		return Jis1
	case "JIS0213":
		return JIS0213
	case "JoyoKanji":
		return JoyoKanji
	case "KangXi":
		return KangXi
	case "Karlgren":
		return Karlgren
	case "Korean":
		return Korean
	case "KoreanEducationHanja":
		return KoreanEducationHanja
	case "KoreanName":
		return KoreanName
	case "KPS0":
		return KPS0
	case "KPS1":
		return KPS1
	case "KSC0":
		return KSC0
	case "KSC1":
		return KSC1
	case "Lau":
		return Lau
	case "MainlandTelegraph":
		return MainlandTelegraph
	case "Mandarin":
		return Mandarin
	case "Matthews":
		return Matthews
	case "MeyerWempe":
		return MeyerWempe
	case "Morohashi":
		return Morohashi
	case "Nelson":
		return Nelson
	case "OtherNumeric":
		return OtherNumeric
	case "Phonetic":
		return Phonetic
	case "PrimaryNumeric":
		return PrimaryNumeric
	case "PseudoGB1":
		return PseudoGB1
	case "RSAdobe_Japan1_6":
		return RSAdobe_Japan1_6
	case "RSKangXi":
		return RSKangXi
	case "RSUnicode":
		return RSUnicode
	case "SBGY":
		return SBGY
	case "SemanticVariant":
		return SemanticVariant
	case "SimplifiedVariant":
		return SimplifiedVariant
	case "SpecializedSemanticVariant":
		return SpecializedSemanticVariant
	case "SpoofingVariant":
		return SpoofingVariant
	case "Strange":
		return Strange
	case "TaiwanTelegraph":
		return TaiwanTelegraph
	case "Tang":
		return Tang
	case "TGH":
		return TGH
	case "TGHZ2013":
		return TGHZ2013
	case "TotalStrokes":
		return TotalStrokes
	case "TraditionalVariant":
		return TraditionalVariant
	case "UnihanCore2020":
		return UnihanCore2020
	case "Vietnamese":
		return Vietnamese
	case "Xerox":
		return Xerox
	case "XHC1983":
		return XHC1983
	case "ZVariant":
		return ZVariant
	default:
		return 0
	}
}

const (
	AccountingNumeric FieldType = iota + 1
	AlternateTotalStrokes
	BigFive
	Cangjie
	Cantonese
	CCCII
	CheungBauer
	CheungBauerIndex
	CihaiT
	CNS1986
	CNS1992
	CompatibilityVariant
	Cowles
	DaeJaweon
	Definition
	EACC
	Fenn
	FennIndex
	FourCornerCode
	Frequency
	GB0
	GB1
	GB3
	GB5
	GB7
	GB8
	GradeLevel
	GSR
	Hangul
	HanYu
	HanyuPinlu
	HanyuPinyin
	HDZRadBreak
	HKGlyph
	HKSCS
	IBMJapan
	IICore
	IRG_GSource
	IRG_HSource
	IRG_JSource
	IRG_KPSource
	IRG_KSource
	IRG_MSource
	IRG_SSource
	IRG_TSource
	IRG_UKSource
	IRG_USource
	IRG_VSource
	IRGDaeJaweon
	IRGDaiKanwaZiten
	IRGHanyuDaZidian
	IRGKangXi
	Ja
	JapaneseKun
	JapaneseOn
	JinmeiyoKanji
	Jis0
	Jis1
	JIS0213
	JoyoKanji
	KangXi
	Karlgren
	Korean
	KoreanEducationHanja
	KoreanName
	KPS0
	KPS1
	KSC0
	KSC1
	Lau
	MainlandTelegraph
	Mandarin
	Matthews
	MeyerWempe
	Morohashi
	Nelson
	OtherNumeric
	Phonetic
	PrimaryNumeric
	PseudoGB1
	RSAdobe_Japan1_6
	RSKangXi
	RSUnicode
	SBGY
	SemanticVariant
	SimplifiedVariant
	SpecializedSemanticVariant
	SpoofingVariant
	Strange
	TaiwanTelegraph
	Tang
	TGH
	TGHZ2013
	TotalStrokes
	TraditionalVariant
	UnihanCore2020
	Vietnamese
	Xerox
	XHC1983
	ZVariant
)
