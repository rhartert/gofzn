package tok

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testCase struct {
	input   string  // the string to tokenize
	want    []Token // excluding last EOF token
	wantErr bool    // whether an error should be returned or not
}

var testCases = []testCase{
	// Valid decimal IntLit.
	{"0", []Token{{IntLit, "0"}}, false},
	{"-0", []Token{{IntLit, "-0"}}, false},
	{"1", []Token{{IntLit, "1"}}, false},
	{"-1", []Token{{IntLit, "-1"}}, false},
	{"012345", []Token{{IntLit, "012345"}}, false},
	{"-012345", []Token{{IntLit, "-012345"}}, false},

	// Valid Hexadecimal IntLit.
	{"0x0", []Token{{IntLit, "0x0"}}, false},
	{"-0x0", []Token{{IntLit, "-0x0"}}, false},
	{"0x1", []Token{{IntLit, "0x1"}}, false},
	{"0x123abc45def", []Token{{IntLit, "0x123abc45def"}}, false},
	{"0x123ABC45DEF", []Token{{IntLit, "0x123ABC45DEF"}}, false},
	{"0x123AbC45dEf", []Token{{IntLit, "0x123AbC45dEf"}}, false},

	// Valid Octal IntLit.
	{"0o0", []Token{{IntLit, "0o0"}}, false},
	{"-0o0", []Token{{IntLit, "-0o0"}}, false},
	{"0o1", []Token{{IntLit, "0o1"}}, false},
	{"0o1234567", []Token{{IntLit, "0o1234567"}}, false},

	// Valid FloatLit
	{"0.0", []Token{{FloatLit, "0.0"}}, false},
	{"-0.0", []Token{{FloatLit, "-0.0"}}, false},
	{"0.1", []Token{{FloatLit, "0.1"}}, false},
	{"-0.1", []Token{{FloatLit, "-0.1"}}, false},
	{"1.0", []Token{{FloatLit, "1.0"}}, false},
	{"-1.0", []Token{{FloatLit, "-1.0"}}, false},
	{"12.345", []Token{{FloatLit, "12.345"}}, false},
	{"-12.345", []Token{{FloatLit, "-12.345"}}, false},
	{"12.34e5", []Token{{FloatLit, "12.34e5"}}, false},
	{"-12.34e5", []Token{{FloatLit, "-12.34e5"}}, false},
	{"12.34e-5", []Token{{FloatLit, "12.34e-5"}}, false},
	{"12.34e+5", []Token{{FloatLit, "12.34e+5"}}, false},
	{"12.34E5", []Token{{FloatLit, "12.34E5"}}, false},
	{"-12.34E5", []Token{{FloatLit, "-12.34E5"}}, false},
	{"12.34E-5", []Token{{FloatLit, "12.34E-5"}}, false},
	{"12.34E+5", []Token{{FloatLit, "12.34E+5"}}, false},

	// Partial numbers.
	{"1x2345", []Token{{IntLit, "1"}, {Identifier, "x2345"}}, false},
	{"0x", []Token{{IntLit, "0"}, {Identifier, "x"}}, false},
	{"0o", []Token{{IntLit, "0"}, {Identifier, "o"}}, false},
	{"1.2e", []Token{{FloatLit, "1.2"}, {Identifier, "e"}}, false},

	// Invalid numbers.
	{"01234_5", nil, true},
	{"0123+45", nil, true},
	{"+1", nil, true},
	{"-", nil, true},
	{"0.", nil, true},
	{"0xaF.2", nil, true},
	{"1.2e-", nil, true},
	{"1.2e+", nil, true},
	{"1.2.3", nil, true},

	// Ranges.
	{"12..345", []Token{{IntLit, "12"}, {Range, ".."}, {IntLit, "345"}}, false},
	{"123..45", []Token{{IntLit, "123"}, {Range, ".."}, {IntLit, "45"}}, false},
	{"0x1..45", []Token{{IntLit, "0x1"}, {Range, ".."}, {IntLit, "45"}}, false},
	{"0o1..45", []Token{{IntLit, "0o1"}, {Range, ".."}, {IntLit, "45"}}, false},

	// Supported keywords.
	{"array", []Token{{Array, "array"}}, false},
	{"bool", []Token{{BoolType, "bool"}}, false},
	{"constraint", []Token{{Constraint, "constraint"}}, false},
	{"false", []Token{{BoolLit, "false"}}, false},
	{"float", []Token{{FloatType, "float"}}, false},
	{"int", []Token{{IntType, "int"}}, false},
	{"maximize", []Token{{Maximize, "maximize"}}, false},
	{"minimize", []Token{{Minimize, "minimize"}}, false},
	{"of", []Token{{Of, "of"}}, false},
	{"predicate", []Token{{Predicate, "predicate"}}, false},
	{"satisfy", []Token{{Satisfy, "satisfy"}}, false},
	{"set", []Token{{Set, "set"}}, false},
	{"solve", []Token{{Solve, "solve"}}, false},
	{"true", []Token{{BoolLit, "true"}}, false},
	{"var", []Token{{Var, "var"}}, false},

	// Unsupported keywords.
	{"ann", []Token{{Error, "ann"}}, false},
	{"annotation", []Token{{Error, "annotation"}}, false},
	{"any", []Token{{Error, "any"}}, false},
	{"case", []Token{{Error, "case"}}, false},
	{"diff", []Token{{Error, "diff"}}, false},
	{"div", []Token{{Error, "div"}}, false},
	{"else", []Token{{Error, "else"}}, false},
	{"elseif", []Token{{Error, "elseif"}}, false},
	{"endif", []Token{{Error, "endif"}}, false},
	{"enum", []Token{{Error, "enum"}}, false},
	{"function", []Token{{Error, "function"}}, false},
	{"if", []Token{{Error, "if"}}, false},
	{"in", []Token{{Error, "in"}}, false},
	{"include", []Token{{Error, "include"}}, false},
	{"intersect", []Token{{Error, "intersect"}}, false},
	{"let", []Token{{Error, "let"}}, false},
	{"list", []Token{{Error, "list"}}, false},
	{"mod", []Token{{Error, "mod"}}, false},
	{"not", []Token{{Error, "not"}}, false},
	{"op", []Token{{Error, "op"}}, false},
	{"opt", []Token{{Error, "opt"}}, false},
	{"output", []Token{{Error, "output"}}, false},
	{"par", []Token{{Error, "par"}}, false},
	{"record", []Token{{Error, "record"}}, false},
	{"string", []Token{{Error, "string"}}, false},
	{"subset", []Token{{Error, "subset"}}, false},
	{"superset", []Token{{Error, "superset"}}, false},
	{"symdiff", []Token{{Error, "symdiff"}}, false},
	{"test", []Token{{Error, "test"}}, false},
	{"then", []Token{{Error, "then"}}, false},
	{"tuple", []Token{{Error, "tuple"}}, false},
	{"type", []Token{{Error, "type"}}, false},
	{"union", []Token{{Error, "union"}}, false},
	{"where", []Token{{Error, "where"}}, false},
	{"xor", []Token{{Error, "xor"}}, false},

	// Strings.
	{`""`, []Token{{StringLit, `""`}}, false},
	{`" foo  "`, []Token{{StringLit, `" foo  "`}}, false},
	{`"foo bar"`, []Token{{StringLit, `"foo bar"`}}, false},
	{`"foo" "bar"`, []Token{{StringLit, `"foo"`}, {StringLit, `"bar"`}}, false},
	{`"foo\" \\"`, []Token{{StringLit, `"foo\" \\"`}}, false},

	// Invalid strings.
	{`"`, nil, true},
	{`"\"`, nil, true},

	// Arrays.
	{"[]", []Token{{ArrayStart, "["}, {ArrayEnd, "]"}}, false},
	{"[0]", []Token{{ArrayStart, "["}, {IntLit, "0"}, {ArrayEnd, "]"}}, false},
	{"[1, 2]", []Token{{ArrayStart, "["}, {IntLit, "1"}, {Comma, ","}, {IntLit, "2"}, {ArrayEnd, "]"}}, false},

	// Tuples.
	{"()", []Token{{TupleStart, "("}, {TupleEnd, ")"}}, false},
	{"(0)", []Token{{TupleStart, "("}, {IntLit, "0"}, {TupleEnd, ")"}}, false},
	{"(1, 2)", []Token{{TupleStart, "("}, {IntLit, "1"}, {Comma, ","}, {IntLit, "2"}, {TupleEnd, ")"}}, false},

	// Set.
	{"{}", []Token{{SetStart, "{"}, {SetEnd, "}"}}, false},
	{"{0}", []Token{{SetStart, "{"}, {IntLit, "0"}, {SetEnd, "}"}}, false},
	{"{1, 2}", []Token{{SetStart, "{"}, {IntLit, "1"}, {Comma, ","}, {IntLit, "2"}, {SetEnd, "}"}}, false},

	// Comment.
	{
		input: "%% comment line; var :: solve : %&_@!",
		want:  []Token{{Comment, "%% comment line; var :: solve : %&_@!"}},
	},
	{
		input: "%% comment 1\n%% comment 2",
		want:  []Token{{Comment, "%% comment 1"}, {Comment, "%% comment 2"}},
	},

	// Sample instructions.
	{
		input: "  var  bool: X_VAR_::foo :: bar = Y_VAR_;",
		want: []Token{
			{Var, "var"},
			{BoolType, "bool"},
			{Colon, ":"},
			{Identifier, "X_VAR_"},
			{AnnStart, "::"},
			{Identifier, "foo"},
			{AnnStart, "::"},
			{Identifier, "bar"},
			{Assign, "="},
			{Identifier, "Y_VAR_"},
			{EOI, ";"},
		},
	},
}

