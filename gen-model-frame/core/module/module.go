package module

import (
	"github.com/tools/gen-model-frame/core/label"
	"github.com/tools/gen-model-frame/input/config"
)

type (
	ModelFrameModule struct {
		// Name is the fully-qualified name of this module
		Name string `json:"name"`
		// Layers is the list of layers defined by this module
		Layers []ModelLayer `json:"layers"`
	}

	ModelSection struct {
		Label label.ModelFrameResourceLabel `json:"label"`
	}

	ModelFunction struct {
		Label label.ModelFrameResourceLabel `json:"label"`
	}

	ModelLayer struct {
		Label           label.ModelFrameResourceLabel `json:"label"`
		PathTemplate    config.ModelFilePathTemplate  `json:"pathTemplate"`
		Implementations []ModelLayerImplementation    `json:"implementations"`
	}

	ModelLayerFunctionRef struct {
		Label label.ModelFrameResourceLabel `json:"label"`
	}

	ModelLayerImplementationModuleDep struct {
		LayerLabel          label.ModelFrameResourceLabel `json:"layerLabel"`
		ImplementationLabel label.ModelFrameResourceLabel `json:"implementationLabel"`
	}

	ModelLayerFile struct {
		PathTemplate config.ModelFilePathTemplate `json:"pathTemplate"`
	}

	ModelLayerImplementation struct {
		Label    label.ModelFrameResourceLabel       `json:"label"`
		File     ModelLayerFile                      `json:"file"`
		Deps     []ModelLayerImplementationModuleDep `json:"deps"`
		ForLayer label.ModelFrameResourceLabel       `json:"forLayer"`
	}
)

func (m ModelFrameModule) FullyQualifyLabels() ModelFrameModule {
	q := m

	for i, l := range q.Layers {
		ql := label.NewModelFrameResourceLabel(string(l.Label), q.Name, "layer")
		q.Layers[i].Label = ql

		for ii, im := range l.Implementations {
			q.Layers[i].Implementations[ii].Label = label.NewModelFrameResourceLabel(string(im.Label), q.Name, "implementation")
			if q.Layers[i].Implementations[ii].ForLayer == "" {
				q.Layers[i].Implementations[ii].ForLayer = ql
			} else {
				q.Layers[i].Implementations[ii].ForLayer = label.NewModelFrameResourceLabel(string(im.ForLayer), q.Name, "layer")
			}

			for di, d := range im.Deps {
				q.Layers[i].Implementations[ii].Deps[di].LayerLabel = label.NewModelFrameResourceLabel(string(d.LayerLabel), q.Name, "layer")
				if d.ImplementationLabel != "" {
					q.Layers[i].Implementations[ii].Deps[di].ImplementationLabel = label.NewModelFrameResourceLabel(string(d.ImplementationLabel), q.Name, "implementation")
				}
			}
		}
	}

	return q
}

func getLayerByLabel(label label.ModelFrameResourceLabel, layers []ModelLayer) *ModelLayer {
	if len(layers) == 0 {
		return nil
	}

	head := layers[0]
	tail := layers[1:]

	if head.Label == label {
		return &head
	}

	return getLayerByLabel(label, tail)
}

func (m ModelFrameModule) GetLayerByLabel(label label.ModelFrameResourceLabel) *ModelLayer {
	return getLayerByLabel(label, m.Layers)
}

func findImplementationInLayer(label label.ModelFrameResourceLabel, impls []ModelLayerImplementation) *ModelLayerImplementation {
	if len(impls) == 0 {
		return nil
	}

	head := impls[0]
	tail := impls[1:]

	if head.Label == label {
		return &head
	}

	return findImplementationInLayer(label, tail)
}

func getLayerImplementationByLabel(label label.ModelFrameResourceLabel, layers []ModelLayer) *ModelLayerImplementation {
	if len(layers) == 0 {
		return nil
	}

	head := layers[0]
	tail := layers[1:]

	match := findImplementationInLayer(label, head.Implementations)

	if match != nil {
		return match
	}

	return getLayerImplementationByLabel(label, tail)
}

func (m ModelFrameModule) GetLayerImplementationByLabel(label label.ModelFrameResourceLabel) *ModelLayerImplementation {
	return getLayerImplementationByLabel(label, m.Layers)
}
