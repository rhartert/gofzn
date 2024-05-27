package fzn

import (
	_ "embed"
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/google/go-cmp/cmp"
	"github.com/rhartert/ptr"
)

// instruction implements the Handler interface and serves as a convenient way
// to compare parsed instructions.
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
	input   string
	want    instruction
	wantErr bool
}

func TestParse_invalidReader(t *testing.T) {
	got := instruction{}
	gotErr := Parse(iotest.ErrReader(errors.New("test error")), &got)

	if gotErr == nil {
		t.Errorf("Parse(): want error, got none")
	}
	if diff := cmp.Diff(instruction{}, got); diff != "" {
		t.Errorf("Parse(): mismatch (-want +got):\n%s", diff)
	}
}

func TestParse_comment(t *testing.T) {
	testParse(t, []testCase{
		{
			input: "%%",
			want:  instruction{}, // drop comment
		},
		{
			input: "%% a comment",
			want:  instruction{}, // drop comment
		},
		{
			input: "%% a comment %% with a comment",
			want:  instruction{}, // drop comment
		},
	})
}

func TestParse_predicate(t *testing.T) {
	testParse(t, []testCase{
		{
			input: "predicate foo(int)",
			want: instruction{
				Predicate: &Predicate{Value: "predicate foo ( int )"},
			},
		},
	})
}

func TestParse_parameter(t *testing.T) {
	testParse(t, []testCase{
		{
			input:   ": foo = 42;",
			wantErr: true,
		},
		{
			input:   "int foo = 42;",
			wantErr: true,
		},
		{
			input:   "int: = 42;",
			wantErr: true,
		},
		{
			input:   "foo = 42;",
			wantErr: true,
		},
		{
			input:   "int: foo;",
			wantErr: true,
		},
		{
			input:   "42;",
			wantErr: true,
		},
		{
			input:   "int: foo = 42",
			wantErr: true,
		},
		{
			input: "int: foo = 42;",
			want: instruction{
				Parameter: &Parameter{
					Identifier: "foo",
					Type:       ParTypeInt,
					Exprs:      []Literal{{Int: ptr.Of(42)}},
				},
			},
		},
		{
			input: "bool: foo = true;",
			want: instruction{
				Parameter: &Parameter{
					Identifier: "foo",
					Type:       ParTypeBool,
					Exprs:      []Literal{{Bool: ptr.Of(true)}},
				},
			},
		},
		{
			input: "float: foo = 42.0;",
			want: instruction{
				Parameter: &Parameter{
					Identifier: "foo",
					Type:       ParTypeFloat,
					Exprs:      []Literal{{Float: ptr.Of(42.0)}},
				},
			},
		},
		{
			input: "set of int: foo = {42, 44, 45};",
			want: instruction{
				Parameter: &Parameter{
					Identifier: "foo",
					Type:       ParTypeSetOfInt,
					Exprs: []Literal{{SetInt: &SetIntLit{
						Values: [][]int{{42, 42}, {44, 45}},
					}}},
				},
			},
		},
		{
			input: "array [1..2] of int: foo = [42, 1337];",
			want: instruction{
				Parameter: &Parameter{
					Identifier: "foo",
					Type:       ParTypeInt,
					Array:      &Array{1, 2},
					Exprs: []Literal{
						{Int: ptr.Of(42)},
						{Int: ptr.Of(1337)},
					},
				},
			},
		},
	})
}

func TestParse_variables(t *testing.T) {
	testParse(t, []testCase{
		{
			input:   "int: X;",
			wantErr: true,
		},
		{
			input:   "var: X;",
			wantErr: true,
		},
		{
			input:   "var int X;",
			wantErr: true,
		},
		{
			input:   "var int: ;",
			wantErr: true,
		},
		{
			input:   "var int: X",
			wantErr: true,
		},
		{
			input: "var int: X;",
			want: instruction{
				Variable: &Variable{
					Identifier: "X",
					Type:       VarTypeIntRange,
					Domain:     VarDomain{},
				},
			},
		},
		{
			input: "var int : X ::foo;",
			want: instruction{
				Variable: &Variable{
					Identifier:  "X",
					Type:        VarTypeIntRange,
					Domain:      VarDomain{},
					Annotations: []Annotation{{Identifier: "foo"}},
				},
			},
		},
		{
			input: "var 1..5: X;",
			want: instruction{
				Variable: &Variable{
					Identifier: "X",
					Type:       VarTypeIntRange,
					Domain:     VarDomain{IntDomain: &SetIntLit{Values: [][]int{{1, 5}}}},
				},
			},
		},
		{
			input: "var 0.1..0.5: X;",
			want: instruction{
				Variable: &Variable{
					Identifier: "X",
					Type:       VarTypeFloatRange,
					Domain:     VarDomain{FloatDomain: &SetFloatLit{Values: [][]float64{{0.1, 0.5}}}},
				},
			},
		},
		{
			input: "var {1, 3}: X;",
			want: instruction{
				Variable: &Variable{
					Identifier: "X",
					Type:       VarTypeIntSet,
					Domain:     VarDomain{IntDomain: &SetIntLit{Values: [][]int{{1, 1}, {3, 3}}}},
				},
			},
		},
		{
			input: "var set of 1..3: X;",
			want: instruction{
				Variable: &Variable{
					Identifier: "X",
					Type:       VarTypeIntSet,
					Domain:     VarDomain{IntDomain: &SetIntLit{Values: [][]int{{1, 3}}}},
				},
			},
		},
		{
			input: "var set of {1, 3} : X;",
			want: instruction{
				Variable: &Variable{
					Identifier: "X",
					Type:       VarTypeIntSet,
					Domain:     VarDomain{IntDomain: &SetIntLit{Values: [][]int{{1, 1}, {3, 3}}}},
				},
			},
		},
		{
			input: "array [1..2] of var int: X;",
			want: instruction{
				Variable: &Variable{
					Identifier: "X",
					Type:       VarTypeIntRange,
					Array:      &Array{1, 2},
					Domain:     VarDomain{},
				},
			},
		},
		{
			input: "array [1..2] of var int: X = [foo, bar];",
			want: instruction{
				Variable: &Variable{
					Identifier: "X",
					Type:       VarTypeIntRange,
					Array:      &Array{1, 2},
					Domain:     VarDomain{},
					Exprs: []BasicExpr{
						{Identifier: "foo"},
						{Identifier: "bar"},
					},
				},
			},
		},
	})
}

