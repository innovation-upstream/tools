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
		LayerSectionsToGoBasicLayoutTemplateInputPtr(layerSections map[string]string) *model.GoBasicLayoutTemplateInput
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
	n := t.model.Name
	return &model.BasicTemplateInput{
		ModCamel:                strcase.ToCamel(n),
		ModLowerCamel:           strcase.ToLowerCamel(n),
		ModSnake:                strcase.ToSnake(n),
		ModKebab:                strcase.ToKebab(n),
		ReferenceTypeCamel:      strcase.ToCamel(string(fp.ReferenceType)),
		ReferenceTypeLowerCamel: strcase.ToLowerCamel(string(fp.ReferenceType)),
		MetaData:                t.model.Metadata,
	}
}

func (t *modelFramePathGoTemplateTransformer) LayerSectionsToGoBasicLayoutTemplateInputPtr(sections map[string]string) *model.GoBasicLayoutTemplateInput {
	n := t.model.Name
	basic := model.BasicLayoutTemplateInput{
		ModCamel:      strcase.ToCamel(n),
		ModLowerCamel: strcase.ToLowerCamel(n),
		ModSnake:      strcase.ToSnake(n),
	}

	return &model.GoBasicLayoutTemplateInput{
		Basic:    basic,
		Sections: sections,
	}
}
