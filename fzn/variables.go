package fzn

import (
	"fmt"

	"github.com/rhartert/gofzn/fzn/tok"
)

func isVarDeclaration(p *parser) bool {
	switch p.lookAhead(0).Type {
	case tok.Var:
		return true
	case tok.Array:
		return p.lookAhead(7).Type == tok.Var
	default:
		return false
	}
}

func parseVarDeclaration(p *parser) (v *VarDeclaration, err error) {
	v = &VarDeclaration{}

	if p.lookAhead(0).Type == tok.Array {
		v.Array, err = parseArrayOf(p, true)
		if err != nil {
			return nil, fmt.Errorf("error parsing variable array: %w", err)
		}
	}

	v.Variable, err = parseVariable(p)
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
func parseVariable(p *parser) (Variable, error) {
	if !p.nextIf(tok.Var) {
		return Variable{}, fmt.Errorf("shoud start with var")
	}

	switch t := p.lookAhead(0); t.Type {
	case tok.BoolType:
		p.next()
		return Variable{Type: VarTypeBool}, nil
	case tok.IntType:
		p.next()
		return Variable{Type: VarTypeIntRange}, nil
	case tok.FloatType:
		p.next()
		return Variable{Type: VarTypeFloatRange}, nil
	case tok.FloatLit:
		r, err := parseFloatRange(p)
		if err != nil {
			return Variable{}, err
		}
		return toFloatDomain(r), nil
	case tok.IntLit:
		r, err := parseIntRange(p)
		if err != nil {
			return Variable{}, err
		}
		return toIntDomain(r), nil
	case tok.SetStart:
		is, err := parseSetIntLit(p)
		if err != nil {
			return Variable{}, err
		}
		return Variable{Type: VarTypeIntSet, IntDomain: &is}, nil
	case tok.Set:
		is, err := parseSetOfInt(p)
		if err != nil {
			return Variable{}, err
		}
		return Variable{Type: VarTypeIntSet, IntDomain: &is}, nil
	default:
		return Variable{}, fmt.Errorf("invalid variable")
	}
}

func toFloatDomain(r rangeFloat) Variable {
	return Variable{
		Type:        VarTypeFloatRange,
		FloatDomain: &SetFloatLit{Values: [][]float64{{r.Min, r.Max}}},
	}
}

func toIntDomain(r rangeInt) Variable {
	return Variable{
		Type:      VarTypeIntRange,
		IntDomain: &SetIntLit{Values: [][]int{{r.Min, r.Max}}},
	}
}
