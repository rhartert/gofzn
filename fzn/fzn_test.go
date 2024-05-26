package fzn

import (
	"errors"
	"io"
	"testing"
	"testing/iotest"

	"github.com/google/go-cmp/cmp"
)

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
