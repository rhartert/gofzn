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
// impacting the result of p.next. In particular, lookAhead(0) peeks at the
// next token without consuming it.
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
	for p.lookAhead(0).Type != tok.EOF {
		switch {
		case isComment(p):
			_, err := parseComment(p) // drop comments
			if err != nil {
				return err
			}
		case isPredicate(p):
			pred, err := parsePredicate(p)
			if err != nil {
				return err
			}
			if err := p.handler.AddPredicate(pred); err != nil {
				return err
			}
		case isParamDeclaration(p):
			param, err := parseParamDeclaration(p)
			if err != nil {
				return err
			}
			if err := p.handler.AddParamDeclaration(param); err != nil {
				return err
			}
		case isVarDeclaration(p):
			v, err := parseVarDeclaration(p)
			if err != nil {
				return err
			}
			if err := p.handler.AddVarDeclaration(v); err != nil {
				return err
			}
		case isConstraint(p):
			c, err := parseConstraint(p)
			if err != nil {
				return err
			}
			if err := p.handler.AddConstraint(c); err != nil {
				return err
			}
		case isSolveGoal(p):
			s, err := parseSolveGoal(p)
			if err != nil {
				return err
			}
			if err := p.handler.AddSolveGoal(s); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unrecognized instruction")
		}
	}

	return nil
}
