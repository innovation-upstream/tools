package generator

import (
	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/code_layer_generator"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model_frame_path"
)

type (
	//go:generate mockgen -destination=../mock/model_frame_generator_mock.go -package=mock gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/generator ModelFrameGenerator
	ModelFrameGenerator interface {
		GenerateFrames(fp model_frame_path.ModelFramePath) ([]code_layer_generator.CodeLayers, error)
	}

	modelFrameGenerator struct {
		CodeLayerGenerator code_layer_generator.CodeLayerGenerator
	}
)

func NewModelFrameGenerator(codeLayerGenerator code_layer_generator.CodeLayerGenerator) ModelFrameGenerator {
	return &modelFrameGenerator{
		CodeLayerGenerator: codeLayerGenerator,
	}
}

func (g *modelFrameGenerator) GenerateFrames(fp model_frame_path.ModelFramePath) ([]code_layer_generator.CodeLayers, error) {
	var out []code_layer_generator.CodeLayers

	// TODO: add feature to allow models.json to specify custom generator binary
	codeLayers, err := g.CodeLayerGenerator.GenerateCodeLayersForFramePath(fp)
	if err != nil {
		return out, errors.WithStack(err)
	}

	// TODO: this is where we would call plugins for post-processing

	out = codeLayers

	return out, nil
}
