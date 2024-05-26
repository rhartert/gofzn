package parser

import (
	"fmt"

	tok "github.com/rhartert/gofzn/fzn/tokenizer"
)

// Parsers for sets
// ----------------
//
//	SetLit      ::= "set" "of" SetIntLit
//                | SetIntLit
//		          | SetFloatLit
//
//	SetIntLit   ::= "{" [ <int-literal> "," ... ] "}"
//		          | <int-literal> ".." <int-literal>
//
//	SetFloatLit ::= "{" [ <float-literal> "," ... ] "}"
//		          | <float-literal> ".." <float-literal>

func isSetLit(p *parser) bool {
	switch p.lookAhead(0).Type {
	case tok.SetStart:
		return true
	case tok.IntLit:
		return p.lookAhead(1).Type == tok.Range
	case tok.FloatLit:
		return p.lookAhead(1).Type == tok.Range
	default:
		return false
	}
}

func parseSetLit(p *parser) (SetLit, error) {
	la0 := p.lookAhead(0).Type
	la1 := p.lookAhead(1).Type

	switch {
	case la0 == tok.IntLit || la0 == tok.SetStart && la1 == tok.IntLit:
		is, err := parseSetIntLit(p)
		if err != nil {
			return SetLit{}, err
		}
		return SetLit{SetInt: &is}, nil
	case la0 == tok.FloatLit || la0 == tok.SetStart && la1 == tok.FloatLit:
		fs, err := parseSetFloatLit(p)
		if err != nil {
			return SetLit{}, err
		}
		return SetLit{SetFloat: &fs}, nil
	default:
		return SetLit{}, fmt.Errorf("not a set literal")
	}
}

func parseSetOfInt(p *parser) (SetIntLit, error) {
	if !p.nextIf(tok.Set) {
		return SetIntLit{}, fmt.Errorf("not a set")
	}
	if !p.nextIf(tok.Of) {
		return SetIntLit{}, fmt.Errorf("not a set")
	}
	return parseSetIntLit(p)
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
