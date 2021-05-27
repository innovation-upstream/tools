package template

import (
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/model_frame_path"
	repoTemplate "gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/template/frame/data/repo"
)

type (
	TemplateTargetLanguage string
	TemplateSection        string
)

const TemplateTargetLanguageGolang = TemplateTargetLanguage("golang")

const TemplateSectionInterface = TemplateSection("interface")
const TemplateSectionMethod = TemplateSection("method")

type TemplateRegistry interface {
	GetTemplatesForModelFrameLayer(layer model_frame_path.ModelFrameLayerType) (templatesForLayer, string)
}

type templatesForLayer struct {
	SectionTemplates map[TemplateSection]string
}

type templatesForFunctionType struct {
	LayerTemplates map[model_frame_path.ModelFrameLayerType]templatesForLayer
}

type golangTemplateRegistry struct {
	FunctionTemplates map[model_frame_path.ModelFunctionType]templatesForFunctionType
	LayoutTemplates   map[model_frame_path.ModelFrameLayerType]string
	ModelFramePath    model_frame_path.ModelFramePath
}

func NewGolangTemplateRegistry(modelFramePath model_frame_path.ModelFramePath) TemplateRegistry {
	dataRepoCreateSectionTmpl := templatesForLayer{
		SectionTemplates: map[TemplateSection]string{
			TemplateSectionInterface: repoTemplate.InterfaceRepoCreate,
			TemplateSectionMethod:    repoTemplate.MethodRepoCreate,
		},
	}

	createLayerTmpl := templatesForFunctionType{
		LayerTemplates: map[model_frame_path.ModelFrameLayerType]templatesForLayer{
			model_frame_path.ModelFrameLayerTypeDataRepo: dataRepoCreateSectionTmpl,
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
		FunctionTemplates: map[model_frame_path.ModelFunctionType]templatesForFunctionType{
			model_frame_path.ModelFunctionTypeCreate: createLayerTmpl,
		},
		LayoutTemplates: map[model_frame_path.ModelFrameLayerType]string{
			model_frame_path.ModelFrameLayerTypeDataRepo: repoTemplate.RepoLayout,
		},
		ModelFramePath: modelFramePath,
	}
}

func (r *golangTemplateRegistry) GetTemplatesForModelFrameLayer(layer model_frame_path.ModelFrameLayerType) (templatesForLayer, string) {
	return r.FunctionTemplates[r.ModelFramePath.FunctionType].LayerTemplates[layer], r.LayoutTemplates[layer]
}
