package fzn

type Predicate struct {
	Value string
}

type Parameter struct {
	Identifier string
	Type       ParType
	Array      *Array
	Exprs      []Literal
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

// Annotation represents an annotation which is an identifier or a function call
// with a list of lists of parameters.
type Annotation struct {
	Identifier string
	Parameters [][]AnnParam
}

// AnnParam represents an Annotation parameter. It can either be a Literal, a
// variable identifier, a string literal or nested annotation.
type AnnParam struct {
	Literal    *Literal
	VarID      *string
	StringLit  *string
	Annotation *Annotation
}

type BasicExpr struct {
	Identifier string
	Literal    Literal
}

// Literal represents a literal which can either be an int, a bool, a float,
// or a set.
type Literal struct {
	Int      *int
	Bool     *bool
	Float    *float64
	SetInt   *SetIntLit
	SetFloat *SetFloatLit
}

type SetIntLit struct {
	Values [][]int
}

type SetFloatLit struct {
	Values [][]float64
}

type Array struct {
	Start int
	End   int
}