func TestParse_constraint(t *testing.T) {
	testParse(t, []testCase{
		{
			input:   "foobar(X_VAR);",
			wantErr: true,
		},
		{
			input:   "constraint (X_VAR);",
			wantErr: true,
		},
		{
			input:   "constraint foobar;",
			wantErr: true,
		},
		{
			input:   "constraint foobar(;",
			wantErr: true,
		},
		{
			input:   "constraint foo bar;",
			wantErr: true,
		},
		{
			input: "constraint foobar ();",
			want: instruction{
				Constraint: &Constraint{
					Identifier:  "foobar",
					Expressions: []Expr{},
				},
			},
		},
		{
			input: "constraint foobar(X_VAR);",
			want: instruction{
				Constraint: &Constraint{
					Identifier: "foobar",
					Expressions: []Expr{
						{Exprs: []BasicExpr{{Identifier: "X_VAR"}}},
					},
				},
			},
		},
		{
			input: "constraint foobar(X_VAR, Y_VAR) ::bar;",
			want: instruction{
				Constraint: &Constraint{
					Identifier: "foobar",
					Expressions: []Expr{
						{Exprs: []BasicExpr{{Identifier: "X_VAR"}}},
						{Exprs: []BasicExpr{{Identifier: "Y_VAR"}}},
					},
					Annotations: []Annotation{{Identifier: "bar"}},
				},
			},
		},
		{
			input: "constraint foobar ([X_VAR, Y_VAR]) :: bar;",
			want: instruction{
				Constraint: &Constraint{
					Identifier: "foobar",
					Expressions: []Expr{{
						IsArray: true,
						Exprs: []BasicExpr{
							{Identifier: "X_VAR"},
							{Identifier: "Y_VAR"},
						},
					}},
					Annotations: []Annotation{{Identifier: "bar"}},
				},
			},
		},
	})
}

func TestParse_solveGoal(t *testing.T) {
	testParse(t, []testCase{
		{
			input:   "solve;",
			wantErr: true,
		},
		{
			input:   "satisfy;",
			wantErr: true,
		},
		{
			input:   "solve satisfy",
			wantErr: true,
		},
		{
			input:   "solve minimize;",
			wantErr: true,
		},
		{
			input:   "solve maximize;",
			wantErr: true,
		},
		{
			input:   "solve maximize minimize;",
			wantErr: true,
		},
		{
			input: "solve satisfy;",
			want: instruction{
				SolveGoal: &SolveGoal{
					SolveMethod: SolveMethodSatisfy,
				},
			},
		},
		{
			input: "solve ::foobar satisfy;",
			want: instruction{
				SolveGoal: &SolveGoal{
					SolveMethod: SolveMethodSatisfy,
					Annotations: []Annotation{{Identifier: "foobar"}},
				},
			},
		},
		{
			input: "solve minimize OBJ_VAR;",
			want: instruction{
				SolveGoal: &SolveGoal{
					SolveMethod: SolveMethodMinimize,
					Objective:   BasicExpr{Identifier: "OBJ_VAR"},
				},
			},
		},
		{
			input: "solve ::foobar maximize OBJ_VAR;",
			want: instruction{
				SolveGoal: &SolveGoal{
					SolveMethod: SolveMethodMaximize,
					Annotations: []Annotation{{Identifier: "foobar"}},
					Objective:   BasicExpr{Identifier: "OBJ_VAR"},
				},
			},
		},
	})
}

func testParse(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			got := instruction{}
			gotErr := Parse(strings.NewReader(tc.input), &got)

			if tc.wantErr && gotErr == nil {
				t.Errorf("Parse(): want error, got nil")
			}
			if !tc.wantErr && gotErr != nil {
				t.Errorf("Parse(): want no error, got %s", gotErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Parse(): mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
