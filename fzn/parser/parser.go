// Package parser provides functionalities to parse sequences of FlatZinc
// lexical tokens into higher level objects such as variables and constraints.
package parser

import (
	"fmt"

	tok "github.com/rhartert/gofzn/fzn/tokenizer"
)

// Handler is an interface that clients must implement to handle the parsed
// components of a FlatZinc model.
type Handler interface {
	AddPredicate(p *Predicate) error
	AddParameter(p *Parameter) error
	AddVariable(v *Variable) error
	AddConstraint(c *Constraint) error
	AddSolveGoal(s *SolveGoal) error
}

type Predicate struct{}
type Parameter struct{}
type Variable struct{}
type Constraint struct{}
type SolveGoal struct{}

// ParseInstruction parses a sequence of tokens representing a FlatZinc
// instruction and uses the provided Handler to manage the parsed elements.
// It returns an error if the parsing fails or if the Handler reports an error.
func ParseInstruction(tokens []tok.Token, handler Handler) error {
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

// parse analyzes the tokens and delegates handling of parsed elements to the
// appropriate Handler method. It returns an error if parsing fails or if the
// Handler reports an error.
func (p *parser) parse() error {
	return fmt.Errorf("not implemented")
}
