package transform

import (
	"strings"

	"github.com/iancoleman/strcase"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/model"
	modelTemplate "gitlab.innovationup.stream/innovation-upstream/gen-model-frame/model/generator/template"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/model/model_function"
)

type (
	ModelFramePathGoTemplateTransformer interface {
		ModelFramePathToBasicTemplateInputPtr(fp model_function.ModelFramePath) *model.BasicTemplateInput
		LayerSectionsToGoBasicLayoutTemplateInputPtr(layerSections map[modelTemplate.TemplateSection]string) *model.GoBasicLayoutTemplateInput
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

func (t *modelFramePathGoTemplateTransformer) ModelFramePathToBasicTemplateInputPtr(fp model_function.ModelFramePath) *model.BasicTemplateInput {
	n := t.model.Name
	updatableFieldsStr := t.model.Metadata[model.ModelMetadataUpdatableFields]
	updatableFields := strings.Split(updatableFieldsStr, ",")
	return &model.BasicTemplateInput{
		ModCamel:                strcase.ToCamel(n),
		ModLowerCamel:           strcase.ToLowerCamel(n),
		ModSnake:                strcase.ToSnake(n),
		ModKebab:                strcase.ToKebab(n),
		ReferenceTypeCamel:      strcase.ToCamel(string(fp.ReferenceType)),
		ReferenceTypeLowerCamel: strcase.ToLowerCamel(string(fp.ReferenceType)),
		UpdateableFields:        updatableFields,
	}
}

func (t *modelFramePathGoTemplateTransformer) LayerSectionsToGoBasicLayoutTemplateInputPtr(layerSections map[modelTemplate.TemplateSection]string) *model.GoBasicLayoutTemplateInput {
	n := t.model.Name
	basic := model.BasicLayoutTemplateInput{
		ModCamel:      strcase.ToCamel(n),
		ModLowerCamel: strcase.ToLowerCamel(n),
		ModSnake:      strcase.ToSnake(n),
	}

	// TODO: fail if either are empty
	goPkg := t.model.Metadata[model.ModelMetadataGolangModelPackage]
	goPkgPath := t.model.Metadata[model.ModelMetadataGolangModelPackagePath]

	return &model.GoBasicLayoutTemplateInput{
		Basic:               basic,
		ModGoPackage:        goPkg,
		ModelGoPackagePath:  goPkgPath,
		InterfaceDefinition: layerSections[modelTemplate.TemplateSectionInterface],
		Methods:             layerSections[modelTemplate.TemplateSectionMethod],
	}
}
