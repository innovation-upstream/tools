package generator

import (
	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/model"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/model_frame_path"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/template_hydrator"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/transform"
)

type (
	GolangModelGenerator interface {
		Generate(m model.Model) (map[model_frame_path.ModelFramePathType]map[model_frame_path.ModelFrameLayerType]string, error)
	}
	golangModelGenerator struct{}
)

func NewGolangModelGenerator() GolangModelGenerator {
	return &golangModelGenerator{}
}

func (g *golangModelGenerator) Generate(m model.Model) (map[model_frame_path.ModelFramePathType]map[model_frame_path.ModelFrameLayerType]string, error) {
	out := make(map[model_frame_path.ModelFramePathType]map[model_frame_path.ModelFrameLayerType]string)
	tr := transform.NewModelFramePathGoTemplateTransformer(&m)
	gen := template_hydrator.NewGoTemplateHydrator(tr)
	for _, n := range m.FramePaths {
		res, err := gen.Generate(n)
		if err != nil {
			return out, errors.WithStack(err)
		}

		out[n.Type] = res
	}

	return out, nil
}
