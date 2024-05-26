package parser

import (
	"fmt"

	tok "github.com/rhartert/gofzn/fzn/tokenizer"
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
