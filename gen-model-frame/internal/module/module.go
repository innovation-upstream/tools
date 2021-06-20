package module

import (
	"fmt"

	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"
)

type (
	ModelFrameModule struct {
		// Name is the fully-qualified name of this module
		Name label.ModelFrameResourceLabel `json:"name"`
		// Functions is the list of functions defined by this module
		Functions []ModelFunction `json:"functions"`
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
		// Functions is the list of functions this layer can implement
		Functions []ModelLayerFunctionRef `json:"functions"`
		Sections  []ModelSection          `json:"sections"`
		Deps      []ModelLayerModuleDep   `json:"deps"`
	}

	ModelLayerFunctionRef struct {
		Label label.ModelFrameResourceLabel `json:"label"`
	}

	ModelLayerModuleDep struct {
		Label label.ModelFrameResourceLabel `json:"label"`
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
		layerLabelIsQualified := label.ModelFrameResourceLabelPattern.MatchString(string(l.Label))
		if layerLabelIsQualified == false {
			fql := label.ModelFrameResourceLabel(fmt.Sprintf("%s::%s/%s", q.Name, "layer", l.Label))
			q.Layers[i].Label = fql
		}

		for si, s := range l.Sections {
			secLabelIsQualified := label.ModelFrameResourceLabelPattern.MatchString(string(s.Label))
			if secLabelIsQualified == false {
				sfql := label.ModelFrameResourceLabel(fmt.Sprintf("%s::%s/%s", q.Name, "section", s.Label))
				q.Layers[i].Sections[si].Label = sfql
			}
		}

		for di, d := range l.Deps {
			secLabelIsQualified := label.ModelFrameResourceLabelPattern.MatchString(string(d.Label))
			if secLabelIsQualified == false {
				dfql := label.ModelFrameResourceLabel(fmt.Sprintf("%s::%s/%s", q.Name, "layer", d.Label))
				q.Layers[i].Deps[di].Label = dfql
			}
		}

		for fi, l := range l.Functions {
			secLabelIsQualified := label.ModelFrameResourceLabelPattern.MatchString(string(l.Label))
			if secLabelIsQualified == false {
				ffql := label.ModelFrameResourceLabel(fmt.Sprintf("%s::%s/%s", q.Name, "function", l.Label))
				q.Layers[i].Functions[fi].Label = ffql
			}
		}
	}

	for fi, f := range q.Functions {
		functionLabelIsQualified := label.ModelFrameResourceLabelPattern.MatchString(string(f.Label))
		if functionLabelIsQualified == false {
			ffql := label.ModelFrameResourceLabel(fmt.Sprintf("%s::%s/%s", q.Name, "function", f.Label))
			q.Functions[fi].Label = ffql
		}
	}

	return q
}
