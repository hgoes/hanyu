//go:generate go run ../cmd/gen-simp-trad
package simplified

import "github.com/hgoes/hanyu/dict"

func To(from string) string {
	runes := []rune(from)
	if ToInplace(runes) {
		return string(runes)
	}
	return from
}

func ToInplace(from []rune) bool {
	replaced := false
	for len(from) > 0 {
		l, m := dict.Main.Lookup(from)
		if l == 0 {
			repl, ok := Replacements[from[0]]
			if ok {
				replaced = true
				from[0] = repl
			}
			from = from[1:]
			continue
		}
		if m[0].Simplified != "" {
			for i, repl := range []rune(m[0].Simplified) {
				from[i] = repl
			}
			replaced = true
		}
		from = from[l:]
	}
	return replaced
}
