package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/hgoes/hanyu/cedict"
	"github.com/hgoes/hanyu/pinyin"
	"github.com/hgoes/hanyu/unihan"
)

func main() {
	rd, err := os.Open("../cedict_1_0_ts_utf-8_mdbg.txt.gz")
	if err != nil {
		panic(err)
	}
	p, err := cedict.New(rd)
	if err != nil {
		panic(err)
	}
	charDB, err := unihan.Open("../Unihan.zip")
	if err != nil {
		panic(err)
	}
	// pre-build prefered readings
	prefered := make(map[rune]pinyin.Pinyin)
	entries := charDB.Get(unihan.Mandarin)
	for {
		c, f, err := entries.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		prefered[c] = f.(*unihan.MandarinF).ReadingCN
	}
	entries.Close()
	isPrefered := func(word []rune, pin []cedict.Pinyin) bool {
		if len(word) != 1 || len(pin) != 1 || pin[0].Literal != "" {
			return false
		}
		c := word[0]
		pref, ok := prefered[c]
		if !ok {
			return false
		}
		return pin[0].Pinyin == pref
	}
	g := node{
		Next: make(map[rune]*node),
	}
	hsk, err := getHSK()
	if err != nil {
		panic(err)
	}
	var m []meaning
	i := 0
NEXT_LINE:
	for {
		ln, err := p.Next()
		if err != nil {
			panic(err)
		}
		switch sub := ln.(type) {
		case nil:
			break NEXT_LINE
		case cedict.Entry:
			trad := []rune(sub.Traditional)
			simp := []rune(sub.Simplified)
			variants := getVariant(trad, simp)
			pref := isPrefered(trad, sub.Pinyin)
			g.Insert(trad, i, pref)
			if len(variants) != 0 {
				pref := isPrefered(simp, sub.Pinyin)
				g.Insert(simp, i, pref)
			}
			m = append(m, meaning{
				Pinyin:   sub.Pinyin,
				Meaning:  sub.Meaning,
				HSKLevel: hsk[sub.Simplified],
				Variant:  variants,
			})
			i++
		}
	}
	// create a rune encoding
	var allRunes []rune
	g.collectRunes(func(r rune) {
		idx := sort.Search(len(allRunes), func(i int) bool {
			return allRunes[i] >= r
		})
		if idx >= len(allRunes) {
			allRunes = append(allRunes, r)
		} else if allRunes[idx] != r {
			allRunes = append(allRunes, 0)
			copy(allRunes[idx+1:], allRunes[idx:len(allRunes)-1])
			allRunes[idx] = r
		}
	})
	wr, err := os.Create("gen.bin")
	if err != nil {
		panic(err)
	}
	err = createBinaryDict(wr, &g, m, allRunes)
	if err != nil {
		panic(err)
	}
	if err := wr.Close(); err != nil {
		panic(err)
	}
}

func getHSK() (map[string]byte, error) {
	result := make(map[string]byte)
	for i := 1; i <= 6; i++ {
		h, err := os.Open(fmt.Sprintf("../hsk%d.txt", i))
		if err != nil {
			return nil, err
		}
		sc := bufio.NewScanner(h)
		for sc.Scan() {
			result[sc.Text()] = byte(i)
		}
		if err := sc.Err(); err != nil {
			return nil, err
		}
	}
	return result, nil
}

type node struct {
	Next    map[rune]*node
	Meaning []nodeMeaning
}

type nodeMeaning struct {
	Meaning  int
	Prefered bool
}

type meaning struct {
	Pinyin   []cedict.Pinyin
	Meaning  []string
	HSKLevel byte
	Variant  []variant
}

func (n *node) Insert(text []rune, meaning int, prefered bool) {
	if len(text) == 0 {
		n.Meaning = append(n.Meaning, nodeMeaning{
			Meaning:  meaning,
			Prefered: prefered,
		})
		return
	}
	next, ok := n.Next[text[0]]
	if !ok {
		next = &node{
			Next: make(map[rune]*node),
		}
		n.Next[text[0]] = next
	}
	next.Insert(text[1:], meaning, prefered)
}

