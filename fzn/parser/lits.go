package parser

import (
	"fmt"
	"strconv"

	"github.com/rhartert/gofzn/fzn/tok"
)

// Parsers for terminal literals
// -----------------------------
//
// Grammar:
//
//  <bool-literal>    ::= "false"
//                      | "true"
//
//  <int-literal>     ::= [-]?[0-9]+
//                      | [-]?0x[0-9A-Fa-f]+
//                      | [-]?0o[0-7]+
//
//  <float-literal>   ::= [-]?[0-9]+.[0-9]+
//                      | [-]?[0-9]+.[0-9]+[Ee][-+]?[0-9]+
//                      | [-]?[0-9]+[Ee][-+]?[0-9]+
//
//  <string-literal>  ::= """   """
//
//  <identifier>      ::= [A-Za-z_][A-Za-z0-9_]*

// parseBoolLit parses a bool. The function expects the parser to be positioned
// on a BoolLit token.
func parseBoolLit(p *parser) (bool, error) {
	t := p.next()
	if t.Type != tok.BoolLit {
		return false, fmt.Errorf("not a BoolLit token %s", t)
	}
	switch t.Value {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("invalid BoolLit token %s", t)
	}
}

// parseIntLit parses an int. The function expects the parser to be positioned
// on a tIntLit token.
func parseIntLit(p *parser) (int, error) {
	t := p.next()
	if t.Type != tok.IntLit {
		return 0, fmt.Errorf("not a IntLit token %s", t)
	}
	i, err := strconv.ParseInt(t.Value, 0, 0)
	if err != nil {
		return 0, fmt.Errorf("invalid IntLit token %s: %w", t, err)
	}
	return int(i), nil
}

// parseFloatLit parses an int. The function expects the parser to be positioned
// on a tFloatLit token.
func parseFloatLit(p *parser) (float64, error) {
	t := p.next()
	if t.Type != tok.FloatLit {
		return 0, fmt.Errorf("not a FloatLit token %s", t)
	}
	f, err := strconv.ParseFloat(t.Value, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid FloatLit token %s: %w", t, err)
	}
	return f, nil
}

func isStringLit(p *parser) bool {
	return p.lookAhead(0).Type == tok.StringLit
}

func parseStringLit(p *parser) (string, error) {
	t := p.next()
	if t.Type != tok.StringLit {
		return "", fmt.Errorf("not a string token %s", t)
	}
	return t.Value, nil
}

func isIdentifier(p *parser) bool {
	return p.lookAhead(0).Type == tok.Identifier
}

func parseIdentifier(p *parser) (string, error) {
	t := p.next()
	if t.Type != tok.Identifier {
		return "", fmt.Errorf("not an identifier token %s", t)
	}
	if t.Value == "" {
		return "", fmt.Errorf("empty identifier %s", t)
	}
	return t.Value, nil
}
