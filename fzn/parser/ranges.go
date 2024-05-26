package parser

import (
	"fmt"

	"github.com/rhartert/gofzn/fzn/tok"
)

// Parsers for ranges
// ------------------
//
// Grammar:
//
//	RangeInt   ::= <int-lit> ".." <int-lit>
//
//	RangeFloat ::= <float-lit> ".." <float-lit>

// parseIntRange parses a range of integers.
func parseFloatRange(p *parser) (r RangeFloat, err error) {
	r.Min, err = parseFloatLit(p)
	if err != nil {
		return RangeFloat{}, err
	}
	if !p.nextIf(tok.Range) {
		return RangeFloat{}, fmt.Errorf("missing range \"..\" separator")
	}
	r.Max, err = parseFloatLit(p)
	if err != nil {
		return RangeFloat{}, err
	}
	return r, nil
}

// parseIntRange parses a range of integers.
func parseIntRange(p *parser) (r RangeInt, err error) {
	r.Min, err = parseIntLit(p)
	if err != nil {
		return RangeInt{}, err
	}
	if !p.nextIf(tok.Range) {
		return RangeInt{}, fmt.Errorf("missing range \"..\" separator")
	}
	r.Max, err = parseIntLit(p)
	if err != nil {
		return RangeInt{}, err
	}
	return r, nil
}
