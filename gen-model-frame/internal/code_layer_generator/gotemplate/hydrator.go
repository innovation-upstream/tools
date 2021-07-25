package gotemplate

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/code_layer_generator"
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

func NewTemplateHydrator(transform transform.ModelFramePathGoTemplateTransformer, TemplatesForModules *module.ModuleTemplates) code_layer_generator.CodeLayerGenerator {
	return &templateHydrator{
		transform:             transform,
		TemplatesForFramePath: TemplatesForModules,
	}
}

func (g *templateHydrator) GenerateCodeLayersForFramePath(framePath model_frame_path.ModelFramePath) ([]code_layer_generator.CodeLayers, error) {
	var out []code_layer_generator.CodeLayers

	for _, templatesForLayers := range g.TemplatesForFramePath.Templates {
		hydratedTemplatesForModuleLayers, err := g.hydrateModuleTemplates(framePath, templatesForLayers)
		if err != nil {
			return out, errors.WithStack(err)
		}

		out = append(out, hydratedTemplatesForModuleLayers)
	}

	return out, nil
}

func (g *templateHydrator) hydrateModuleTemplates(framePath model_frame_path.ModelFramePath, templatesForLayers module.TemplatesForLayers) (map[label.ModelFrameResourceLabel]string, error) {
	templatesForModuleLayers := make(map[label.ModelFrameResourceLabel]string)
	tmplData := g.transform.ModelFramePathToBasicTemplateInputPtr(framePath)

	for k, templatesForLayer := range templatesForLayers.LayerTemplates {
		layerTmplSections, err := g.hydrateLayerTemplates(templatesForLayer, tmplData, k)
		if err != nil {
			return templatesForModuleLayers, errors.WithStack(err)
		}

		layoutTemplate := templatesForLayer.LayoutTemplate
		t := template.Must(
			template.New(fmt.Sprintf("layout_%s", k.GetFileFriendlyName())).
				Funcs(sprig.TxtFuncMap()).
				Funcs(TxtFuncMap()).
				Parse(layoutTemplate),
		)

		layoutTmplData := g.transform.LayerSectionsToGoBasicLayoutTemplateInputPtr(*tmplData, layerTmplSections)
		var buff bytes.Buffer
		err = t.Execute(&buff, layoutTmplData)
		if err != nil {
			return templatesForModuleLayers, errors.WithStack(err)
		}

		templatesForModuleLayers[k] = buff.String()
	}

	return templatesForModuleLayers, nil
}

// TODO: move this to a layer hydrator struct
func (g *templateHydrator) hydrateLayerTemplates(templatesForLayer module.TemplatesForLayer, tmplData *model.BasicTemplateInput, layerLabel label.ModelFrameResourceLabel) (map[string]string, error) {
	hydratedLayerSections := make(map[string]string)

	for section, tmpl := range templatesForLayer.SectionTemplates {
		trimmedOut, err := g.hydrateLayerSectionTemplate(tmpl, tmplData, section, layerLabel)
		if err != nil {
			return hydratedLayerSections, errors.WithStack(err)
		}

		hydratedLayerSections[section.GetResourceName()] = trimmedOut
	}

	return hydratedLayerSections, nil
}

// TODO: move this to a section hydrator struct
func (g *templateHydrator) hydrateLayerSectionTemplate(tmpl string, data *model.BasicTemplateInput, sectionLabel label.ModelFrameResourceLabel, layerLabel label.ModelFrameResourceLabel) (string, error) {
	var hydratedSection string

	t := template.Must(
		template.New(fmt.Sprintf("tmpl_%s_%s", layerLabel, sectionLabel)).
			Funcs(sprig.TxtFuncMap()).
			Funcs(TxtFuncMap()).
			Parse(tmpl),
	)
	var buff bytes.Buffer
	err := t.Execute(&buff, data)
	if err != nil {
		return hydratedSection, errors.WithStack(err)
	}

	hydratedSection = strings.Trim(buff.String(), "\n")

	return hydratedSection, nil
}
