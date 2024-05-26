package fzn

import (
	"fmt"

	"github.com/rhartert/gofzn/fzn/tok"
)

func isComment(p *parser) bool {
	return p.lookAhead(0).Type == tok.Comment
}

func parseComment(p *parser) (string, error) {
	t := p.next()
	if t.Type != tok.Comment {
		return "", fmt.Errorf("comment should start with '%%'")
	}
	return t.Value, nil
}
