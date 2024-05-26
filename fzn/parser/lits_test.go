package parser

import (
	"testing"

	tok "github.com/rhartert/gofzn/fzn/tokenizer"
)

func TestParserInstruction_parseBoolLit(t *testing.T) {
	testCases := []struct {
		tokens  []tok.Token
		want    bool
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
		{
			tokens:  []tok.Token{{Type: tok.EOF}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.EOI}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.BoolLit}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.BoolLit}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.IntLit, Value: "true"}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.BoolLit, Value: "foo"}},
			wantErr: true,
		},
		{
			tokens: []tok.Token{{Type: tok.BoolLit, Value: "true"}},
			want:   true,
		},
		{
			tokens: []tok.Token{{Type: tok.BoolLit, Value: "false"}},
			want:   false,
		},
	}

	for _, tc := range testCases {
		p := &parser{tokens: tc.tokens}

		got, gotErr := parseBoolLit(p)

		if tc.wantErr && gotErr == nil {
			t.Errorf("parseBoolLit: want error, got nil")
		}
		if !tc.wantErr && gotErr != nil {
			t.Errorf("parseBoolLit: want no error, got %s", gotErr)
		}
		if tc.want != got {
			t.Errorf("parseBoolLit: want %v, got %v", tc.want, got)
		}
	}
}

func TestParserInstruction_parseIntLit(t *testing.T) {
	testCases := []struct {
		tokens  []tok.Token
		want    int
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
		{
			tokens:  []tok.Token{{Type: tok.EOF}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.EOI}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.BoolLit}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.IntLit}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.IntLit, Value: "foo"}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.IntLit, Value: "42foo"}},
			wantErr: true,
		},
		{
			tokens: []tok.Token{{Type: tok.IntLit, Value: "42"}},
			want:   42,
		},
		{
			tokens: []tok.Token{{Type: tok.IntLit, Value: "0x2A"}},
			want:   42,
		},
		{
			tokens: []tok.Token{{Type: tok.IntLit, Value: "0o52"}},
			want:   42,
		},
	}

	for _, tc := range testCases {
		p := &parser{tokens: tc.tokens}

		got, gotErr := parseIntLit(p)

		if tc.wantErr && gotErr == nil {
			t.Errorf("parseIntLit: want error, got nil")
		}
		if !tc.wantErr && gotErr != nil {
			t.Errorf("parseIntLit: want no error, got %s", gotErr)
		}
		if tc.want != got {
			t.Errorf("parseIntLit: want %d, got %d", tc.want, got)
		}
	}
}

func TestParserInstruction_parseFloatLit(t *testing.T) {
	testCases := []struct {
		tokens  []tok.Token
		want    float64
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
		{
			tokens:  []tok.Token{{Type: tok.EOF}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.EOI}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.BoolLit}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.FloatLit}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.FloatLit, Value: "foo"}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.FloatLit, Value: "42foo"}},
			wantErr: true,
		},
		{
			tokens: []tok.Token{{Type: tok.FloatLit, Value: "42.0"}},
			want:   42,
		},
		{
			tokens: []tok.Token{{Type: tok.FloatLit, Value: "420E-1"}},
			want:   42,
		},
		{
			tokens: []tok.Token{{Type: tok.FloatLit, Value: "0.042e3"}},
			want:   42,
		},
	}

	for _, tc := range testCases {
		p := &parser{tokens: tc.tokens}

		got, gotErr := parseFloatLit(p)

		if tc.wantErr && gotErr == nil {
			t.Errorf("parseFloatLit: want error, got nil")
		}
		if !tc.wantErr && gotErr != nil {
			t.Errorf("parseFloatLit: want no error, got %s", gotErr)
		}
		if tc.want != got {
			t.Errorf("parseFloatLit: want %f, got %f", tc.want, got)
		}
	}
}

func TestParserInstruction_parseIdentifier(t *testing.T) {
	testCases := []struct {
		tokens  []tok.Token
		want    string
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
		{
			tokens:  []tok.Token{{Type: tok.EOF}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.EOI}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.BoolLit}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.Identifier, Value: ""}},
			wantErr: true,
		},
		{
			tokens: []tok.Token{{Type: tok.Identifier, Value: "foo_bar42_"}},
			want:   "foo_bar42_",
		},
	}

	for _, tc := range testCases {
		p := &parser{tokens: tc.tokens}

		got, gotErr := parseIdentifier(p)

		if tc.wantErr && gotErr == nil {
			t.Errorf("parseIdentifier: want error, got nil")
		}
		if !tc.wantErr && gotErr != nil {
			t.Errorf("parseIdentifier: want no error, got %s", gotErr)
		}
		if tc.want != got {
			t.Errorf("parseIdentifier: want %s, got %s", tc.want, got)
		}
	}
}

func TestParserInstruction_parseStringLit(t *testing.T) {
	testCases := []struct {
		tokens  []tok.Token
		want    string
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
		{
			tokens:  []tok.Token{{Type: tok.EOF}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.EOI}},
			wantErr: true,
		},
		{
			tokens:  []tok.Token{{Type: tok.BoolLit}},
			wantErr: true,
		},
		{
			tokens: []tok.Token{{Type: tok.StringLit, Value: ""}},
			want:   "",
		},
		{
			tokens: []tok.Token{{Type: tok.StringLit, Value: "foobar42"}},
			want:   "foobar42",
		},
	}

	for _, tc := range testCases {
		p := &parser{tokens: tc.tokens}

		got, gotErr := parseStringLit(p)

		if tc.wantErr && gotErr == nil {
			t.Errorf("parseStringLit: want error, got nil")
		}
		if !tc.wantErr && gotErr != nil {
			t.Errorf("parseStringLit: want no error, got %s", gotErr)
		}
		if tc.want != got {
			t.Errorf("parseStringLit: want %s, got %s", tc.want, got)
		}
	}
}
