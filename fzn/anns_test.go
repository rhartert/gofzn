package fzn

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rhartert/gofzn/fzn/tok"
	"github.com/rhartert/ptr"
)

func TestParserInstruction_parseAnnotations(t *testing.T) {
	testCases := []struct {
		tokens  []tok.Token
		want    []Annotation
		wantErr bool
	}{
		{
			tokens: []tok.Token{
				{Type: tok.AnnStart, Value: "::"},
			},
			wantErr: true,
		},
		{
			tokens: []tok.Token{
				{Type: tok.AnnStart, Value: "::"},
				{Type: tok.AnnStart, Value: "::"},
			},
			wantErr: true,
		},
		{
			tokens: []tok.Token{
				{Type: tok.AnnStart, Value: "::"},
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.AnnStart, Value: "::"},
			},
			wantErr: true,
		},
		{
			tokens: []tok.Token{},
			want:   nil,
		},
		{
			tokens: []tok.Token{
				{Type: tok.AnnStart, Value: "::"},
				{Type: tok.Identifier, Value: "foo"},
			},
			want: []Annotation{
				{Identifier: "foo"},
			},
		},
		{
			tokens: []tok.Token{
				{Type: tok.AnnStart, Value: "::"},
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.AnnStart, Value: "::"},
				{Type: tok.Identifier, Value: "bar"},
			},
			want: []Annotation{
				{Identifier: "foo"},
				{Identifier: "bar"},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			p := &parser{tokens: tc.tokens}

			got, gotErr := parseAnnotations(p)

			if tc.wantErr && gotErr == nil {
				t.Errorf("parseAnnotations(): want error, got nil")
			}
			if !tc.wantErr && gotErr != nil {
				t.Errorf("parseAnnotations(): want no error, got %s", gotErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("parseAnnotations(): mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParserInstruction_parseAnnotation(t *testing.T) {
	testCases := []struct {
		desc    string
		tokens  []tok.Token
		want    Annotation
		wantErr bool
	}{
		{
			desc:    "no token",
			tokens:  []tok.Token{},
			wantErr: true,
		},
		{
			desc: "not starting with identifier (\"::\")",
			tokens: []tok.Token{
				{Type: tok.AnnStart, Value: "::"},
			},
			wantErr: true,
		},
		{
			desc: "not starting with identifier (\"(\")",
			tokens: []tok.Token{
				{Type: tok.TupleStart, Value: "("},
			},
			wantErr: true,
		},
		{
			desc: "missing \")\"",
			tokens: []tok.Token{
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.TupleStart, Value: "("},
			},
			wantErr: true,
		},
		{
			desc: "missing \")\" after expression",
			tokens: []tok.Token{
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.TupleStart, Value: "("},
				{Type: tok.Identifier, Value: "bar"},
				{Type: tok.Comma, Value: ","},
			},
			wantErr: true,
		},
		{
			desc: "valid identifier",
			tokens: []tok.Token{
				{Type: tok.Identifier, Value: "foo"},
			},
			want: Annotation{
				Identifier: "foo",
			},
		},
		{
			desc: "valid empty call",
			tokens: []tok.Token{
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.TupleStart, Value: "("},
				{Type: tok.TupleEnd, Value: ")"},
			},
			want: Annotation{
				Identifier: "foo",
				Exprs:      [][]AnnExpr{},
			},
		},
		{
			desc: "valid call with one IntLit",
			tokens: []tok.Token{
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.TupleStart, Value: "("},
				{Type: tok.IntLit, Value: "42"},
				{Type: tok.TupleEnd, Value: ")"},
			},
			want: Annotation{
				Identifier: "foo",
				Exprs: [][]AnnExpr{{{
					BasicLitExpr: &BasicLitExpr{
						Int: ptr.Of(42),
					},
				}}},
			},
		},
		{
			desc: "valid call with multiple literal expressions",
			tokens: []tok.Token{
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.TupleStart, Value: "("},
				{Type: tok.IntLit, Value: "42"},
				{Type: tok.Comma, Value: ","},
				{Type: tok.Identifier, Value: "bar"},
				{Type: tok.Comma, Value: ","},
				{Type: tok.TupleEnd, Value: ")"},
			},
			want: Annotation{
				Identifier: "foo",
				Exprs: [][]AnnExpr{
					{{BasicLitExpr: &BasicLitExpr{Int: ptr.Of(42)}}},
					{{VarID: ptr.Of("bar")}},
				},
			},
		},
		{
			desc: "valid nested annotation",
			tokens: []tok.Token{
				{Type: tok.Identifier, Value: "foo"},
				{Type: tok.TupleStart, Value: "("},
				{Type: tok.IntLit, Value: "42"},
				{Type: tok.Comma, Value: ","},
				{Type: tok.Identifier, Value: "bar"},
				{Type: tok.TupleStart, Value: "("},
				{Type: tok.IntLit, Value: "1337"},
				{Type: tok.TupleEnd, Value: ")"},
				{Type: tok.Comma, Value: ","},
				{Type: tok.TupleEnd, Value: ")"},
			},
			want: Annotation{
				Identifier: "foo",
				Exprs: [][]AnnExpr{
					{{BasicLitExpr: &BasicLitExpr{Int: ptr.Of(42)}}},
					{{Annotation: &Annotation{
						Identifier: "bar",
						Exprs: [][]AnnExpr{{
							{BasicLitExpr: &BasicLitExpr{Int: ptr.Of(1337)}},
						}},
					}}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			p := &parser{tokens: tc.tokens}

			got, gotErr := parseAnnotation(p)

			if tc.wantErr && gotErr == nil {
				t.Errorf("parseAnnotation(): want error, got nil")
			}
			if !tc.wantErr && gotErr != nil {
				t.Errorf("parseAnnotation(): want no error, got %s", gotErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("parseAnnotation(): mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
