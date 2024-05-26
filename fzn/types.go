package fzn

type Predicate struct {
	Value string
}

type Parameter struct {
	Identifier string
	Type       ParType
	Array      *Array
	Exprs      []BasicLitExpr
}

type ParType int

const (
	ParTypeUnknown ParType = iota
	ParTypeInt
	ParTypeBool
	ParTypeFloat
	ParTypeSetOfInt
)

type Variable struct {
	Identifier  string
	Type        VarType
	Domain      VarDomain
	Array       *Array
	Annotations []Annotation
	Exprs       []BasicExpr
}

type VarType int

const (
	VarTypeUnknown VarType = iota
	VarTypeIntRange
	VarTypeIntSet
	VarTypeFloatRange
	VarTypeBool
)

type VarDomain struct {
	IntDomain   *SetIntLit
	FloatDomain *SetFloatLit
}

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