func createBinaryDict(
	wr io.WriterAt,
	root *node,
	meanings []meaning,
	runes []rune,
) error {
	// write the length of the rune index
	if len(runes) > 0xFFFFFF {
		panic("too many runes")
	}
	_, err := putUint24(wr, 0, uint32(len(runes)))
	if err != nil {
		return err
	}
	for i, r := range runes {
		_, err = putUint24(wr, 3+int64(i)*3, uint32(r))
		if err != nil {
			return err
		}
	}
	idxOffset := 3 + int64(len(runes))*3
	// leave 3 bytes for the meanings offset
	idxSize, err := root.binary(wr, idxOffset+3, runes)
	if err != nil {
		return err
	}
	// write the meanings offset
	if idxSize > 0xFFFFFF {
		panic("meanings offset too big")
	}
	_, err = putUint24(wr, idxOffset, uint32(idxSize))
	if err != nil {
		return err
	}
	// write the meanings
	meaningOffset := idxOffset + idxSize + 3
	offsetIndex := 3 * int64(len(meanings))
	for i, m := range meanings {
		if offsetIndex > 0xFFFFFF {
			panic(fmt.Sprintf("offset index too big: %d", offsetIndex))
		}
		_, err = putUint24(wr, meaningOffset+int64(i)*3, uint32(offsetIndex))
		if err != nil {
			return err
		}
		hskSz, err := wr.WriteAt([]byte{m.HSKLevel}, meaningOffset+offsetIndex)
		if err != nil {
			return err
		}
		offsetIndex += int64(hskSz)
		pinsz, err := binaryPinyins(wr, meaningOffset+offsetIndex, m.Pinyin)
		if err != nil {
			return err
		}
		offsetIndex += int64(pinsz)
		msz, err := binaryMeanings(wr, meaningOffset+offsetIndex, m.Meaning)
		if err != nil {
			return err
		}
		offsetIndex += int64(msz)
		vsz, err := binaryVariants(wr, meaningOffset+offsetIndex, m.Variant, runes)
		if err != nil {
			return err
		}
		offsetIndex += int64(vsz)
	}
	return nil
}

func binaryPinyin(
	wr io.WriterAt,
	offset int64,
	p cedict.Pinyin,
) (int, error) {
	var buf [2]byte
	if p.Literal == "" {
		binary.BigEndian.PutUint16(buf[:], uint16(p.Pinyin)|0x8000)
		return wr.WriteAt(buf[:], offset)
	}
	litBytes := []byte(p.Literal)
	if len(litBytes) > 127 {
		panic(fmt.Sprintf("pinyin literal %q not representable", p.Literal))
	}
	buf[0] = byte(len(litBytes))
	c, err := wr.WriteAt(buf[:1], offset)
	if err != nil {
		return 0, err
	}
	sz := c
	c, err = wr.WriteAt(litBytes, offset+1)
	if err != nil {
		return 0, err
	}
	sz += c
	return sz, nil
}

func binaryPinyins(
	wr io.WriterAt,
	offset int64,
	ps []cedict.Pinyin,
) (int, error) {
	// write the size
	c, err := wr.WriteAt([]byte{byte(len(ps))}, offset)
	if err != nil {
		return 0, err
	}
	sz := c
	for _, p := range ps {
		c, err = binaryPinyin(wr, offset+int64(sz), p)
		if err != nil {
			return 0, err
		}
		sz += c
	}
	return sz, nil
}

func binaryMeanings(
	wr io.WriterAt,
	offset int64,
	meanings []string,
) (int, error) {
	if len(meanings) > 255 {
		panic("too many meanings")
	}
	c, err := wr.WriteAt([]byte{byte(len(meanings))}, offset)
	if err != nil {
		return 0, err
	}
	sz := c
	offset += int64(c)
	var buf [2]byte
	for _, meaning := range meanings {
		b := []byte(meaning)
		if len(b) > 0xFFFF {
			panic("meaning too large")
		}
		binary.BigEndian.PutUint16(buf[:], uint16(len(b)))
		c, err = wr.WriteAt(buf[:], offset)
		if err != nil {
			return 0, err
		}
		sz += c
		offset += int64(c)
		c, err = wr.WriteAt(b, offset)
		if err != nil {
			return 0, err
		}
		sz += c
		offset += int64(c)
	}
	return sz, nil
}

func encodeRune(r rune, allRunes []rune) uint16 {
	idx := sort.Search(len(allRunes), func(i int) bool {
		return allRunes[i] >= r
	})
	if idx >= len(allRunes) || allRunes[idx] != r {
		panic("rune index not found")
	}
	return uint16(idx)
}

func binaryVariants(
	wr io.WriterAt,
	offset int64,
	variants []variant,
	allRunes []rune,
) (int, error) {
	if len(variants) > 0xFF {
		panic("too many variants")
	}
	sz, err := wr.WriteAt([]byte{byte(len(variants))}, offset)
	if err != nil {
		return 0, err
	}
	offset += int64(sz)
	for _, v := range variants {
		var buf [5]byte
		buf[0] = v.Pos
		binary.BigEndian.PutUint16(buf[1:], encodeRune(v.Traditional, allRunes))
		binary.BigEndian.PutUint16(buf[3:], encodeRune(v.Simplified, allRunes))
		c, err := wr.WriteAt(buf[:], offset)
		if err != nil {
			return 0, err
		}
		sz += c
		offset += int64(c)
	}
	return sz, nil
}

