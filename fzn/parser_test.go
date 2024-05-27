package fzn

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rhartert/gofzn/fzn/tok"
	"github.com/rhartert/ptr"
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

type testCaseParser struct {
	desc    string
	tokens  []tok.Token
	want    instruction
	wantErr bool
}

func TestParseInstruction_invalid(t *testing.T) {
	testParseInstruction(t, []testCaseParser{
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
func TestParseInstruction_comment(t *testing.T) {
	testParseInstruction(t, []testCaseParser{
		{
			desc: "drop comment",
			tokens: []tok.Token{
				{Type: tok.Comment, Value: "test comment"},
			},
			want: instruction{},
		},
	})
}

func TestParseInstruction_parameter(t *testing.T) {
	testParseInstruction(t, []testCaseParser{
		{
			desc: "missing type",
			tokens: []tok.Token{
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.Assign},
				{Type: tok.IntLit, Value: "42"},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "missing colon",
			tokens: []tok.Token{
				{Type: tok.IntType},
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.Assign},
				{Type: tok.IntLit, Value: "42"},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "missing identifier",
			tokens: []tok.Token{
				{Type: tok.IntType},
				{Type: tok.Colon},
				{Type: tok.Assign},
				{Type: tok.IntLit, Value: "42"},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "missing assign",
			tokens: []tok.Token{
				{Type: tok.IntType},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.IntLit, Value: "42"},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "missing assigned expression",
			tokens: []tok.Token{
				{Type: tok.IntType},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.Assign},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "missing end of instruction",
			tokens: []tok.Token{
				{Type: tok.IntType},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.Assign},
				{Type: tok.IntLit, Value: "42"},
			},
			wantErr: true,
		},
		{
			desc: "invalid type",
			tokens: []tok.Token{
				{Type: tok.Error},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.Assign},
				{Type: tok.IntLit, Value: "42"},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "invalid array",
			tokens: []tok.Token{
				{Type: tok.Array},
				{Type: tok.ArrayStart},
				{Type: tok.ArrayEnd},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.Assign},
				{Type: tok.IntLit, Value: "42"},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "valid int parameter",
			tokens: []tok.Token{
				{Type: tok.IntType},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.Assign},
				{Type: tok.IntLit, Value: "42"},
				{Type: tok.EOI},
			},
			want: instruction{
				Parameter: &Parameter{
					Identifier: "foo",
					Type:       ParTypeInt,
					Exprs:      []Literal{{Int: ptr.Of(42)}},
				},
			},
		},
		{
			desc: "valid bool parameter",
			tokens: []tok.Token{
				{Type: tok.BoolType},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.Assign},
				{Type: tok.BoolLit, Value: "true"},
				{Type: tok.EOI},
			},
			want: instruction{
				Parameter: &Parameter{
					Identifier: "foo",
					Type:       ParTypeBool,
					Exprs:      []Literal{{Bool: ptr.Of(true)}},
				},
			},
		},
		{
			desc: "valid float parameter",
			tokens: []tok.Token{
				{Type: tok.FloatType},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.Assign},
				{Type: tok.FloatLit, Value: "42.0"},
				{Type: tok.EOI},
			},
			want: instruction{
				Parameter: &Parameter{
					Identifier: "foo",
					Type:       ParTypeFloat,
					Exprs:      []Literal{{Float: ptr.Of(42.0)}},
				},
			},
		},
		{
			desc: "valid set of int parameter",
			tokens: []tok.Token{
				{Type: tok.Set},
				{Type: tok.Of},
				{Type: tok.IntType},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.Assign},
				{Type: tok.SetStart},
				{Type: tok.IntLit, Value: "42"},
				{Type: tok.Comma},
				{Type: tok.IntLit, Value: "44"},
				{Type: tok.Comma},
				{Type: tok.IntLit, Value: "45"},
				{Type: tok.SetEnd},
				{Type: tok.EOI},
			},
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
			desc: "valid array of parameters",
			tokens: []tok.Token{
				{Type: tok.Array},
				{Type: tok.ArrayStart},
				{Type: tok.IntLit, Value: "1"},
				{Type: tok.Range},
				{Type: tok.IntLit, Value: "2"},
				{Type: tok.ArrayEnd},
				{Type: tok.Of},
				{Type: tok.IntType},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.Assign},
				{Type: tok.ArrayStart},
				{Type: tok.IntLit, Value: "42"},
				{Type: tok.Comma},
				{Type: tok.IntLit, Value: "1337"},
				{Type: tok.ArrayEnd},
				{Type: tok.EOI},
			},
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

func TestParseInstruction_variable(t *testing.T) {
	testParseInstruction(t, []testCaseParser{
		{
			desc: "missign var keyword",
			tokens: []tok.Token{
				{Type: tok.IntType},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "X"},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "missing type",
			tokens: []tok.Token{
				{Type: tok.Var},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "X"},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "missign colon",
			tokens: []tok.Token{
				{Type: tok.Var},
				{Type: tok.IntType},
				{Type: tok.Identifier, Value: "X"},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "missing identifier",
			tokens: []tok.Token{
				{Type: tok.Var},
				{Type: tok.IntType},
				{Type: tok.Colon},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "missing end of instruction",
			tokens: []tok.Token{
				{Type: tok.Var},
				{Type: tok.IntType},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "X"},
			},
			wantErr: true,
		},
		{
			desc: "valid int variable (no domain)",
			tokens: []tok.Token{
				{Type: tok.Var},
				{Type: tok.IntType},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "X"},
				{Type: tok.EOI},
			},
			want: instruction{
				Variable: &Variable{
					Identifier: "X",
					Type:       VarTypeIntRange,
					Domain:     VarDomain{},
				},
			},
		},
		{
			desc: "valid int variable (with annotation)",
			tokens: []tok.Token{
				{Type: tok.Var},
				{Type: tok.IntType},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "X"},
				{Type: tok.AnnStart},
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.EOI},
			},
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
			desc: "valid int variable (range)",
			tokens: []tok.Token{
				{Type: tok.Var},
				{Type: tok.IntLit, Value: "1"},
				{Type: tok.Range},
				{Type: tok.IntLit, Value: "5"},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "X"},
				{Type: tok.EOI},
			},
			want: instruction{
				Variable: &Variable{
					Identifier: "X",
					Type:       VarTypeIntRange,
					Domain:     VarDomain{IntDomain: &SetIntLit{Values: [][]int{{1, 5}}}},
				},
			},
		},
		{
			desc: "valid float variable (range)",
			tokens: []tok.Token{
				{Type: tok.Var},
				{Type: tok.FloatLit, Value: "0.1"},
				{Type: tok.Range},
				{Type: tok.FloatLit, Value: "0.5"},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "X"},
				{Type: tok.EOI},
			},
			want: instruction{
				Variable: &Variable{
					Identifier: "X",
					Type:       VarTypeFloatRange,
					Domain:     VarDomain{FloatDomain: &SetFloatLit{Values: [][]float64{{0.1, 0.5}}}},
				},
			},
		},
		{
			desc: "valid int var (set)",
			tokens: []tok.Token{
				{Type: tok.Var},
				{Type: tok.SetStart},
				{Type: tok.IntLit, Value: "1"},
				{Type: tok.Comma},
				{Type: tok.IntLit, Value: "3"},
				{Type: tok.SetEnd},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "X"},
				{Type: tok.EOI},
			},
			want: instruction{
				Variable: &Variable{
					Identifier: "X",
					Type:       VarTypeIntSet,
					Domain:     VarDomain{IntDomain: &SetIntLit{Values: [][]int{{1, 1}, {3, 3}}}},
				},
			},
		},
		{
			desc: "valid set of int var (range)",
			tokens: []tok.Token{
				{Type: tok.Var},
				{Type: tok.Set},
				{Type: tok.Of},
				{Type: tok.IntLit, Value: "1"},
				{Type: tok.Range},
				{Type: tok.IntLit, Value: "3"},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "X"},
				{Type: tok.EOI},
			},
			want: instruction{
				Variable: &Variable{
					Identifier: "X",
					Type:       VarTypeIntSet,
					Domain:     VarDomain{IntDomain: &SetIntLit{Values: [][]int{{1, 3}}}},
				},
			},
		},
		{
			desc: "valid set of int var (set)",
			tokens: []tok.Token{
				{Type: tok.Var},
				{Type: tok.Set},
				{Type: tok.Of},
				{Type: tok.SetStart},
				{Type: tok.IntLit, Value: "1"},
				{Type: tok.Comma},
				{Type: tok.IntLit, Value: "3"},
				{Type: tok.SetEnd},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "X"},
				{Type: tok.EOI},
			},
			want: instruction{
				Variable: &Variable{
					Identifier: "X",
					Type:       VarTypeIntSet,
					Domain:     VarDomain{IntDomain: &SetIntLit{Values: [][]int{{1, 1}, {3, 3}}}},
				},
			},
		},
		{
			desc: "valid array of int var",
			tokens: []tok.Token{
				{Type: tok.Array},
				{Type: tok.ArrayStart},
				{Type: tok.IntLit, Value: "1"},
				{Type: tok.Range},
				{Type: tok.IntLit, Value: "2"},
				{Type: tok.ArrayEnd},
				{Type: tok.Of},
				{Type: tok.Var},
				{Type: tok.IntType},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "X"},
				{Type: tok.EOI},
			},
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
			desc: "valid array of int var (with assign)",
			tokens: []tok.Token{
				{Type: tok.Array},
				{Type: tok.ArrayStart},
				{Type: tok.IntLit, Value: "1"},
				{Type: tok.Range},
				{Type: tok.IntLit, Value: "2"},
				{Type: tok.ArrayEnd},
				{Type: tok.Of},
				{Type: tok.Var},
				{Type: tok.IntType},
				{Type: tok.Colon},
				{Type: tok.Identifier, Value: "X"},
				{Type: tok.Assign},
				{Type: tok.ArrayStart},
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.Comma},
				{Type: tok.Identifier, Value: "bar"},
				{Type: tok.ArrayEnd},
				{Type: tok.EOI},
			},
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

