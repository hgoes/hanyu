//go:generate go run ../cmd/gen-numbers
package numbers

type Parser struct {
	value      int64
	positional bool
	digits     bool
}

func (p *Parser) Consume(r rune) bool {
	val, ok := All[r]
	if !ok {
		return false
	}
	if val < 10 {
		if p.digits {
			p.value = p.value*10 + val
			return true
		}
		if p.positional {
			if val == 0 {
				return true
			}
			if p.value%10 != 0 {
				return false
			}
		} else if p.value%10 != 0 {
			p.digits = true
			p.value = p.value*10 + val
			return true
		}
		p.value += val
		return true
	}
	if p.digits {
		return false
	}
	p.positional = true
	before := (p.value / (val * 10)) * val * 10
	after := p.value % val
	if after == 0 {
		// assume that 一 has been omitted
		p.value = before + val
		return true
	}
	p.value = before + val*after
	return true
}

func (p *Parser) Value() int64 {
	return p.value
}
