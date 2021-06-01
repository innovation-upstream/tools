package transform

import (
	"strings"

	"github.com/iancoleman/strcase"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/model"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/model_frame_path"
)

type (
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

func (t *modelFramePathGoTemplateTransformer) LayerSectionsToGoBasicLayoutTemplateInputPtr(sections map[string]string) *model.GoBasicLayoutTemplateInput {
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
		Basic:              basic,
		ModGoPackage:       goPkg,
		ModelGoPackagePath: goPkgPath,
		Sections:           sections,
	}
}
