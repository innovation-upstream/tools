package renderer

import (
	"github.com/iancoleman/strcase"
	"innovationup.stream/tools/gen-model-frame/core/model"
)

type (
	GoTemplateModelLayerRendererTransformer interface {
		ModelFramePathToBasicTemplateInputPtr(
			fp model.ModelLayers,
		) *BasicTemplateInput
		LayerSectionsToGoBasicLayoutTemplateInputPtr(
			basic BasicTemplateInput,
			layerSections map[string]string,
		) *GoBasicLayoutTemplateInput
	}

	goTemplateModelLayerRendererTransformer struct {
		model *model.Model
	}
)

func NewModelFramePathGoTemplateTransformer(
	model *model.Model,
) GoTemplateModelLayerRendererTransformer {
	return &goTemplateModelLayerRendererTransformer{
		model,
	}
}

func (t *goTemplateModelLayerRendererTransformer) ModelFramePathToBasicTemplateInputPtr(
	fp model.ModelLayers,
) *BasicTemplateInput {
	n := t.model.Label
	return &BasicTemplateInput{
		ModCamel:      strcase.ToCamel(n.GetName()),
		ModLowerCamel: strcase.ToLowerCamel(n.GetName()),
		ModSnake:      strcase.ToSnake(n.GetName()),
		ModKebab:      strcase.ToKebab(n.GetName()),
		Options:       t.model.Options,
		Hooks:         t.model.Hooks,
	}
}

func (t *goTemplateModelLayerRendererTransformer) LayerSectionsToGoBasicLayoutTemplateInputPtr(
	basic BasicTemplateInput,
	sections map[string]string,
) *GoBasicLayoutTemplateInput {
	return &GoBasicLayoutTemplateInput{
		Basic:    basic,
		Sections: sections,
	}
}
