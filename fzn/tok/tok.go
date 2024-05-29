// Package tok provides the Tokenizer struct to parse strings into sequences
// of FlatZinc lexical tokens.
package tok

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const eof = rune(0)

//go:generate stringer -type=Type
type Type int

const (
	Error Type = iota
	AnnStart
	Array
	ArrayEnd
	ArrayStart
	Assign
	BoolLit
	BoolType
	Colon
	Comma
	Comment
	Constraint
	EOF
	EOI
	FloatLit
	FloatType
	Identifier
	IntLit
	IntType
	Maximize
	Minimize
	Of
	Predicate
	Range
	Satisfy
	Set
	SetEnd
	SetStart
	Solve
	StringLit
	TupleEnd
	TupleStart
	Var
)

type Token struct {
	Type  Type
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("Token{%s %q}", t.Type, t.Value)
}

type Tokenizer struct {
	input  string  // the string being tokenized
	start  int     // start position of the token being parsed
	pos    int     // current position in the input
	width  int     // width of the last rune read
	tokens []Token // in order sequence of parsed tokens
}

// Tokenize returns the in-order sequence of tokens extracted from the input
// string.
func (t *Tokenizer) Tokenize(input string) ([]Token, error) {
	t.input = input                   // set the new input
	t.start, t.pos, t.width = 0, 0, 0 // reset the tokenizer
	t.tokens = t.tokens[:0]           // avoid reallocating a new slice

	t.run() // start the state machine

	// If a parsing error occurred, the last token will be an Error token with
	// the error message as value.
	if last := len(t.tokens) - 1; last >= 0 && t.tokens[last].Type == Error {
		return nil, fmt.Errorf(t.tokens[last].Value)
	}
	return t.tokens, nil
}

// errorf returns an error token and terminates the tokenization by passing
// back a nil pointer that will be the next state, terminating t.run.
func (t *Tokenizer) errorf(format string, args ...any) stateFn {
	t.tokens = append(t.tokens, Token{Error, fmt.Sprintf(format, args...)})
	return nil
}

// emit adds a new token of the given type with the current run (i.e. the
// runes between t.start and t.pos) to t.tokens. It updates t.start to the
// current position t.pos to be ready for the next token.
func (t *Tokenizer) emit(tt Type) {
	t.tokens = append(t.tokens, Token{tt, t.input[t.start:t.pos]})
	t.start = t.pos
}

// next returns the next rune in the input.
func (t *Tokenizer) next() (r rune) {
	if t.pos >= len(t.input) {
		t.width = 0
		return eof
	}
	r, t.width = utf8.DecodeRuneInString(t.input[t.pos:])
	t.pos += t.width
	return r
}

// backup steps back one rune, it can be called only once per call of next.
func (t *Tokenizer) backup() {
	t.pos -= t.width
}

// trim drops the current run and skips all the following space runes until a
// non-space rune is found or the end of t.input is reached.
func (t *Tokenizer) trim() {
	for unicode.IsSpace(t.next()) {
		t.start = t.pos // ignore the last rune
	}
	t.backup()
}

// accept returns true and consumes the next rune if it is part of the valid
// runes. It returns false and does nothing otherwise.
func (t *Tokenizer) accept(valid string) bool {
	if strings.ContainsRune(valid, t.next()) {
		return true
	}
	t.backup()
	return false
}

// acceptRun consumes the next runes until it encounters a rune that is not in
// valid. It is equivalent to calling t.accept until it returns false.
func (t *Tokenizer) acceptRun(valid string) {
	for strings.ContainsRune(valid, t.next()) {
		// nothing
	}
	t.backup()
}

// State machine
// -------------

type stateFn func(*Tokenizer) stateFn

func (t *Tokenizer) run() {
	for state := tokenizeAnything; state != nil; {
		state = state(t)
	}
}

// tokenizeAnything represents the generic state in which no assumption is made
// on the nature of the next token in the sequence.
func tokenizeAnything(t *Tokenizer) stateFn {
	for {
		t.trim()

		switch s := t.input[t.pos:]; {
		case strings.HasPrefix(s, "%"):
			return tokenizeComment
		case strings.HasPrefix(s, "\""):
			return tokenizeString
		case strings.HasPrefix(s, ".."):
			return tokenizeRange
		case strings.HasPrefix(s, "::"):
			return tokenizeAnnotation
		}

		r := t.next()
		if r == eof {
			break
		}

		switch {
		case unicode.IsLetter(r):
			t.backup()
			return tokenizeIdentifierOrKeyword
		case unicode.IsNumber(r) || r == '-':
			t.backup()
			return tokenizeNumber
		}

		tt := tokenTypeFromRune(r)
		if tt == Error {
			return t.errorf("unexpected rune %q (%v)", string(r), r)
		}
		t.emit(tt)
	}

	t.emit(EOF) // cleanly indicate the end of the input
	return nil  // stops the state machine
}

func tokenTypeFromRune(r rune) Type {
	switch r {
	case ':':
		return Colon
	case '{':
		return SetStart
	case '}':
		return SetEnd
	case '[':
		return ArrayStart
	case ']':
		return ArrayEnd
	case '(':
		return TupleStart
	case ')':
		return TupleEnd
	case ',':
		return Comma
	case '=':
		return Assign
	case ';':
		return EOI
	default:
		return Error
	}
}

// tokenizeComment parses a Comment token made of all the incoming runes until
// a new line rune is found.
func tokenizeComment(t *Tokenizer) stateFn {
	for {
		r := t.next()
		if r == eof {
			t.emit(Comment)
			return tokenizeAnything
		}
		if r == '\n' {
			t.backup() // do not include the \n rune
			t.emit(Comment)
			t.trim() // drop the \n rune
			return tokenizeAnything
		}
	}
}

