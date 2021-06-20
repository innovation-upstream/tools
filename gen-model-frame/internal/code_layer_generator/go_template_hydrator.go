package code_layer_generator

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model_frame_path"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/transform"
)

type templateHydrator struct {
	transform             transform.ModelFramePathGoTemplateTransformer
	TemplatesForFramePath *module.ModuleTemplates
}

func NewTemplateHydrator(transform transform.ModelFramePathGoTemplateTransformer, TemplatesForModules *module.ModuleTemplates) CodeLayerGenerator {
	return &templateHydrator{
		transform:             transform,
		TemplatesForFramePath: TemplatesForModules,
	}
}

func (g *templateHydrator) GenerateCodeLayersForFramePath(framePath model_frame_path.ModelFramePath) (ModuleCodeLayers, error) {
	out := make(ModuleCodeLayers)

	hydratedTemplatesForModuleLayers, err := g.hydrateModuleTemplates(framePath, g.TemplatesForFramePath)
	if err != nil {
		return out, errors.WithStack(err)
	}

	fpOutDirName := framePath.Function.GetFileFriendlyName()

	out[fpOutDirName] = hydratedTemplatesForModuleLayers

	return out, nil
}

func (g *templateHydrator) hydrateModuleTemplates(framePath model_frame_path.ModelFramePath, m *module.ModuleTemplates) (map[label.ModelFrameResourceLabel]string, error) {
	templatesForModuleLayers := make(map[label.ModelFrameResourceLabel]string)
	tmplData := g.transform.ModelFramePathToBasicTemplateInputPtr(framePath)
	layerTemplates := m.Templates[framePath.Function]

	for l, templatesForLayer := range layerTemplates.LayerTemplates {
		layerTmplSections, err := g.hydrateLayerTemplates(templatesForLayer, tmplData)
		if err != nil {
			return templatesForModuleLayers, errors.WithStack(err)
		}

		layoutTemplate := layerTemplates.LayoutTemplates[l]
		t := template.Must(
			template.New(fmt.Sprintf("layout_%s", l)).Funcs(sprig.TxtFuncMap()).Parse(layoutTemplate),
		)
		var buff bytes.Buffer
		layoutTmplData := g.transform.LayerSectionsToGoBasicLayoutTemplateInputPtr(*tmplData, layerTmplSections)
		err = t.Execute(&buff, layoutTmplData)
		if err != nil {
			return templatesForModuleLayers, errors.WithStack(err)
		}

		templatesForModuleLayers[l] = buff.String()
	}

	return templatesForModuleLayers, nil
}

// TODO: move this to a layer hydrator struct
func (g *templateHydrator) hydrateLayerTemplates(templatesForLayer module.TemplatesForLayer, tmplData *model.BasicTemplateInput) (map[string]string, error) {
	hydratedLayerSections := make(map[string]string)

	for section, tmpl := range templatesForLayer.SectionTemplates {
		trimmedOut, err := g.hydrateLayerSectionTemplate(tmpl, tmplData, section)
		if err != nil {
			return hydratedLayerSections, errors.WithStack(err)
		}

		hydratedLayerSections[section.GetResourceName()] = trimmedOut
	}

	return hydratedLayerSections, nil
}

// TODO: move this to a section hydrator struct
func (g *templateHydrator) hydrateLayerSectionTemplate(tmpl string, data *model.BasicTemplateInput, sectionLabel label.ModelFrameResourceLabel) (string, error) {
	var hydratedSection string

	t := template.Must(
		template.New(fmt.Sprintf("tmpl_%s", sectionLabel)).Funcs(sprig.TxtFuncMap()).Parse(tmpl),
	)
	var buff bytes.Buffer
	err := t.Execute(&buff, data)
	if err != nil {
		return hydratedSection, errors.WithStack(err)
	}

	hydratedSection = strings.Trim(buff.String(), "\n")

	return hydratedSection, nil
}
