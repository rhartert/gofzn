package fzn

import (
	"fmt"

	"github.com/rhartert/gofzn/fzn/tok"
)

// parseInstruction parses a sequence of tokens representing a FlatZinc
// instruction and uses the provided Handler to manage the parsed elements.
// It returns an error if the parsing fails or if the Handler reports an error.
func parseInstruction(tokens []tok.Token, handler Handler) error {
	p := parser{
		handler: handler,
		tokens:  tokens,
		pos:     0,
	}
	return p.parse()
}

type parser struct {
	handler Handler
	tokens  []tok.Token
	pos     int
}

// next returns the next token or a tEOF token if there's no token left.
func (p *parser) next() tok.Token {
	if p.pos >= len(p.tokens) {
		return tok.Token{Type: tok.EOF}
	}
	p.pos++
	return p.tokens[p.pos-1]
}

// nextIf returns true and consumes the next token iff it is of type tt.
func (p *parser) nextIf(tt tok.Type) bool {
	if p.pos >= len(p.tokens) {
		return tt == tok.EOF
	}
	if p.tokens[p.pos].Type == tt {
		p.pos++
		return true
	}
	return false
}

// lookAhead returns the token at n positions from the current position without
// impacting the result of p.next.
func (p *parser) lookAhead(n int) tok.Token {
	if i := p.pos + n; i < len(p.tokens) {
		return p.tokens[i]
	}
	return tok.Token{Type: tok.EOF}
}

// parse analyzes the tokens and delegates handling of parsed elements to the
// appropriate Handler method. It returns an error if parsing fails or if the
// Handler reports an error.
func (p *parser) parse() error {
	switch {
	case isComment(p):
		_, err := parseComment(p) // drop comments
		return err
	case isPredicate(p):
		pred, err := parsePredicate(p)
		if err != nil {
			return err
		}
		return p.handler.AddPredicate(pred)
	case isParameter(p):
		param, err := parseParameter(p)
		if err != nil {
			return err
		}
		return p.handler.AddParameter(param)
	case isVariable(p):
		v, err := parseVariable(p)
		if err != nil {
			return err
		}
		return p.handler.AddVariable(v)
	case isConstraint(p):
		c, err := parseConstraint(p)
		if err != nil {
			return err
		}
		return p.handler.AddConstraint(c)
	case isSolveGoal(p):
		s, err := parseSolveGoal(p)
		if err != nil {
			return err
		}
		return p.handler.AddSolveGoal(s)
	default:
		return fmt.Errorf("unrecognized instruction")
	}
}
