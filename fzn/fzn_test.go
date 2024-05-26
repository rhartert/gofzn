package fzn

import (
	_ "embed"
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/google/go-cmp/cmp"
	"github.com/rhartert/ptr"
)

//go:embed testdata/cakes.fzn
var testCakesFZN string

type testCase struct {
	desc    string
	input   io.Reader
	want    *Model
	wantErr bool
}

var testCases = []testCase{
	{
		desc:    "error reader",
		input:   iotest.ErrReader(errors.New("test error")),
		wantErr: true,
	},
	{
		desc:    "invalid input (tokenizer)",
		input:   strings.NewReader("@@foo-bar%**!!"),
		wantErr: true,
	},
	{
		desc:    "incorrect input (parser)",
		input:   strings.NewReader("var [[ array ))"),
		wantErr: true,
	},
	{
		desc:  "cake.fzn",
		input: strings.NewReader(testCakesFZN),
		want: &Model{
			Parameters: []Parameter{
				{
					Identifier: "X_INTRODUCED_2_",
					Array:      &Array{Start: 1, End: 2},
					Type:       ParTypeInt,
					Exprs: []BasicLitExpr{
						{Int: ptr.Of(250)},
						{Int: ptr.Of(200)},
					},
				},
				{
					Identifier: "X_INTRODUCED_6_",
					Array:      &Array{Start: 1, End: 2},
					Type:       ParTypeInt,
					Exprs: []BasicLitExpr{
						{Int: ptr.Of(75)},
						{Int: ptr.Of(150)},
					},
				},
				{
					Identifier: "X_INTRODUCED_8_",
					Array:      &Array{Start: 1, End: 2},
					Type:       ParTypeInt,
					Exprs: []BasicLitExpr{
						{Int: ptr.Of(100)},
						{Int: ptr.Of(150)},
					},
				},
			},
			Variables: []Variable{
				{
					Identifier:  "b",
					Type:        VarTypeIntRange,
					Domain:      VarDomain{IntDomain: &SetIntLit{Values: [][]int{{0, 3}}}},
					Annotations: []Annotation{{Identifier: "output_var"}},
				},
				{
					Identifier:  "c",
					Type:        VarTypeIntRange,
					Domain:      VarDomain{IntDomain: &SetIntLit{Values: [][]int{{0, 6}}}},
					Annotations: []Annotation{{Identifier: "output_var"}},
				},
				{
					Identifier:  "X_INTRODUCED_0_",
					Type:        VarTypeIntRange,
					Domain:      VarDomain{IntDomain: &SetIntLit{Values: [][]int{{0, 85000}}}},
					Annotations: []Annotation{{Identifier: "is_defined_var"}},
				},
			},
			Constraints: []Constraint{
				{
					Identifier: "int_lin_le",
					Expressions: []Expr{
						{Exprs: []BasicExpr{{Identifier: "X_INTRODUCED_2_"}}},
						{IsArray: true, Exprs: []BasicExpr{{Identifier: "b"}, {Identifier: "c"}}},
						{Exprs: []BasicExpr{{LiteralExpr: BasicLitExpr{Int: ptr.Of(4000)}}}},
					},
				},
				{
					Identifier: "int_lin_le",
					Expressions: []Expr{
						{Exprs: []BasicExpr{{Identifier: "X_INTRODUCED_6_"}}},
						{IsArray: true, Exprs: []BasicExpr{{Identifier: "b"}, {Identifier: "c"}}},
						{Exprs: []BasicExpr{{LiteralExpr: BasicLitExpr{Int: ptr.Of(2000)}}}},
					},
				},
				{
					Identifier: "int_lin_le",
					Expressions: []Expr{
						{Exprs: []BasicExpr{{Identifier: "X_INTRODUCED_8_"}}},
						{IsArray: true, Exprs: []BasicExpr{{Identifier: "b"}, {Identifier: "c"}}},
						{Exprs: []BasicExpr{{LiteralExpr: BasicLitExpr{Int: ptr.Of(500)}}}},
					},
				},
				{
					Identifier: "int_lin_eq",
					Expressions: []Expr{
						{IsArray: true, Exprs: []BasicExpr{
							{LiteralExpr: BasicLitExpr{Int: ptr.Of(400)}},
							{LiteralExpr: BasicLitExpr{Int: ptr.Of(450)}},
							{LiteralExpr: BasicLitExpr{Int: ptr.Of(-1)}},
						}},
						{IsArray: true, Exprs: []BasicExpr{
							{Identifier: "b"},
							{Identifier: "c"},
							{Identifier: "X_INTRODUCED_0_"},
						}},
						{Exprs: []BasicExpr{{LiteralExpr: BasicLitExpr{Int: ptr.Of(0)}}}},
					},
					Annotations: []Annotation{
						{Identifier: "ctx_pos"},
						{Identifier: "defines_var", Exprs: [][]AnnExpr{
							{{VarID: ptr.Of("X_INTRODUCED_0_")}},
						}},
					},
				},
			},
			SolveGoals: []SolveGoal{{
				SolveMethod: SolveMethodMaximize,
				Objective:   BasicExpr{Identifier: "X_INTRODUCED_0_"},
			}},
		},
	},
}

func TestParseModel(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, gotErr := ParseModel(tc.input)

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
