package parser

type Predicate struct{}
type Parameter struct{}
type Variable struct{}
type Constraint struct{}
type SolveGoal struct{}

type Annotation struct {
	Identifier string
	Exprs      [][]AnnExpr
}

type AnnExpr struct {
	BasicLitExpr *BasicLitExpr
	VarID        *string
	StringLit    *string
	Annotation   *Annotation
}

type BasicExpr struct {
	Identifier  string
	LiteralExpr BasicLitExpr
}

type BasicLitExpr struct {
	Int   *int
	Bool  *bool
	Float *float64
	Set   *SetLit
}

type SetLit struct {
	SetInt   *SetIntLit
	SetFloat *SetFloatLit
}

type SetIntLit struct {
	Values [][]int
}

type SetFloatLit struct {
	Values [][]float64
}

type RangeInt struct {
	Min, Max int
}

type RangeFloat struct {
	Min, Max float64
}
