package template

import (
	repoTemplate "gitlab.innovationup.stream/innovation-upstream/gen-model-frame/model/generator/template/frame/data/repo"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/model/model_function"
)

type (
	TemplateTargetLanguage string
	TemplateSection        string
)

const TemplateTargetLanguageGolang = TemplateTargetLanguage("golang")

const TemplateSectionInterface = TemplateSection("interface")
const TemplateSectionMethod = TemplateSection("method")

type TemplateRegistry interface {
	GetTemplatesForModelFrameLayer(layer model_function.ModelFrameLayerType) (templatesForLayer, string)
}

type templatesForLayer struct {
	SectionTemplates map[TemplateSection]string
}

type templatesForFunctionType struct {
	LayerTemplates map[model_function.ModelFrameLayerType]templatesForLayer
}

type golangTemplateRegistry struct {
	FunctionTemplates map[model_function.ModelFunctionType]templatesForFunctionType
	LayoutTemplates   map[model_function.ModelFrameLayerType]string
	ModelFramePath    model_function.ModelFramePath
}

func NewGolangTemplateRegistry(modelFramePath model_function.ModelFramePath) TemplateRegistry {
	dataRepoCreateSectionTmpl := templatesForLayer{
		SectionTemplates: map[TemplateSection]string{
			TemplateSectionInterface: repoTemplate.InterfaceRepoCreate,
			TemplateSectionMethod:    repoTemplate.MethodRepoCreate,
		},
	}

	createLayerTmpl := templatesForFunctionType{
		LayerTemplates: map[model_function.ModelFrameLayerType]templatesForLayer{
			model_function.ModelFrameLayerTypeDataRepo: dataRepoCreateSectionTmpl,
		},
	}

	// TODO: write the rest of the layer section templates + their template
	// structs here

	/*
		dataRepoUpdateSectionTmpl := templatesForLayer{
			SectionTemplates: map[TemplateSection]string{
				TemplateSectionInterface: repoTemplate.InterfaceRepoUpdate,
				TemplateSectionMethod:    repoTemplate.MethodRepoUpdate,
			},
		}

		dataRepoDeleteSectionTmpl := templatesForLayer{
			SectionTemplates: map[TemplateSection]string{
				TemplateSectionInterface: repoTemplate.InterfaceRepoDelete,
				TemplateSectionMethod:    repoTemplate.MethodRepoDelete,
			},
		}

		dataRepoReadAllSectionTmpl := templatesForLayer{
			SectionTemplates: map[TemplateSection]string{
				TemplateSectionInterface: repoTemplate.InterfaceRepoReadAll,
				TemplateSectionMethod:    repoTemplate.MethodRepoReadAll,
			},
		}
	*/

	return &golangTemplateRegistry{
		FunctionTemplates: map[model_function.ModelFunctionType]templatesForFunctionType{
			model_function.ModelFunctionTypeCreate: createLayerTmpl,
		},
		LayoutTemplates: map[model_function.ModelFrameLayerType]string{
			model_function.ModelFrameLayerTypeDataClient: repoTemplate.RepoLayout,
		},
		ModelFramePath: modelFramePath,
	}
}

func (r *golangTemplateRegistry) GetTemplatesForModelFrameLayer(layer model_function.ModelFrameLayerType) (templatesForLayer, string) {
	return r.FunctionTemplates[r.ModelFramePath.FunctionType].LayerTemplates[layer], r.LayoutTemplates[layer]
}
