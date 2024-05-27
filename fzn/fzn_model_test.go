package fzn

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rhartert/ptr"
)

var testCakesFZN = `
array [1..2] of int: X_INTRODUCED_2_ = [250,200];
array [1..2] of int: X_INTRODUCED_6_ = [75,150];
array [1..2] of int: X_INTRODUCED_8_ = [100,150];
var 0..3: b:: output_var;
var 0..6: c:: output_var;
var 0..85000: X_INTRODUCED_0_:: is_defined_var;
constraint int_lin_le(X_INTRODUCED_2_,[b,c],4000);
constraint int_lin_le(X_INTRODUCED_6_,[b,c],2000);
constraint int_lin_le(X_INTRODUCED_8_,[b,c],500);
constraint int_lin_eq([400,450,-1],[b,c,X_INTRODUCED_0_],0):: ctx_pos:: defines_var(X_INTRODUCED_0_);
solve  maximize X_INTRODUCED_0_;
`

var testCakesModel = Model{
	ParamDeclarations: []ParamDeclaration{
		{
			Identifier: "X_INTRODUCED_2_",
			Array:      &Array{Start: 1, End: 2},
			Type:       ParTypeInt,
			Literals: []Literal{
				{Int: ptr.Of(250)},
				{Int: ptr.Of(200)},
			},
		},
		{
			Identifier: "X_INTRODUCED_6_",
			Array:      &Array{Start: 1, End: 2},
			Type:       ParTypeInt,
			Literals: []Literal{
				{Int: ptr.Of(75)},
				{Int: ptr.Of(150)},
			},
		},
		{
			Identifier: "X_INTRODUCED_8_",
			Array:      &Array{Start: 1, End: 2},
			Type:       ParTypeInt,
			Literals: []Literal{
				{Int: ptr.Of(100)},
				{Int: ptr.Of(150)},
			},
		},
	},
	VarDeclarations: []VarDeclaration{
		{
			Identifier: "b",
			Variable: Variable{
				Type:      VarTypeIntRange,
				IntDomain: &SetIntLit{Values: [][]int{{0, 3}}},
			},
			Annotations: []Annotation{{Identifier: "output_var"}},
		},
		{
			Identifier: "c",
			Variable: Variable{
				Type:      VarTypeIntRange,
				IntDomain: &SetIntLit{Values: [][]int{{0, 6}}},
			},
			Annotations: []Annotation{{Identifier: "output_var"}},
		},
		{
			Identifier: "X_INTRODUCED_0_",
			Variable: Variable{
				Type:      VarTypeIntRange,
				IntDomain: &SetIntLit{Values: [][]int{{0, 85000}}},
			},
			Annotations: []Annotation{{Identifier: "is_defined_var"}},
		},
	},
	Constraints: []Constraint{
		{
			Identifier: "int_lin_le",
			Expressions: []Expr{
				{Exprs: []BasicExpr{{Identifier: "X_INTRODUCED_2_"}}},
				{IsArray: true, Exprs: []BasicExpr{{Identifier: "b"}, {Identifier: "c"}}},
				{Exprs: []BasicExpr{{Literal: Literal{Int: ptr.Of(4000)}}}},
			},
		},
		{
			Identifier: "int_lin_le",
			Expressions: []Expr{
				{Exprs: []BasicExpr{{Identifier: "X_INTRODUCED_6_"}}},
				{IsArray: true, Exprs: []BasicExpr{{Identifier: "b"}, {Identifier: "c"}}},
				{Exprs: []BasicExpr{{Literal: Literal{Int: ptr.Of(2000)}}}},
			},
		},
		{
			Identifier: "int_lin_le",
			Expressions: []Expr{
				{Exprs: []BasicExpr{{Identifier: "X_INTRODUCED_8_"}}},
				{IsArray: true, Exprs: []BasicExpr{{Identifier: "b"}, {Identifier: "c"}}},
				{Exprs: []BasicExpr{{Literal: Literal{Int: ptr.Of(500)}}}},
			},
		},
		{
			Identifier: "int_lin_eq",
			Expressions: []Expr{
				{IsArray: true, Exprs: []BasicExpr{
					{Literal: Literal{Int: ptr.Of(400)}},
					{Literal: Literal{Int: ptr.Of(450)}},
					{Literal: Literal{Int: ptr.Of(-1)}},
				}},
				{IsArray: true, Exprs: []BasicExpr{
					{Identifier: "b"},
					{Identifier: "c"},
					{Identifier: "X_INTRODUCED_0_"},
				}},
				{Exprs: []BasicExpr{{Literal: Literal{Int: ptr.Of(0)}}}},
			},
			Annotations: []Annotation{
				{Identifier: "ctx_pos"},
				{Identifier: "defines_var", Parameters: [][]AnnParam{
					{{VarID: ptr.Of("X_INTRODUCED_0_")}},
				}},
			},
		},
	},
	SolveGoals: []SolveGoal{{
		SolveMethod: SolveMethodMaximize,
		Objective:   BasicExpr{Identifier: "X_INTRODUCED_0_"},
	}},
}

func TestParseModel(t *testing.T) {
	testCases := []struct {
		desc    string
		input   string
		want    *Model
		wantErr bool
	}{
		{
			desc:  "cakes.fzn",
			input: testCakesFZN,
			want:  &testCakesModel,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, gotErr := ParseModel(strings.NewReader(tc.input))

			if tc.wantErr && gotErr == nil {
				t.Errorf("ParseModel(): want error, got nil")
			}
			if !tc.wantErr && gotErr != nil {
				t.Errorf("ParseModel(): want no error, got %s", gotErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("ParseModel(): mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
