package fzn

import (
	"fmt"

	"github.com/rhartert/gofzn/fzn/tok"
)

func parseArrayOf(p *parser, requireIndexSet bool) (*Array, error) {
	if p.next().Type != tok.Array {
		return nil, fmt.Errorf("should start with array")
	}
	if p.next().Type != tok.ArrayStart {
		return nil, fmt.Errorf("should be '['")
	}

	a := Array{}
	if requireIndexSet || !p.nextIf(tok.IntType) {
		r, err := parseIntRange(p)
		if err != nil {
			return nil, err
		}
		a.IndexSet = &IndexSet{
			Start: r.Min,
			End:   r.Max,
		}
	}

	if p.next().Type != tok.ArrayEnd {
		return nil, fmt.Errorf("should be ]")
	}
	if p.next().Type != tok.Of {
		return nil, fmt.Errorf("should be of")
	}

	return &a, nil
}

func parseArrayLit(p *parser) ([]BasicExpr, error) {
	if tt := p.next().Type; tt != tok.ArrayStart {
		return nil, fmt.Errorf("array literal should start with tArrayStart, got %s", tt)
	}

	bes := make([]BasicExpr, 0, 8)
	for !p.nextIf(tok.ArrayEnd) {
		be, err := parseBasicExpr(p)
		if err != nil {
			return nil, err
		}
		bes = append(bes, be)

		if !p.nextIf(tok.Comma) && p.lookAhead(0).Type != tok.ArrayEnd {
			return nil, fmt.Errorf("missing comma")
		}
	}

	return bes, nil
}
