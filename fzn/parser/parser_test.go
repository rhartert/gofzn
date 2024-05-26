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

type testCase struct {
	desc    string
	tokens  []tok.Token
	want    instruction
	wantErr bool
}

func TestParseInstruction_invalid(t *testing.T) {
	testParseInstruction(t, []testCase{
		{
			desc:    "nil sequence of token",
			tokens:  nil,
			wantErr: true,
		},
		{
			desc:    "empty sequence of token",
			tokens:  []tok.Token{},
			wantErr: true,
		},
		{
			desc: "unrecognized",
			tokens: []tok.Token{
				{Type: tok.Identifier, Value: "not an instruction"},
			},
			wantErr: true,
		},
	})
}

func TestParseInstruction_solveGoal(t *testing.T) {
	testParseInstruction(t, []testCase{
		{
			desc: "missing solve method",
			tokens: []tok.Token{
				{Type: tok.Solve},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "missing solve keyword",
			tokens: []tok.Token{
				{Type: tok.Satisfy},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "missing end of instruction",
			tokens: []tok.Token{
				{Type: tok.Solve},
				{Type: tok.Satisfy},
			},
			wantErr: true,
		},
		{
			desc: "missing minimize objective",
			tokens: []tok.Token{
				{Type: tok.Solve},
				{Type: tok.Minimize},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "missing maximize objective",
			tokens: []tok.Token{
				{Type: tok.Solve},
				{Type: tok.Maximize},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "invalid objective",
			tokens: []tok.Token{
				{Type: tok.Solve},
				{Type: tok.Maximize},
				{Type: tok.Error},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "valid solve satisfy (no annotation)",
			tokens: []tok.Token{
				{Type: tok.Solve},
				{Type: tok.Satisfy},
				{Type: tok.EOI},
			},
			want: instruction{
				SolveGoal: &SolveGoal{
					SolveMethod: SolveMethodSatisfy,
				},
			},
		},
		{
			desc: "valid solve satisfy (with annotation)",
			tokens: []tok.Token{
				{Type: tok.Solve},
				{Type: tok.AnnStart},
				{Type: tok.Identifier, Value: "foobar"},
				{Type: tok.Satisfy},
				{Type: tok.EOI},
			},
			want: instruction{
				SolveGoal: &SolveGoal{
					SolveMethod: SolveMethodSatisfy,
					Annotations: []Annotation{{Identifier: "foobar"}},
				},
			},
		},

		{
			desc: "valid solve minimize (no annotation)",
			tokens: []tok.Token{
				{Type: tok.Solve},
				{Type: tok.Minimize},
				{Type: tok.Identifier, Value: "OBJ_VAR_"},
				{Type: tok.EOI},
			},
			want: instruction{
				SolveGoal: &SolveGoal{
					SolveMethod: SolveMethodMinimize,
					Objective:   BasicExpr{Identifier: "OBJ_VAR_"},
				},
			},
		},
		{
			desc: "valid solve maximize (with annotation)",
			tokens: []tok.Token{
				{Type: tok.Solve},
				{Type: tok.AnnStart},
				{Type: tok.Identifier, Value: "foobar"},
				{Type: tok.Maximize},
				{Type: tok.Identifier, Value: "OBJ_VAR_"},
				{Type: tok.EOI},
			},
			want: instruction{
				SolveGoal: &SolveGoal{
					SolveMethod: SolveMethodMaximize,
					Annotations: []Annotation{{Identifier: "foobar"}},
					Objective:   BasicExpr{Identifier: "OBJ_VAR_"},
				},
			},
		},
	})
}

func testParseInstruction(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
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
		})
	}
}
