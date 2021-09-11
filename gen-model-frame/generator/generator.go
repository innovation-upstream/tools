package generator

import (
	"github.com/pkg/errors"
	"innovationup.stream/tools/gen-model-frame/core/model"
	"innovationup.stream/tools/gen-model-frame/generator/renderer"
)

type (
	ModelFrameGenerator interface {
		GenerateFrames(fp model.ModelLayers) ([]renderer.RenderedCodeLayers, error)
	}

	modelFrameGenerator struct {
		CodeLayerGenerator renderer.CodeLayerRenderer
	}
)

func NewModelFrameGenerator(codeLayerGenerator renderer.CodeLayerRenderer) ModelFrameGenerator {
	return &modelFrameGenerator{
		CodeLayerGenerator: codeLayerGenerator,
	}
}

func (g *modelFrameGenerator) GenerateFrames(fp model.ModelLayers) ([]renderer.RenderedCodeLayers, error) {
	var out []renderer.RenderedCodeLayers

	// TODO: add feature to allow models.json to specify custom generator binary
	codeLayers, err := g.CodeLayerGenerator.GenerateCodeLayersForFramePath(fp)
	if err != nil {
		return out, errors.WithStack(err)
	}

	// TODO: this is where we would call plugins for post-processing

	out = codeLayers

	return out, nil
}
