package model_gen

import (
	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/model"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/model/generator"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/model/model_function"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/model/transform"
)

type (
	GolangModelGenerator interface {
		Generate(m model.Model) (map[model_function.ModelFunctionType]map[model_function.ModelFrameLayerType]string, error)
	}
	golangModelGenerator struct{}
)

func NewGolangModelGenerator() GolangModelGenerator {
	return &golangModelGenerator{}
}

func (g *golangModelGenerator) Generate(m model.Model) (map[model_function.ModelFunctionType]map[model_function.ModelFrameLayerType]string, error) {
	var out map[model_function.ModelFunctionType]map[model_function.ModelFrameLayerType]string
	tr := transform.NewModelFramePathGoTemplateTransformer(&m)
	gen := generator.NewGoTemplateGenerator(tr)
	for _, n := range m.FramePaths {
		res, err := gen.Generate(n)
		if err != nil {
			return out, errors.WithStack(err)
		}

		out[n.FunctionType] = res
	}

	return out, nil
}