func putUint24(wr io.WriterAt, at int64, val uint32) (int, error) {
	if val > 0x00FFFFFF {
		panic("out of bounds for uint24")
	}
	var enc [3]byte
	enc[0], enc[1], enc[2] = byte(val>>16), byte(val>>8), byte(val)
	return wr.WriteAt(enc[:], at)
}

func putUint16(wr io.WriterAt, at int64, val uint16) (int, error) {
	var enc [2]byte
	enc[0], enc[1] = byte(val>>8), byte(val)
	return wr.WriteAt(enc[:], at)
}

func (n *node) binary(
	wr io.WriterAt,
	off int64,
	runes []rune,
) (sz int64, err error) {
	type elem struct {
		key     rune
		runeIdx int
		nd      *node
	}
	elems := make([]elem, 0, len(n.Next))
	for key, nd := range n.Next {
		elems = append(elems, elem{
			key: key,
			nd:  nd,
		})
	}
	sort.Slice(elems, func(i, j int) bool {
		return elems[i].key < elems[j].key
	})
	offset := 0
	for i := range elems {
		idx := sort.Search(len(runes[offset:]), func(j int) bool {
			return runes[offset+j] >= elems[i].key
		})
		if idx >= len(runes)-offset || runes[offset+idx] != elems[i].key {
			panic("rune index not found")
		}
		elems[i].runeIdx = offset + idx
		offset = idx + 1
	}
	var buf [4]byte
	// write the length of the index
	if len(n.Next) > 0xFFFF {
		panic("node has too many successors")
	}
	binary.BigEndian.PutUint16(buf[:2], uint16(len(n.Next)))
	c, err := wr.WriteAt(buf[:2], off)
	if err != nil {
		return 0, err
	}
	sz += int64(c)
	off += int64(c)
	// calculate the size of the index
	offsetIndex := off + int64(len(n.Next))*5
	for i := range elems {
		// write the rune
		if elems[i].runeIdx > 0xFFFF {
			panic("rune too big")
		}
		c, err = putUint16(wr, off+int64(i)*5, uint16(elems[i].runeIdx))
		if err != nil {
			return 0, err
		}
		sz += int64(c)
		// write the offset of the node's content
		if offsetIndex > 0xFFFFFF {
			panic(fmt.Sprintf("offset index too large: %d", offsetIndex))
		}
		c, err = putUint24(wr, off+int64(i)*5+2, uint32(offsetIndex))
		if err != nil {
			return 0, err
		}
		sz += int64(c)
		// write the meanings
		c, err = wr.WriteAt([]byte{byte(len(elems[i].nd.Meaning))}, offsetIndex)
		if err != nil {
			return 0, err
		}
		sz += int64(c)
		offsetIndex += int64(c)
		// sort the meanings so that the preferred one's is on top
		sort.Slice(elems[i].nd.Meaning, func(x, y int) bool {
			if elems[i].nd.Meaning[x].Prefered {
				return true
			}
			return false
		})
		for _, m := range elems[i].nd.Meaning {
			if m.Meaning > 0xFFFFFF {
				panic("meaning too large")
			}
			c, err = putUint24(wr, offsetIndex, uint32(m.Meaning))
			if err != nil {
				return 0, err
			}
			sz += int64(c)
			offsetIndex += int64(c)
		}
		// recursively write the node's content
		ndSz, err := elems[i].nd.binary(wr, offsetIndex, runes)
		if err != nil {
			return 0, err
		}
		sz += ndSz
		// update the offset
		offsetIndex += ndSz
	}
	return sz, nil
}

func (n *node) collectRunes(f func(rune)) {
	for r, nxt := range n.Next {
		f(r)
		nxt.collectRunes(f)
	}
}

type variant struct {
	Pos         byte
	Traditional rune
	Simplified  rune
}

func getVariant(traditional, simplified []rune) []variant {
	if len(traditional) != len(simplified) {
		panic("not of equal length")
	}
	var result []variant
	for i, c := range traditional {
		if c == simplified[i] {
			continue
		}
		result = append(result, variant{
			Pos:         byte(i),
			Traditional: c,
			Simplified:  simplified[i],
		})
	}
	return result
}
