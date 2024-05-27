package fzn

import (
	"fmt"

	"github.com/rhartert/gofzn/fzn/tok"
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
		return nil, fmt.Errorf("invalid solve method %s", t)
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

	if !p.nextIf(tok.EOI) {
		return nil, fmt.Errorf("missing end of instruction")
	}

	sg.Objective = expr
	return sg, nil
}
