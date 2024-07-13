package fzn

import (
	"fmt"

	"github.com/rhartert/gofzn/fzn/tok"
)

func isPredicate(p *parser) bool {
	return p.lookAhead(0).Type == tok.Predicate
}

// parsePredicate returns a string of all the token value contained in the
// predicate. Said otherwise, the parser effectively treats predicates as
// comments
func parsePredicate(p *parser) (pred *Predicate, err error) {
	if p.next().Type != tok.Predicate {
		return nil, fmt.Errorf("not a predicate")
	}

	pred = &Predicate{}
	pred.Identifier, err = parseIdentifier(p)
	if err != nil {
		return nil, fmt.Errorf("error parsing predicate identifier: %w", err)
	}

	if !p.nextIf(tok.TupleStart) {
		return nil, fmt.Errorf("missing ( after predicate identifier")
	}

	for !p.nextIf(tok.TupleEnd) {
		pp, err := parsePredicateParam(p)
		if err != nil {
			return nil, fmt.Errorf("error parsing predicate parameter: %w", err)
		}
		pred.Parameters = append(pred.Parameters, pp)

		if !p.nextIf(tok.Comma) && p.lookAhead(0).Type != tok.TupleEnd {
			return nil, fmt.Errorf("missing comma")
		}
	}

	if !p.nextIf(tok.EOI) {
		return nil, fmt.Errorf("missing end of predicate instruction")
	}

	return pred, nil
}

func parsePredicateParam(p *parser) (PredParam, error) {
	pp := PredParam{}

	if p.lookAhead(0).Type == tok.Array {
		a, err := parseArrayOf(p, false)
		if err != nil {
			return PredParam{}, err
		}
		pp.Array = a
	}

	switch p.lookAhead(0).Type {
	case tok.Var:
		v, err := parseVariable(p)
		if err != nil {
			return PredParam{}, err
		}
		pp.VarType = v.Type
	default: // parameter
		pt, err := parseParType(p)
		if err != nil {
			return PredParam{}, err
		}
		pp.ParType = pt
	}

	if !p.nextIf(tok.Colon) {
		return PredParam{}, fmt.Errorf("missing ':'")
	}

	id, err := parseIdentifier(p)
	if err != nil {
		return PredParam{}, err
	}
	pp.Identifier = id

	return pp, nil
}
