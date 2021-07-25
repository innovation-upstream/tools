package transform

import (
	"github.com/iancoleman/strcase"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model_frame_path"
)

type (
	//go:generate mockgen -destination=../mock/model_frame_path_go_template_transformer_mock.go -package=mock gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/transform ModelFramePathGoTemplateTransformer
	ModelFramePathGoTemplateTransformer interface {
		ModelFramePathToBasicTemplateInputPtr(fp model_frame_path.ModelFramePath) *model.BasicTemplateInput
		LayerSectionsToGoBasicLayoutTemplateInputPtr(basic model.BasicTemplateInput, layerSections map[string]string) *model.GoBasicLayoutTemplateInput
	}
	modelFramePathGoTemplateTransformer struct {
		model *model.Model
	}
)

func NewModelFramePathGoTemplateTransformer(model *model.Model) ModelFramePathGoTemplateTransformer {
	return &modelFramePathGoTemplateTransformer{
		model,
	}
}

func (t *modelFramePathGoTemplateTransformer) ModelFramePathToBasicTemplateInputPtr(fp model_frame_path.ModelFramePath) *model.BasicTemplateInput {
	n := t.model.Label
	return &model.BasicTemplateInput{
		ModCamel:      strcase.ToCamel(n.GetName()),
		ModLowerCamel: strcase.ToLowerCamel(n.GetName()),
		ModSnake:      strcase.ToSnake(n.GetName()),
		ModKebab:      strcase.ToKebab(n.GetName()),
		Options:       t.model.Options,
	}
}

func (t *modelFramePathGoTemplateTransformer) LayerSectionsToGoBasicLayoutTemplateInputPtr(basic model.BasicTemplateInput, sections map[string]string) *model.GoBasicLayoutTemplateInput {
	return &model.GoBasicLayoutTemplateInput{
		Basic:    basic,
		Sections: sections,
	}
}
