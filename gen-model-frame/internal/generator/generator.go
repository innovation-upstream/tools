package generator

import (
	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/analyze"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/code_layer_generator"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module"
)

type (
	//go:generate mockgen -destination=../mock/model_frame_generator_mock.go -package=mock gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/generator ModelFrameGenerator
	ModelFrameGenerator interface {
		GenerateFrames() (map[module.ModelFunctionLabel]code_layer_generator.ModuleCodeLayers, error)
	}

	modelFrameGenerator struct {
		Model              model.Model
		ModelAnalyzer      analyze.ModelAnalyzer
		CodeLayerGenerator code_layer_generator.CodeLayerGenerator
	}
)

func NewModelFrameGenerator(m model.Model, modelAnalyzer analyze.ModelAnalyzer, codeLayerGenerator code_layer_generator.CodeLayerGenerator) ModelFrameGenerator {
	return &modelFrameGenerator{
		Model:              m,
		ModelAnalyzer:      modelAnalyzer,
		CodeLayerGenerator: codeLayerGenerator,
	}
}

func (g *modelFrameGenerator) GenerateFrames() (map[module.ModelFunctionLabel]code_layer_generator.ModuleCodeLayers, error) {
	out := make(map[module.ModelFunctionLabel]code_layer_generator.ModuleCodeLayers)

	for _, n := range g.Model.FramePaths {
		res, err := g.CodeLayerGenerator.GenerateCodeLayersForFramePath(n)
		if err != nil {
			return out, errors.WithStack(err)
		}

		out[n.FunctionType] = res
	}

	return out, nil
}
