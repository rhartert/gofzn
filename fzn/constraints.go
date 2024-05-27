package fzn

import (
	"fmt"

	"github.com/rhartert/gofzn/fzn/tok"
)

// Parsers for Constraints
// -----------------------
//
// Grammar:
//
//  "constraint" <identifier> "(" [ <expr> "," ... ] ")" <annotations> ";"

func isConstraint(p *parser) bool {
	return p.lookAhead(0).Type == tok.Constraint
}

func parseConstraint(p *parser) (*Constraint, error) {
	if !p.nextIf(tok.Constraint) {
		return nil, fmt.Errorf("constraints should start with tConstraint")
	}

	id, err := parseIdentifier(p)
	if err != nil {
		return nil, fmt.Errorf("error parsing constraint identifier: %w", err)
	}

	if !p.nextIf(tok.TupleStart) {
		return nil, fmt.Errorf("missing '('")
	}
	exprs := make([]Expr, 0, 8)
	for !p.nextIf(tok.TupleEnd) {
		expr, err := parseExpr(p)
		if err != nil {
			return nil, fmt.Errorf("error parsing constraint expression: %w", err)
		}
		exprs = append(exprs, expr)

		if !p.nextIf(tok.Comma) && p.lookAhead(0).Type != tok.TupleEnd {
			return nil, fmt.Errorf("missing comma")
		}
	}

	anns, err := parseAnnotations(p)
	if err != nil {
		return nil, fmt.Errorf("error parsing constraint annotations: %w", err)
	}

	if !p.nextIf(tok.EOI) {
		return nil, fmt.Errorf("missing end of constraint declaration ';'")
	}

	c := &Constraint{
		Identifier:  id,
		Expressions: exprs,
	}
	if len(anns) != 0 {
		c.Annotations = anns
	}

	return c, nil
}

func parseExpr(p *parser) (Expr, error) {
	if p.lookAhead(0).Type == tok.ArrayStart {
		es, err := parseArrayLit(p)
		if err != nil {
			return Expr{}, err
		}
		return Expr{
			Exprs: es,
		}, nil
	}

	e, err := parseBasicExpr(p)
	if err != nil {
		return Expr{}, err
	}
	return Expr{
		Expr: &e,
	}, nil
}
