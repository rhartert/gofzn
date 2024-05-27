package fzn

type Predicate struct {
	Identifier string
	Parameters []PredParam
}

type PredParam struct {
	Identifier string
	Array      *Array
	VarType    VarType
	ParType    ParType
}

type ParamDeclaration struct {
	Identifier string
	Type       ParType
	Array      *Array
	Literals   []Literal
}

type ParType int

const (
	ParTypeUnknown ParType = iota
	ParTypeInt
	ParTypeBool
	ParTypeFloat
	ParTypeSetOfInt
)

type VarDeclaration struct {
	Identifier  string
	Variable    Variable
	Array       *Array
	Annotations []Annotation
	Exprs       []BasicExpr
}

type Variable struct {
	Type        VarType
	IntDomain   *SetIntLit
	FloatDomain *SetFloatLit
}

type VarType int

const (
	VarTypeUnknown VarType = iota
	VarTypeIntRange
	VarTypeIntSet
	VarTypeFloatRange
	VarTypeBool
)

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