func TestParseInstruction_constraint(t *testing.T) {
	testParseInstruction(t, []testCaseParser{
		{
			desc: "missing constraint keyword",
			tokens: []tok.Token{
				{Type: tok.Identifier, Value: "foobar"},
				{Type: tok.TupleStart},
				{Type: tok.Identifier, Value: "X_VAR"},
				{Type: tok.TupleEnd},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "missing constraint identifier",
			tokens: []tok.Token{
				{Type: tok.Constraint},
				{Type: tok.TupleStart},
				{Type: tok.Identifier, Value: "X_VAR"},
				{Type: tok.TupleEnd},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "missing parameters",
			tokens: []tok.Token{
				{Type: tok.Constraint},
				{Type: tok.Identifier, Value: "foobar"},
				{Type: tok.EOI},
			},
			wantErr: true,
		},
		{
			desc: "missing end of instruction",
			tokens: []tok.Token{
				{Type: tok.Constraint},
				{Type: tok.Identifier, Value: "foobar"},
				{Type: tok.TupleStart},
				{Type: tok.TupleEnd},
			},
			wantErr: true,
		},
		{
			desc: "valid constraint (no parameter)",
			tokens: []tok.Token{
				{Type: tok.Constraint},
				{Type: tok.Identifier, Value: "foobar"},
				{Type: tok.TupleStart},
				{Type: tok.TupleEnd},
				{Type: tok.EOI},
			},
			want: instruction{
				Constraint: &Constraint{
					Identifier:  "foobar",
					Expressions: []Expr{},
				},
			},
		},
		{
			desc: "valid constraint (one parameter)",
			tokens: []tok.Token{
				{Type: tok.Constraint},
				{Type: tok.Identifier, Value: "foobar"},
				{Type: tok.TupleStart},
				{Type: tok.Identifier, Value: "X_VAR"},
				{Type: tok.TupleEnd},
				{Type: tok.EOI},
			},
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
			desc: "valid constraint (two parameter with annotation)",
			tokens: []tok.Token{
				{Type: tok.Constraint},
				{Type: tok.Identifier, Value: "foobar"},
				{Type: tok.TupleStart},
				{Type: tok.Identifier, Value: "X_VAR"},
				{Type: tok.Comma},
				{Type: tok.Identifier, Value: "Y_VAR"},
				{Type: tok.TupleEnd},
				{Type: tok.AnnStart},
				{Type: tok.Identifier, Value: "annotation"},
				{Type: tok.EOI},
			},
			want: instruction{
				Constraint: &Constraint{
					Identifier: "foobar",
					Expressions: []Expr{
						{Exprs: []BasicExpr{{Identifier: "X_VAR"}}},
						{Exprs: []BasicExpr{{Identifier: "Y_VAR"}}},
					},
					Annotations: []Annotation{{Identifier: "annotation"}},
				},
			},
		},
		{
			desc: "valid constraint (array of parameters)",
			tokens: []tok.Token{
				{Type: tok.Constraint},
				{Type: tok.Identifier, Value: "foobar"},
				{Type: tok.TupleStart},
				{Type: tok.ArrayStart},
				{Type: tok.Identifier, Value: "X_VAR"},
				{Type: tok.Comma},
				{Type: tok.Identifier, Value: "Y_VAR"},
				{Type: tok.ArrayEnd},
				{Type: tok.TupleEnd},
				{Type: tok.AnnStart},
				{Type: tok.Identifier, Value: "annotation"},
				{Type: tok.EOI},
			},
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
					Annotations: []Annotation{{Identifier: "annotation"}},
				},
			},
		},
	})
}

func TestParseInstruction_solveGoal(t *testing.T) {
	testParseInstruction(t, []testCaseParser{
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

func testParseInstruction(t *testing.T, testCases []testCaseParser) {
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := instruction{}
			gotErr := parseInstruction(tc.tokens, &got)

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
