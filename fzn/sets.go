package fzn

import (
	"fmt"

	"github.com/rhartert/gofzn/fzn/tok"
)

// Parsers for sets
// ----------------
//
//	SetIntLit   ::= "{" [ <int-literal> "," ... ] "}"
//		          | <int-literal> ".." <int-literal>
//
//	SetFloatLit ::= "{" [ <float-literal> "," ... ] "}"
//		          | <float-literal> ".." <float-literal>

func parseSetOfInt(p *parser) (SetIntLit, error) {
	if !p.nextIf(tok.Set) {
		return SetIntLit{}, fmt.Errorf("not a set")
	}
	if !p.nextIf(tok.Of) {
		return SetIntLit{}, fmt.Errorf("not a set")
	}
	return parseSetIntLit(p)
}

func isSetIntLit(p *parser) bool {
	switch p.lookAhead(0).Type {
	case tok.IntLit:
		return p.lookAhead(1).Type == tok.Range
	case tok.SetStart:
		return p.lookAhead(1).Type == tok.IntLit
	default:
		return false
	}
}

// parseSetIntLit parses a set of int either represented as a range or
// a list of values.
func parseSetIntLit(p *parser) (SetIntLit, error) {
	if p.lookAhead(0).Type == tok.IntLit {
		r, err := parseIntRange(p)
		if err != nil {
			return SetIntLit{}, err
		}
		return SetIntLit{Values: [][]int{{r.Min, r.Max}}}, nil
	}

	if p.next().Type != tok.SetStart {
		return SetIntLit{}, fmt.Errorf("not a set")
	}

	values := make([]int, 0, 8)
	for !p.nextIf(tok.SetEnd) {
		i, err := parseIntLit(p)
		if err != nil {
			return SetIntLit{}, err
		}
		values = append(values, i)

		if !p.nextIf(tok.Comma) && p.lookAhead(0).Type != tok.SetEnd {
			return SetIntLit{}, fmt.Errorf("missing comma")
		}
	}

	return SetIntLit{Values: toSetRanges(values)}, nil
}

func isSetFloatLit(p *parser) bool {
	switch p.lookAhead(0).Type {
	case tok.FloatLit:
		return p.lookAhead(1).Type == tok.Range
	case tok.SetStart:
		return p.lookAhead(1).Type == tok.FloatLit
	default:
		return false
	}
}

// parseSetFloatLit parses a set of float64 either represented as a range or
// a list of values.
func parseSetFloatLit(p *parser) (SetFloatLit, error) {
	if p.lookAhead(0).Type == tok.IntLit {
		r, err := parseFloatRange(p)
		if err != nil {
			return SetFloatLit{}, err
		}
		return SetFloatLit{Values: [][]float64{{r.Min, r.Max}}}, nil
	}

	if p.next().Type != tok.SetStart {
		return SetFloatLit{}, fmt.Errorf("not a set")
	}

	values := make([]float64, 0, 8)
	for !p.nextIf(tok.SetEnd) {
		f, err := parseFloatLit(p)
		if err != nil {
			return SetFloatLit{}, err
		}
		values = append(values, f)

		if !p.nextIf(tok.Comma) && p.lookAhead(0).Type != tok.SetEnd {
			return SetFloatLit{}, fmt.Errorf("missing comma")
		}
	}

	return SetFloatLit{Values: toSetRanges(values)}, nil
}

func toSetRanges[T float64 | int](values []T) [][]T {
	if len(values) == 0 {
		return [][]T{}
	}

	ranges := make([][]T, 0, 8)
	start := values[0]
	last := start

	for _, v := range values[1:] {
		if v != last+1 {
			ranges = append(ranges, []T{start, last})
			start = v
		}
		last = v
	}
	ranges = append(ranges, []T{start, last})
	return ranges
}
