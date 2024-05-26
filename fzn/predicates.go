package fzn

import (
	"fmt"
	"strings"

	"github.com/rhartert/gofzn/fzn/tok"
)

func isPredicate(p *parser) bool {
	return p.lookAhead(0).Type == tok.Predicate
}

// parsePredicate returns a string of all the token value contained in the
// predicate. Said otherwise, the parser effectively treats predicates as
// comments
func parsePredicate(p *parser) (*Predicate, error) {
	if p.next().Type != tok.Predicate {
		return nil, fmt.Errorf("not a predicate")
	}
	sb := strings.Builder{}
	sb.WriteString("predicate")
	for t := p.next(); t.Type != tok.EOF; t = p.next() {
		sb.WriteByte(' ')
		sb.WriteString(t.Value)
	}
	return &Predicate{Value: sb.String()}, nil
}
