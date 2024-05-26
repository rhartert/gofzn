package parser

import (
	"fmt"

	"github.com/rhartert/gofzn/fzn/tok"
)

func isVariable(p *parser) bool {
	switch p.lookAhead(0).Type {
	case tok.Var:
		return true
	case tok.Array:
		return p.lookAhead(7).Type == tok.Var
	default:
		return false
	}
}

func parseVariable(p *parser) (v *Variable, err error) {
	v = &Variable{}

	if p.lookAhead(0).Type == tok.Array {
		v.Array, err = parseArrayOf(p)
		if err != nil {
			return nil, fmt.Errorf("error parsing variable array: %w", err)
		}
	}

	v.Domain, v.Type, err = parseBasicVarType(p)
	if err != nil {
		return nil, err
	}

	if !p.nextIf(tok.Colon) {
		return nil, fmt.Errorf("missing colon")
	}

	v.Identifier, err = parseIdentifier(p)
	if err != nil {
		return nil, fmt.Errorf("error parsing variable identifier: %w", err)
	}

	v.Annotations, err = parseAnnotations(p)
	if err != nil {
		return nil, fmt.Errorf("error parsing variable annotations: %w", err)
	}

	if p.nextIf(tok.Assign) {
		v.Exprs, err = parseArrayLit(p)
		if err != nil {
			return nil, fmt.Errorf("error parsing variable expressions: %w", err)
		}
	}

	if !p.nextIf(tok.EOI) {
		return nil, fmt.Errorf("missing ';'")
	}

	return v, nil
}

// Grammar:
//
//	<basic-var-type> ::= "var" <basic-par-type>
//	                   | "var" <float-literal> ".." <float-literal>
//	                   | "var" <int-literal> ".." <int-literal>
//	                   | "var" "{" <int-literal> "," ... "}"
//	                   | "var" "set" "of" <int-literal> ".." <int-literal>
//	                   | "var" "set" "of" "{" [ <int-literal> "," ... ] "}"
func parseBasicVarType(p *parser) (VarDomain, VarType, error) {
	if !p.nextIf(tok.Var) {
		return VarDomain{}, VarTypeUnknown, fmt.Errorf("shoud start with var")
	}

	switch t := p.lookAhead(0); t.Type {
	case tok.BoolType:
		p.next()
		return VarDomain{}, VarTypeBool, nil
	case tok.IntType:
		p.next()
		return VarDomain{}, VarTypeIntRange, nil
	case tok.FloatType:
		p.next()
		return VarDomain{}, VarTypeFloatRange, nil
	case tok.FloatLit:
		r, err := parseFloatRange(p)
		if err != nil {
			return VarDomain{}, VarTypeUnknown, err
		}
		return toFloatDomain(r), VarTypeFloatRange, nil
	case tok.IntLit:
		r, err := parseIntRange(p)
		if err != nil {
			return VarDomain{}, VarTypeUnknown, err
		}
		return toIntDomain(r), VarTypeIntRange, nil
	case tok.SetStart:
		is, err := parseSetIntLit(p)
		if err != nil {
			return VarDomain{}, VarTypeUnknown, err
		}
		return VarDomain{IntDomain: &is}, VarTypeIntSet, nil
	case tok.Set:
		is, err := parseSetOfInt(p)
		if err != nil {
			return VarDomain{}, VarTypeUnknown, err
		}
		return VarDomain{IntDomain: &is}, VarTypeIntSet, nil
	default:
		return VarDomain{}, VarTypeUnknown, fmt.Errorf("invalid variable")
	}
}

func toFloatDomain(r RangeFloat) VarDomain {
	return VarDomain{FloatDomain: &SetFloatLit{Values: [][]float64{{r.Min, r.Max}}}}
}

func toIntDomain(r RangeInt) VarDomain {
	return VarDomain{IntDomain: &SetIntLit{Values: [][]int{{r.Min, r.Max}}}}
}
