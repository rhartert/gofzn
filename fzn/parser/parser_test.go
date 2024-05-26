package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	tok "github.com/rhartert/gofzn/fzn/tokenizer"
)

// instruction implements the Handler interface.
type instruction struct {
	Predicate  *Predicate
	Parameter  *Parameter
	Variable   *Variable
	Constraint *Constraint
	SolveGoal  *SolveGoal
}

func (i *instruction) AddPredicate(p *Predicate) error {
	i.Predicate = p
	return nil
}

func (i *instruction) AddParameter(p *Parameter) error {
	i.Parameter = p
	return nil
}

func (i *instruction) AddVariable(v *Variable) error {
	i.Variable = v
	return nil
}

func (i *instruction) AddConstraint(c *Constraint) error {
	i.Constraint = c
	return nil
}

func (i *instruction) AddSolveGoal(sg *SolveGoal) error {
	i.SolveGoal = sg
	return nil
}

func TestParseInstruction(t *testing.T) {
	testCases := []struct {
		tokens  []tok.Token
		want    instruction
		wantErr bool
	}{
		{
			tokens:  nil,
			wantErr: true,
		},
		{
			tokens:  []tok.Token{},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		got := instruction{}
		gotErr := ParseInstruction(tc.tokens, &got)

		if tc.wantErr && gotErr == nil {
			t.Errorf("ParseInstruction(): want error, got nil")
		}
		if !tc.wantErr && gotErr != nil {
			t.Errorf("ParseInstruction(): want no error, got %s", gotErr)
		}
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("ParseInstruction(): mismatch (-want +got):\n%s", diff)
		}
	}
}
