// Package fzn contains functionality for parsing FlatZinc models.
package fzn

import (
	"bufio"
	"fmt"
	"io"

	"github.com/rhartert/gofzn/fzn/parser"
	"github.com/rhartert/gofzn/fzn/tokenizer"
)

// Model represents a FlatZinc model.
type Model struct {
	Predicates  []parser.Predicate
	Parameters  []parser.Parameter
	Variables   []parser.Variable
	Constraints []parser.Constraint
	SolveGoals  []parser.SolveGoal
}

// Parse parses an FZN model from the given reader, calling the handler
// functions along the way. The parser stops and returns an error if the
// handler returns an error.
//
// The function does not verify that the model itself is valid. For example,
// it does not verify that variable domains are consistent or that referenced
// entities have been declared. It is the Handler's responsibility to perform
// these verifications.
func Parse(reader io.Reader, handler parser.Handler) error {
	tok := tokenizer.Tokenizer{}
	scanner := bufio.NewScanner(reader)

	i := 0 // line number
	for scanner.Scan() {
		i++
		line := scanner.Text()
		tokens, err := tok.Tokenize(line)
		if err != nil {
			return fmt.Errorf("tokenizer error at line %d: %w", i, err)
		}
		if err := parser.ParseInstruction(tokens, handler); err != nil {
			return fmt.Errorf("parser error at line %d: %w", i, err)
		}
	}

	return nil
}

// ParseModel reads a FlatZinc model from the provided reader and returns a
// fully constructed Model.
//
// The function does not verify that the model is valid. For example, it does
// not verify that variable domains are valid or that referenced entities have
// been declared.
func ParseModel(reader io.Reader) (*Model, error) {
	mb := &modelBuilder{}
	if err := Parse(reader, mb); err != nil {
		return nil, fmt.Errorf("error reading model: %w", err)
	}
	return &mb.Model, nil
}

// modelBuilder wraps Model to implement the Handler interface.
type modelBuilder struct {
	Model Model
}

func (mb *modelBuilder) AddPredicate(p *parser.Predicate) error {
	mb.Model.Predicates = append(mb.Model.Predicates, *p)
	return nil
}

func (mb *modelBuilder) AddParameter(p *parser.Parameter) error {
	mb.Model.Parameters = append(mb.Model.Parameters, *p)
	return nil
}

func (mb *modelBuilder) AddVariable(v *parser.Variable) error {
	mb.Model.Variables = append(mb.Model.Variables, *v)
	return nil
}

func (mb *modelBuilder) AddConstraint(c *parser.Constraint) error {
	mb.Model.Constraints = append(mb.Model.Constraints, *c)
	return nil
}

func (mb *modelBuilder) AddSolveGoal(s *parser.SolveGoal) error {
	mb.Model.SolveGoals = append(mb.Model.SolveGoals, *s)
	return nil
}