// tokenizeAnnotation parses a AnnStart token.
func tokenizeAnnotation(t *Tokenizer) stateFn {
	t.pos += len("::")
	t.emit(AnnStart)
	t.trim() // remove optional spaces before identifier
	return tokenizeIdentifierOrKeyword
}

// tokenizeRange parses a Range token.
func tokenizeRange(t *Tokenizer) stateFn {
	t.pos += len("..")
	t.emit(Range)
	return tokenizeNumber
}

// tokenizeString parses a String token.
func tokenizeString(t *Tokenizer) stateFn {
	if !t.accept("\"") {
		return t.errorf("string should start with \"")
	}

	escape := false
	for {
		r := t.next()
		if r == eof {
			return t.errorf("string wasn't closed")
		}
		if r == '"' && !escape {
			break
		}
		if r == '\\' && !escape {
			escape = true
			continue
		}
		escape = false
	}

	t.emit(StringLit)
	return tokenizeAnything
}

var keywords = map[string]Type{
	"array":      Array,
	"bool":       BoolType,
	"constraint": Constraint,
	"false":      BoolLit,
	"float":      FloatType,
	"int":        IntType,
	"maximize":   Maximize,
	"minimize":   Minimize,
	"of":         Of,
	"predicate":  Predicate,
	"satisfy":    Satisfy,
	"set":        Set,
	"solve":      Solve,
	"true":       BoolLit,
	"var":        Var,

	// Reserved MiniZinc keywords that are not part of FlatZinc.
	"ann":        Error,
	"annotation": Error,
	"any":        Error,
	"case":       Error,
	"diff":       Error,
	"div":        Error,
	"else":       Error,
	"elseif":     Error,
	"endif":      Error,
	"enum":       Error,
	"function":   Error,
	"if":         Error,
	"in":         Error,
	"include":    Error,
	"intersect":  Error,
	"let":        Error,
	"list":       Error,
	"mod":        Error,
	"not":        Error,
	"op":         Error,
	"opt":        Error,
	"output":     Error,
	"par":        Error,
	"record":     Error,
	"string":     Error,
	"subset":     Error,
	"superset":   Error,
	"symdiff":    Error,
	"test":       Error,
	"then":       Error,
	"tuple":      Error,
	"type":       Error,
	"union":      Error,
	"where":      Error,
	"xor":        Error,
}

// tokenizeIdentifierOrKeyword parses either a tIdentifier token or one of the
// reserved keyword tokens defined in keywords.
func tokenizeIdentifierOrKeyword(t *Tokenizer) stateFn {
	r := t.next()
	for unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' {
		r = t.next()
	}
	t.backup()

	s := t.input[t.start:t.pos]
	if tt, ok := keywords[s]; ok {
		t.emit(tt) // the identifier is a reserved keyword
	} else {
		t.emit(Identifier)
	}

	return tokenizeAnything
}

const (
	digitsOctal       = "01234567"
	digitsDecimal     = "0123456789"
	digitsHexadecimal = "0123456789abcdefABCDEF"
)

// tokenizeNumber parses either an IntLit or FloatLit token. The function
// assumes that the tokenizer is positioned on the first rune of a valid
// number as defined by the following expressions:
//
//	IntLit   ::= [-]?[0-9]+
//	           | [-]?0x[0-9A-Fa-f]+
//	           | [-]?0o[0-7]+
//
//	FloatLit ::= [-]?[0-9]+.[0-9]+
//	           | [-]?[0-9]+.[0-9]+[Ee][-+]?[0-9]+
//	           | [-]?[0-9]+[Ee][-+]?[0-9]+
func tokenizeNumber(t *Tokenizer) stateFn {
	t.accept("-")

	r := t.next()
	if !unicode.IsNumber(r) {
		return t.errorf("invalid number")
	}

	// Parsing FlatZinc numbers requires either a two-runes look ahead or
	// bactkracking (see "0x" or "1.2Ea"). The function relies on lastOK to
	// keep track of the last rune at which the current run was a valid number.
	// This enables easy backtracking by resetting the t.pos to lastOK.
	lastOK := t.pos

	digits := digitsDecimal
	if r == '0' {
		switch {
		case t.accept("x"):
			digits = digitsHexadecimal
		case t.accept("o"):
			digits = digitsOctal
		}
	}

	// Verify that there is at least one digit after 'x' or 'o' for
	// hexadecimal and octal numbers. If not, this means that the number
	// stopped before reaching 'x' or 'o'.
	if digits != digitsDecimal && !t.accept(digits) {
		goto backtrackEmitInt
	}

	t.acceptRun(digits)

	// This point marks the end of IntLit. The rest of the function attempts
	// to keep parsing the stream of runes to build a FloatInt. If it fails,
	// it backtracks and emits an IntLit.
	lastOK = t.pos

	if !t.accept(".") {
		goto backtrackEmitInt
	}
	if digits != digitsDecimal { // FloatLit only contains decimal numbers
		goto backtrackEmitInt
	}
	if !t.accept(digitsDecimal) {
		goto backtrackEmitInt
	}

	t.acceptRun(digitsDecimal)
	lastOK = t.pos // first valid float

	if t.accept("eE") {
		t.accept("-+")
		if t.accept(digitsDecimal) {
			t.acceptRun(digitsDecimal)
		} else {
			t.pos = lastOK // missing digit after exponent
		}
	}

	t.emit(FloatLit)
	return tokenizeAnything

backtrackEmitInt:
	t.pos = lastOK
	t.emit(IntLit)
	return tokenizeAnything
}
