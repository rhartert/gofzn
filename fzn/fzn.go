// Package fzn contains functionality for parsing FlatZinc models.
package fzn

import (
	"bufio"
	"fmt"
	"io"

	"github.com/rhartert/gofzn/fzn/tok"
)

// Model represents a FlatZinc model.
type Model struct {
	Predicates        []Predicate
	ParamDeclarations []ParamDeclaration
	VarDeclarations   []VarDeclaration
	Constraints       []Constraint
	SolveGoals        []SolveGoal
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
		return nil, err
	}
	return &mb.Model, nil
}

// Handler is an interface that clients must implement to handle the parsed
// components of a FlatZinc model.
type Handler interface {
	AddPredicate(p *Predicate) error
	AddParamDeclaration(p *ParamDeclaration) error
	AddVarDeclaration(v *VarDeclaration) error
	AddConstraint(c *Constraint) error
	AddSolveGoal(s *SolveGoal) error
}

// Parse parses an FZN model from the given reader, calling the handler
// functions along the way. The parser stops and returns an error if the
// handler returns an error.
//
// The function does not verify that the model itself is valid. For example,
// it does not verify that variable domains are consistent or that referenced
// entities have been declared. It is the Handler's responsibility to perform
// these verifications.
func Parse(reader io.Reader, handler Handler) error {
	tokenizer := tok.Tokenizer{}
	scanner := bufio.NewScanner(reader)

	i := 0 // line number
	for scanner.Scan() {
		i++

		line := scanner.Text()
		if line == "" { // TODO: this should ideally be done in the parser itself
			continue
		}

		tokens, err := tokenizer.Tokenize(line)
		if err != nil {
			return fmt.Errorf("tokenizer error at line %d: %w", i, err)
		}
		if err := parseInstruction(tokens, handler); err != nil {
			return fmt.Errorf("parser error at line %d: %w", i, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading FlatZinc model: %w", err)
	}

	return nil
}

// modelBuilder wraps Model to implement the Handler interface.
type modelBuilder struct {
	Model Model
}

func (mb *modelBuilder) AddPredicate(p *Predicate) error {
	mb.Model.Predicates = append(mb.Model.Predicates, *p)
	return nil
}

func (mb *modelBuilder) AddParamDeclaration(p *ParamDeclaration) error {
	mb.Model.ParamDeclarations = append(mb.Model.ParamDeclarations, *p)
	return nil
}

func (mb *modelBuilder) AddVarDeclaration(v *VarDeclaration) error {
	mb.Model.VarDeclarations = append(mb.Model.VarDeclarations, *v)
	return nil
}

func (mb *modelBuilder) AddConstraint(c *Constraint) error {
	mb.Model.Constraints = append(mb.Model.Constraints, *c)
	return nil
}

func (mb *modelBuilder) AddSolveGoal(s *SolveGoal) error {
	mb.Model.SolveGoals = append(mb.Model.SolveGoals, *s)
	return nil
}
