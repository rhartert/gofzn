package parser

import (
	"fmt"

	"github.com/rhartert/gofzn/fzn/tok"
)

// Parsers for Parameters
// ----------------------
//
// Grammar:
//
//  <par-decl-item>  ::= <par-type> ":" <var-par-identifier> "=" <par-expr> ";"
//
//  <basic-par-type> ::= "bool" | "int" | "float" | "set of int"
//
//  <par-type>       ::= <basic-par-type>
//                     | "array" "[" <index-set> "]" "of" <basic-par-type>
//

func isParameter(p *parser) bool {
	switch p.lookAhead(0).Type {
	case tok.IntType:
		return true
	case tok.BoolType:
		return true
	case tok.FloatType:
		return true
	case tok.Set:
		return true
	case tok.Array:
		return p.lookAhead(7).Type != tok.Var
	default:
		return false
	}
}

func parseParameter(p *parser) (param *Parameter, err error) {
	param = &Parameter{}

	if p.lookAhead(0).Type == tok.Array {
		param.Array, err = parseArrayOf(p)
		if err != nil {
			return nil, fmt.Errorf("error parsing parameter array: %w", err)
		}
	}

	param.Type, err = parseParType(p)
	if err != nil {
		return nil, fmt.Errorf("error parsing parameter type: %w", err)
	}

	if !p.nextIf(tok.Colon) {
		return nil, fmt.Errorf("missing colon")
	}

	param.Identifier, err = parseIdentifier(p)
	if err != nil {
		return nil, fmt.Errorf("error parsing parameter identifier: %w", err)
	}

	if !p.nextIf(tok.Assign) {
		return nil, fmt.Errorf("missing assign")
	}

	param.Exprs, err = parseParamExpr(p)
	if err != nil {
		return nil, fmt.Errorf("error parsing parameter expressions: %w", err)
	}

	if !p.nextIf(tok.EOI) {
		return nil, fmt.Errorf("missing end of parameter declaration ';'")
	}

	return param, nil
}

func parseParType(p *parser) (ParType, error) {
	t := p.next()
	switch t.Type {
	case tok.IntType:
		return ParTypeInt, nil
	case tok.BoolType:
		return ParTypeBool, nil
	case tok.FloatType:
		return ParTypeFloat, nil
	case tok.Set:
		if !p.nextIf(tok.Of) {
			return ParTypeUnknown, fmt.Errorf("invalid set of int type")
		}
		if !p.nextIf(tok.IntType) {
			return ParTypeUnknown, fmt.Errorf("invalid set of int type")
		}
		return ParTypeSetOfInt, nil
	default:
		return ParTypeUnknown, fmt.Errorf("unknown par type: %s", t)
	}
}

func parseParamExpr(p *parser) ([]BasicLitExpr, error) {
	if !p.nextIf(tok.ArrayStart) {
		expr, err := parseBasicLiteralExpr(p)
		if err != nil {
			return nil, err
		}
		return []BasicLitExpr{expr}, nil
	}

	exprs := make([]BasicLitExpr, 0, 8)
	for !p.nextIf(tok.ArrayEnd) {
		expr, err := parseBasicLiteralExpr(p)
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)

		if !p.nextIf(tok.Comma) && p.lookAhead(0).Type != tok.ArrayEnd {
			return nil, fmt.Errorf("missing comma")
		}
	}

	return exprs, nil
}
