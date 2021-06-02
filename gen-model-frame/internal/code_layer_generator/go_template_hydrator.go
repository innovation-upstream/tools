package code_layer_generator

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model_frame_path"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/transform"
)

type templateHydrator struct {
	transform           transform.ModelFramePathGoTemplateTransformer
	TemplatesForModules []*module.ModuleTemplates
}

func NewTemplateHydrator(transform transform.ModelFramePathGoTemplateTransformer, TemplatesForModules []*module.ModuleTemplates) CodeLayerGenerator {
	return &templateHydrator{
		transform:           transform,
		TemplatesForModules: TemplatesForModules,
	}
}

func (g *templateHydrator) GenerateCodeLayersForFramePath(framePath model_frame_path.ModelFramePath) (ModuleCodeLayers, error) {
	out := make(ModuleCodeLayers)

	for _, m := range g.TemplatesForModules {
		templatesForModuleLayers, err := g.hydrateModuleTemplates(framePath, m)
		if err != nil {
			return out, errors.WithStack(err)
		}

		moduleName := m.Module.Name.GetFileFriendlyName()

		out[moduleName] = templatesForModuleLayers
	}

	return out, nil
}

func (g *templateHydrator) hydrateModuleTemplates(framePath model_frame_path.ModelFramePath, m *module.ModuleTemplates) (map[module.ModelFrameLayerLabel]string, error) {
	templatesForModuleLayers := make(map[module.ModelFrameLayerLabel]string)
	templatesForFunction := m.GetTemplatesForFunctionLabel(framePath.FunctionType)
	tmplData := g.transform.ModelFramePathToBasicTemplateInputPtr(framePath)

	for _, l := range framePath.Layers {
		layerName := module.ModelFrameLayerLabel(strings.Split(string(l), "::")[1])
		templatesForLayer := templatesForFunction.LayerTemplates[layerName]

		layerTmplSections, err := g.hydrateLayerTemplates(templatesForLayer, tmplData)
		if err != nil {
			return templatesForModuleLayers, errors.WithStack(err)
		}

		layoutTemplate := templatesForFunction.LayoutTemplates[layerName]
		t := template.Must(template.New("mod_fp_layout_tmpl").Parse(layoutTemplate))
		var buff bytes.Buffer
		layoutTmplData := g.transform.LayerSectionsToGoBasicLayoutTemplateInputPtr(layerTmplSections)
		err = t.Execute(&buff, layoutTmplData)
		if err != nil {
			return templatesForModuleLayers, errors.WithStack(err)
		}

		templatesForModuleLayers[layerName] = buff.String()
	}

	return templatesForModuleLayers, nil
}

func (g *templateHydrator) hydrateLayerTemplates(templatesForLayer module.TemplatesForLayer, tmplData *model.BasicTemplateInput) (map[string]string, error) {
	hydratedLayerSections := make(map[string]string)

	for section, tmpl := range templatesForLayer.SectionTemplates {
		trimmedOut, err := g.hydrateLayerSectionTemplate(tmpl, tmplData)
		if err != nil {
			return hydratedLayerSections, errors.WithStack(err)
		}

		hydratedLayerSections[section] = trimmedOut
	}

	return hydratedLayerSections, nil
}

func (g *templateHydrator) hydrateLayerSectionTemplate(tmpl string, data *model.BasicTemplateInput) (string, error) {
	var hydratedSection string

	t := template.Must(template.New("layer_section_tmpl").Parse(tmpl))
	var buff bytes.Buffer
	err := t.Execute(&buff, data)
	if err != nil {
		return hydratedSection, errors.WithStack(err)
	}

	hydratedSection = strings.Trim(buff.String(), "\n")

	return hydratedSection, nil
}
