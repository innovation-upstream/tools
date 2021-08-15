package module

import (
	"fmt"

	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/config"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/regexp"
)

type (
	ModelFrameModule struct {
		// Name is the fully-qualified name of this module
		Name label.ModelFrameResourceLabel `json:"name"`
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
		Label label.ModelFrameResourceLabel `json:"label"`
		Deps  []ModelLayerModuleDep         `json:"deps"`
		File  ModelLayerFile                `json:"file"`
	}

	ModelLayerFunctionRef struct {
		Label label.ModelFrameResourceLabel `json:"label"`
	}

	ModelLayerModuleDep struct {
		Label label.ModelFrameResourceLabel `json:"label"`
	}

	ModelLayerFile struct {
		PathTemplate config.ModelFilePathTemplate `json:"pathTemplate"`
	}
)

func (m ModelFrameModule) GetTransitiveModuleLabels() []label.ModelFrameResourceLabel {
	var trLabels []label.ModelFrameResourceLabel

	for _, l := range m.Layers {
		for _, trm := range l.Deps {
			trLabels = append(trLabels, trm.Label)
		}
	}

	return trLabels
}

func (m ModelFrameModule) FullyQualifyLabels() ModelFrameModule {
	q := m

	for i, l := range q.Layers {
		layerLabelIsQualified := regexp.ModelFrameResourceLabelPattern.MatchString(string(l.Label))
		if layerLabelIsQualified == false {
			fql := label.ModelFrameResourceLabel(fmt.Sprintf("%s::%s/%s", q.Name, "layer", l.Label))
			q.Layers[i].Label = fql
		}

		for di, d := range l.Deps {
			secLabelIsQualified := regexp.ModelFrameResourceLabelPattern.MatchString(string(d.Label))
			if secLabelIsQualified == false {
				dfql := label.ModelFrameResourceLabel(fmt.Sprintf("%s::%s/%s", q.Name, "layer", d.Label))
				q.Layers[i].Deps[di].Label = dfql
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
