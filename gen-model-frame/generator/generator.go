package generator

import (
	"github.com/pkg/errors"
	"innovationup.stream/tools/gen-model-frame/core/label"
	"innovationup.stream/tools/gen-model-frame/core/model"
	"innovationup.stream/tools/gen-model-frame/generator/renderer"
)

type (
	ModelFrameGenerator interface {
		GenerateFrames(fp model.ModelLayers) (map[label.ModelFrameResourceLabel]renderer.RenderedModelLayers, error)
	}

	modelFrameGenerator struct {
		CodeLayerGenerator renderer.ModelLayerRenderer
	}
)

func NewModelFrameGenerator(codeLayerGenerator renderer.ModelLayerRenderer) ModelFrameGenerator {
	return &modelFrameGenerator{
		CodeLayerGenerator: codeLayerGenerator,
	}
}

func (g *modelFrameGenerator) GenerateFrames(fp model.ModelLayers) (map[label.ModelFrameResourceLabel]renderer.RenderedModelLayers, error) {
	out := make(map[label.ModelFrameResourceLabel]renderer.RenderedModelLayers)

	// TODO: add feature to allow models.json to specify custom generator binary
	codeLayers, err := g.CodeLayerGenerator.GenerateCodeLayersForFramePath(fp)
	if err != nil {
		return out, errors.WithStack(err)
	}

	// TODO: this is where we would call plugins for post-processing

	out = codeLayers

	return out, nil
}
