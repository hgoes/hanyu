// Provides a parser for the CEDICT (https://www.mdbg.net) dictionary format
package cedict

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"regexp"

	"github.com/hgoes/hanyu/pinyin"
)

// Parser reads a CEDICT dictionary file entry by entry
type Parser struct {
	reader *bufio.Scanner
	lineNr int
}

// New creates a new [Parser]. The source must be a gzip'ed
// dictionary, containing a CEDICT entry in each line.
func New(src io.Reader) (*Parser, error) {
	rd, err := gzip.NewReader(src)
	if err != nil {
		return nil, err
	}
	return &Parser{
		reader: bufio.NewScanner(rd),
	}, nil
}

var regEntry = regexp.MustCompile(`^([^ ]+) ([^ ]+) \[([^]]*)\] /+((:?[^/]+/?)+)/$`)

var regMD = regexp.MustCompile(`^ *([^ =]+) *= *(.*)$`)

// Next yields the next [Line] from the dictionary.
func (p *Parser) Next() (Line, error) {
	for {
		if !p.reader.Scan() {
			if err := p.reader.Err(); err != nil {
				return nil, err
			}
			return nil, nil
		}
		p.lineNr++
		ln := p.reader.Bytes()
		if len(ln) == 0 {
			continue
		}
		if ln[0] == '#' {
			if len(ln) > 1 && ln[1] == '!' {
				match := regMD.FindSubmatch(ln[2:])
				if len(match) == 0 {
					return nil, fmt.Errorf(
						"line %d: invalid metadata: %q",
						p.lineNr, ln[2:])
				}
				return Metadata{
					Key:   string(match[1]),
					Value: string(match[2]),
				}, nil
			}
			return Comment(string(ln[1:])), nil
		}
		match := regEntry.FindSubmatch(ln)
		if len(match) == 0 {
			return nil, fmt.Errorf(
				"line %d: invalid entry: %q", p.lineNr, ln)
		}
		trad := string(match[1])
		simp := string(match[2])
		rawPinyins := bytes.Fields(match[3])
		pinyins := make([]Pinyin, len(rawPinyins))
		for i, raw := range rawPinyins {
			r := bytes.Runes(raw)
			ok, pin, rest := pinyin.Parse(r)
			if !ok || len(rest) != 0 {
				pinyins[i].Literal = string(r)
			} else {
				pinyins[i].Pinyin = pin
			}
		}
		rawMeanings := bytes.Split(match[4], []byte{'/'})
		meanings := make([]string, len(rawMeanings))
		for i, m := range rawMeanings {
			meanings[i] = string(m)
		}
		return Entry{
			Traditional: trad,
			Simplified:  simp,
			Pinyin:      pinyins,
			Meaning:     meanings,
		}, nil
	}
}
