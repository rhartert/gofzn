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

	le, err := parseLiteral(p)
	if err != nil {
		return BasicExpr{}, fmt.Errorf("invalid basic expression: %w", err)
	}
	return BasicExpr{Literal: le}, nil
}
