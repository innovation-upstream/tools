package model_frame_path

// Represent a collection of model frames and how they relate/interact with
// each other
// Can be used to generate a chain of related functions (io, logic, data)
type ModelFramePath struct {
	FunctionType ModelFunctionType `json:"function_type"`
	// What reference is used in this frame path
	ReferenceType ReferenceType `json:"reference_type"`
	DataFrameType DataFrameType `json:"data_frame_type"`
}
