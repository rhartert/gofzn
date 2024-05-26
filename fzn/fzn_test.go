package fzn

import (
	_ "embed"
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/google/go-cmp/cmp"
)

type testCaseInstruction struct {
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

func TestParse_error(t *testing.T) {
	testParse(t, []testCaseInstruction{
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
	})
}

func testParse(t *testing.T, testCases []testCaseInstruction) {
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
