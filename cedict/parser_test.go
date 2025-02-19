package cedict

import (
	"bytes"
	"compress/gzip"
	"testing"
)

func Test(t *testing.T) {
	input := "#! someKey=someVal\n" +
		"做事 做事 [zuo4 shi4] /to work/to handle matters/to have a job/\n"
	var buf bytes.Buffer
	wr := gzip.NewWriter(&buf)
	wr.Write([]byte(input))
	wr.Close()
	p, err := New(&buf)
	if err != nil {
		t.Fatal(err)
	}
	ln1, err := p.Next()
	if err != nil {
		t.Fatal(err)
	}
	switch sub := ln1.(type) {
	case Metadata:
		if sub.Key != "someKey" {
			t.Errorf("wrong metadata key: %q", sub.Key)
		}
		if sub.Value != "someVal" {
			t.Errorf("wrong metadata value: %q", sub.Value)
		}
	default:
		t.Fatalf("wrong line type: %T", sub)
	}
	ln2, err := p.Next()
	if err != nil {
		t.Fatal(err)
	}
	switch sub := ln2.(type) {
	case Entry:
		if len(sub.Meaning) != 3 {
			t.Errorf("wrong number of meanings: %d", len(sub.Meaning))
		}
		if sub.Meaning[0] != "to work" {
			t.Errorf("wrong meaning 0: %q", sub.Meaning[0])
		}
		if sub.Meaning[1] != "to handle matters" {
			t.Errorf("wrong meaning 1: %q", sub.Meaning[1])
		}
	}
}
