package parser

import (
	"fmt"

	"github.com/rhartert/gofzn/fzn/tok"
)

// parseAnnotations parses a sequence of annotations. It returns nil if the
// sequence is empty.
func parseAnnotations(p *parser) ([]Annotation, error) {
	var annotations []Annotation
	for p.nextIf(tok.AnnStart) {
		a, err := parseAnnotation(p)
		if err != nil {
			return nil, err
		}
		annotations = append(annotations, a)
	}
	return annotations, nil
}

func parseAnnotation(p *parser) (Annotation, error) {
	id, err := parseIdentifier(p)
	if err != nil {
		return Annotation{}, err
	}

	a := Annotation{Identifier: id}

	if !p.nextIf(tok.TupleStart) {
		return a, nil
	}

	a.Exprs = make([][]AnnExpr, 0, 8)
	for !p.nextIf(tok.TupleEnd) {
		ae, err := parseAnnExpr(p)
		if err != nil {
			return Annotation{}, err
		}
		a.Exprs = append(a.Exprs, ae)

		if !p.nextIf(tok.Comma) && p.lookAhead(0).Type != tok.TupleEnd {
			return Annotation{}, fmt.Errorf("missing comma")
		}
	}

	return a, nil
}

func parseAnnExpr(p *parser) ([]AnnExpr, error) {
	if !p.nextIf(tok.ArrayStart) {
		ae, err := basicAnnExpr(p)
		if err != nil {
			return nil, err
		}
		return []AnnExpr{*ae}, nil
	}

	aes := []AnnExpr{}
	for !p.nextIf(tok.ArrayEnd) {
		ae, err := basicAnnExpr(p)
		if err != nil {
			return nil, err
		}
		aes = append(aes, *ae)

		if !p.nextIf(tok.Comma) && p.lookAhead(0).Type != tok.ArrayEnd {
			return nil, fmt.Errorf("missing comma")
		}
	}

	return aes, nil
}

func basicAnnExpr(p *parser) (*AnnExpr, error) {
	switch {
	case isBasicLiteralExpr(p):
		ble, err := parseBasicLiteralExpr(p)
		if err != nil {
			return nil, err
		}
		return &AnnExpr{BasicLitExpr: &ble}, nil
	case isIdentifier(p):
		if p.lookAhead(1).Type == tok.TupleStart {
			a, err := parseAnnotation(p)
			if err != nil {
				return nil, err
			}
			return &AnnExpr{Annotation: &a}, nil
		}
		id, err := parseIdentifier(p)
		if err != nil {
			return nil, err
		}
		return &AnnExpr{VarID: &id}, nil
	case isStringLit(p):
		sl, err := parseStringLit(p)
		if err == nil {
			return nil, err
		}
		return &AnnExpr{StringLit: &sl}, nil
	default:
		return nil, fmt.Errorf("unknown basicAnnExpr: %s", p.lookAhead(0))
	}
}
