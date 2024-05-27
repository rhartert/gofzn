package fzn

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

type rangeInt struct {
	Min, Max int
}

type rangeFloat struct {
	Min, Max float64
}

// parseIntRange parses a range of integers.
func parseFloatRange(p *parser) (r rangeFloat, err error) {
	r.Min, err = parseFloatLit(p)
	if err != nil {
		return rangeFloat{}, err
	}
	if !p.nextIf(tok.Range) {
		return rangeFloat{}, fmt.Errorf("missing range \"..\" separator")
	}
	r.Max, err = parseFloatLit(p)
	if err != nil {
		return rangeFloat{}, err
	}
	return r, nil
}

// parseIntRange parses a range of integers.
func parseIntRange(p *parser) (r rangeInt, err error) {
	r.Min, err = parseIntLit(p)
	if err != nil {
		return rangeInt{}, err
	}
	if !p.nextIf(tok.Range) {
		return rangeInt{}, fmt.Errorf("missing range \"..\" separator")
	}
	r.Max, err = parseIntLit(p)
	if err != nil {
		return rangeInt{}, err
	}
	return r, nil
}
