// Code generated by "stringer -type=VarType"; DO NOT EDIT.

package fzn

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[VarTypeUnknown-0]
	_ = x[VarTypeIntRange-1]
	_ = x[VarTypeIntSet-2]
	_ = x[VarTypeFloatRange-3]
	_ = x[VarTypeBool-4]
}

const _VarType_name = "VarTypeUnknownVarTypeIntRangeVarTypeIntSetVarTypeFloatRangeVarTypeBool"

var _VarType_index = [...]uint8{0, 14, 29, 42, 59, 70}

func (i VarType) String() string {
	if i < 0 || i >= VarType(len(_VarType_index)-1) {
		return "VarType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _VarType_name[_VarType_index[i]:_VarType_index[i+1]]
}
