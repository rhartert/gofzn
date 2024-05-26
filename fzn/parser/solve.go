package parser

import (
	"fmt"

	tok "github.com/rhartert/gofzn/fzn/tokenizer"
)

// Parsers for SolveGoals
// ----------------------
//
// Grammar:
//
//  <solve-item> ::= "solve" <annotations> "satisfy" ";"
//                 | "solve" <annotations> "minimize" <basic-expr> ";"
//                 | "solve" <annotations> "maximize" <basic-expr> ";"
//

func isSolveGoal(p *parser) bool {
	return p.lookAhead(0).Type == tok.Solve
}

func parseSolveGoal(p *parser) (*SolveGoal, error) {
	if !p.nextIf(tok.Solve) {
		return nil, fmt.Errorf("solve declaration should start with \"solve\"")
	}

	anns, err := parseAnnotations(p)
	if err != nil {
		return nil, fmt.Errorf("error parsing solve annotations: %w", err)
	}

	sg := &SolveGoal{
		Annotations: anns,
	}

	switch t := p.next(); t.Type {
	case tok.Satisfy:
		sg.SolveMethod = SolveMethodSatisfy
	case tok.Minimize:
		sg.SolveMethod = SolveMethodMinimize
	case tok.Maximize:
		sg.SolveMethod = SolveMethodMaximize
	default:
		return nil, fmt.Errorf("invalid solve method: [%s] %q", t.Type, t.Value)
	}

	// No objective to parse for satisfy goals.
	if sg.SolveMethod == SolveMethodSatisfy {
		if !p.nextIf(tok.EOI) {
			return nil, fmt.Errorf("missing end of solve declaration ';'")
		}
		return sg, nil
	}

	expr, err := parseBasicExpr(p)
	if err != nil {
		return nil, fmt.Errorf("error parsing solve objective: %w", err)
	}

	sg.Objective = expr
	return sg, nil
}
