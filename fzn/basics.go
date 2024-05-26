package fzn

import (
	"fmt"

	"github.com/rhartert/gofzn/fzn/tok"
)

func parseBasicExpr(p *parser) (BasicExpr, error) {
	if p.lookAhead(0).Type == tok.Identifier {
		t := p.next()
		return BasicExpr{Identifier: t.Value}, nil
	}

	le, err := parseBasicLiteralExpr(p)
	if err != nil {
		return BasicExpr{}, fmt.Errorf("invalid basic expression: %w", err)
	}
	return BasicExpr{LiteralExpr: le}, nil
}

func isBasicLiteralExpr(p *parser) bool {
	switch p.lookAhead(0).Type {
	case tok.BoolLit, tok.IntLit, tok.FloatLit:
		return true
	default:
		return isSetLit(p)
	}
}

func parseBasicLiteralExpr(p *parser) (BasicLitExpr, error) {
	if isSetLit(p) {
		s, err := parseSetLit(p)
		if err != nil {
			return BasicLitExpr{}, err
		}
		return BasicLitExpr{Set: &s}, nil
	}

	switch tt := p.lookAhead(0).Type; tt {
	case tok.BoolLit:
		b, err := parseBoolLit(p)
		if err != nil {
			return BasicLitExpr{}, err
		}
		return BasicLitExpr{Bool: &b}, nil
	case tok.IntLit:
		i, err := parseIntLit(p)
		if err != nil {
			return BasicLitExpr{}, err
		}
		return BasicLitExpr{Int: &i}, nil
	case tok.FloatLit:
		f, err := parseFloatLit(p)
		if err != nil {
			return BasicLitExpr{}, err
		}
		return BasicLitExpr{Float: &f}, nil
	default:
		return BasicLitExpr{}, fmt.Errorf("token is not part of valid literal: %s", tt)
	}
}
