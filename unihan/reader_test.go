package unihan

import (
	"io"
	"testing"
)

func TestReader(t *testing.T) {
	rd, err := Open("../Unihan.zip")
	if err != nil {
		t.Fatal(err)
	}
	entries := rd.Get(Mandarin)
	entries.Filter = func(r rune) bool {
		return r == '听'
	}
	_, field, err := entries.Next()
	if err != nil {
		t.Fatal(err)
	}
	if sub, ok := field.(*MandarinF); ok {
		if sub.ReadingCN.String() != "tīng" {
			t.Errorf("wrong reading returned: %q", sub.ReadingCN.String())
		}
	} else {
		t.Fatalf("wrong type returned: %T", field)
	}
	_, _, err = entries.Next()
	if err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}
}
