package model

import (
	"innovationup.stream/tools/gen-model-frame/core/module"
)

// Represent a collection of model frames and how they relate/interact with
// each other
// Can be used to generate a chain of related functions (io, logic, data)
type ModelLayers struct {
	Layers []module.ModelLayerImplementationModuleDep `json:"layers"`
}
