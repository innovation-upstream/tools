package model_function

// Represent a collection of model frames and how they relate/interact with
// each other
// Can be used to generate a chain of related functions (io, logic, data)
type ModelFramePath struct {
	FunctionType ModelFunctionType
	// What reference is used in this frame path
	ReferenceType ReferenceType
	DataFrameType DataFrameType
}
