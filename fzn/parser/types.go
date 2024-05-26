package parser

type Predicate struct {
	Value string
}

type Parameter struct{}

type Variable struct{}
type Constraint struct {
	Identifier  string
	Expressions []Expr
	Annotations []Annotation
}

type Expr struct {
	IsArray bool
	Exprs   []BasicExpr
}

type SolveMethod int

const (
	SolveMethodSatisfy SolveMethod = iota
	SolveMethodMinimize
	SolveMethodMaximize
)

type SolveGoal struct {
	SolveMethod SolveMethod
	Objective   BasicExpr
	Annotations []Annotation
}

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
