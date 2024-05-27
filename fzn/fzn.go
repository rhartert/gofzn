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
// This function only checks for syntactic correctness and does not verify
// that the Model is semantically correct (see [Parse] for details).
func ParseModel(reader io.Reader) (*Model, error) {
	mb := &modelBuilder{}
	if err := Parse(reader, mb); err != nil {
		return nil, err
	}
	return &mb.Model, nil
}

// Handler is an interface that works with the parser (see [Parse]) to handle
// parsed FlatZinc model components such as predicates, parameters, variables,
// constraints, and solve goals. Implementations of this interface define how
// these parsed components should be processed.
type Handler interface {
	HandlePredicate(p *Predicate) error
	HandleParamDeclaration(p *ParamDeclaration) error
	HandleVarDeclaration(v *VarDeclaration) error
	HandleConstraint(c *Constraint) error
	HandleSolveGoal(s *SolveGoal) error
}

// Parse parses a FlatZinc model from the reader, actioning the given [Handler]
// interface to handle the parsed model items. It stops and returns an error if
// the model is syntactically incorrect or if the handler returns an error.
//
// This function only checks for syntactic correctness and does not verify
// that the model is semantically correct. For instance, the following
// instruction will be parsed successfully despite defining a variable with an
// invalid domain:
//
//	var 10..0: X; // semantically invalid domain
//
// It is the responsibility of the given Handler's implementation to validate
// the model's semantic to meet its need.
func Parse(reader io.Reader, handler Handler) error {
	tokenizer := tok.Tokenizer{}
	scanner := bufio.NewScanner(reader)

	i := 0 // line number
	for scanner.Scan() {
		i++

		line := scanner.Text()
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

func (mb *modelBuilder) HandlePredicate(p *Predicate) error {
	mb.Model.Predicates = append(mb.Model.Predicates, *p)
	return nil
}

func (mb *modelBuilder) HandleParamDeclaration(p *ParamDeclaration) error {
	mb.Model.ParamDeclarations = append(mb.Model.ParamDeclarations, *p)
	return nil
}

func (mb *modelBuilder) HandleVarDeclaration(v *VarDeclaration) error {
	mb.Model.VarDeclarations = append(mb.Model.VarDeclarations, *v)
	return nil
}

func (mb *modelBuilder) HandleConstraint(c *Constraint) error {
	mb.Model.Constraints = append(mb.Model.Constraints, *c)
	return nil
}

func (mb *modelBuilder) HandleSolveGoal(s *SolveGoal) error {
	mb.Model.SolveGoals = append(mb.Model.SolveGoals, *s)
	return nil
}