func TestTokenizer_Tokenize(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			tok := Tokenizer{}
			got, gotErr := tok.Tokenize(tc.input)

			if tc.wantErr && gotErr == nil {
				t.Errorf("Tokenize(%q): want error, got nil", tc.input)
			}
			if !tc.wantErr && gotErr != nil {
				t.Errorf("Tokenize(%q): want no error, got %s", tc.input, gotErr)
			}

			want := tc.want
			if !tc.wantErr {
				want = append(want, Token{EOF, ""})
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("Tokenize(%q): mismatch (-want +got):\n%s", tc.input, diff)
			}
		})
	}
}

func TestToken_String(t *testing.T) {
	testCases := []struct {
		token Token
		want  string
	}{
		{
			token: Token{},
			want:  `Token{Error ""}`,
		},
		{
			token: Token{EOF, ""},
			want:  `Token{EOF ""}`,
		},
		{
			token: Token{IntLit, "123"},
			want:  `Token{IntLit "123"}`,
		},
		{
			token: Token{Type(-1), "foobar"},
			want:  `Token{Type(-1) "foobar"}`,
		},
	}

	for _, tc := range testCases {
		if got := tc.token.String(); got != tc.want {
			t.Errorf("(%#v).String(): want %q, got %q", tc.token, tc.want, got)
		}
	}
}
