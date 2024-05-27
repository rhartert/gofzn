package fzn

// Predicate represents a FlatZinc predicate.
type Predicate struct {
	Identifier string      // Name of the predicate.
	Parameters []PredParam // List of parameters for the predicate.
}

// PredParam represents a parameter in a predicate.
type PredParam struct {
	Identifier string  // Name of the parameter.
	Array      *Array  // Optional array information.
	VarType    VarType // Variable type of the parameter (Unknown if none).
	ParType    ParType // Parameter type of the parameter (Unknown if none).
}

// ParamDeclaration represents a parameter declaration in FlatZinc.
type ParamDeclaration struct {
	Identifier string    // Name of the parameter.
	Type       ParType   // Type of the parameter.
	Array      *Array    // Optional array information.
	Literals   []Literal // List of literals associated with the parameter.
}

// ParType represents the type of a parameter in FlatZinc.
type ParType int

const (
	ParTypeUnknown ParType = iota
	ParTypeInt
	ParTypeBool
	ParTypeFloat
	ParTypeSetOfInt
)

// VarDeclaration represents a variable declaration in FlatZinc.
type VarDeclaration struct {
	Identifier  string       // Name of the variable.
	Variable    Variable     // Variable information.
	Array       *Array       // Optional array information.
	Annotations []Annotation // List of annotations associated with the variable.
	Exprs       []BasicExpr  // List of basic expressions associated with the variable.
}

// Variable represents a variable in FlatZinc.
type Variable struct {
	Type        VarType      // Type of the variable.
	IntDomain   *SetIntLit   // Integer domain of the variable, if applicable.
	FloatDomain *SetFloatLit // Float domain of the variable, if applicable.
}

// VarType represents the type of a variable in FlatZinc.
type VarType int

const (
	VarTypeUnknown VarType = iota
	VarTypeIntRange
	VarTypeIntSet
	VarTypeFloatRange
	VarTypeBool
)

// Constraint represents a FlatZinc constraint.
type Constraint struct {
	Identifier  string       // Name of the constraint.
	Expressions []Expr       // List of expressions (e.g. variable identifiers).
	Annotations []Annotation // List of annotations.
}

// Expr represents an expression in FlatZinc, which is either a single basic
// expression or an array of basic expressions.
type Expr struct {
	Expr  *BasicExpr  // Single basic expression
	Exprs []BasicExpr // List of basic expressions.
}

// SolveMethod represents the method to solve a FlatZinc model.
type SolveMethod int

const (
	SolveMethodSatisfy SolveMethod = iota
	SolveMethodMinimize
	SolveMethodMaximize
)

// SolveGoal represents the FlatZinc solve goal.
type SolveGoal struct {
	SolveMethod SolveMethod  // Method to solve the model.
	Objective   BasicExpr    // Objective expression for optimization.
	Annotations []Annotation // List of solve goal annotations.
}

// Annotation represents an annotation which is an identifier or a function call
// with a list of lists of parameters.
type Annotation struct {
	Identifier string       // Name of the annotation.
	Parameters [][]AnnParam // List of lists of parameters for the annotation.
}

// AnnParam represents an Annotation parameter. It can either be a Literal, a
// variable identifier, a string literal or nested annotation.
type AnnParam struct {
	Literal    *Literal    // Optional literal value.
	VarID      *string     // Optional variable identifier.
	StringLit  *string     // Optional string literal.
	Annotation *Annotation // Optional nested annotation.
}

// BasicExpr is either an identifier or a literal.
type BasicExpr struct {
	Identifier string  // Name of the basic expression.
	Literal    Literal // Literal value of the basic expression.
}

// Literal represents a literal in FlatZinc.
type Literal struct {
	Int      *int         // Optional integer value.
	Bool     *bool        // Optional boolean value.
	Float    *float64     // Optional float value.
	SetInt   *SetIntLit   // Optional set of integers.
	SetFloat *SetFloatLit // Optional set of floats.
}

// SetIntLit is a set of int in FlatZinc.
type SetIntLit struct {
	// Values is set represented as a list of continuous range of integers.
	// For example, set {1, 2, 3, 5} is represented as [[1, 3], [5, 5]].
	Values [][]int
}

// SetFloatLit is a set of float in FlatZinc.
type SetFloatLit struct {
	// Values is a set represented as a list of continuous range of floats.
	// For example float set [[1.0, 2.0], [3.0, 3.0]] contains all the floats
	// between 1.0 and 2.0 (inclusive) and 3.0.
	Values [][]float64
}

// Array represents the index set of a FlatZinc array.
type Array struct {
	Start, End int // Start and end indexes (inclusive) of the array.
}
