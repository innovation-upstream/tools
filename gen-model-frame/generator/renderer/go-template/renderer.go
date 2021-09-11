package renderer

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/pkg/errors"
	"innovationup.stream/tools/gen-model-frame/core/label"
	"innovationup.stream/tools/gen-model-frame/core/model"
	moduleTmpl "innovationup.stream/tools/gen-model-frame/core/template"
	"innovationup.stream/tools/gen-model-frame/generator/renderer"
)

type templateHydrator struct {
	transform             ModelFramePathGoTemplateTransformer
	TemplatesForFramePath *moduleTmpl.ModuleTemplates
}

func NewTemplateHydrator(transform ModelFramePathGoTemplateTransformer, TemplatesForModules *moduleTmpl.ModuleTemplates) renderer.CodeLayerRenderer {
	return &templateHydrator{
		transform:             transform,
		TemplatesForFramePath: TemplatesForModules,
	}
}

func (g *templateHydrator) GenerateCodeLayersForFramePath(framePath model.ModelLayers) ([]renderer.RenderedCodeLayers, error) {
	var out []renderer.RenderedCodeLayers

	for _, templatesForLayers := range g.TemplatesForFramePath.Templates {
		hydratedTemplatesForModuleLayers, err := g.hydrateModuleTemplates(framePath, templatesForLayers)
		if err != nil {
			return out, errors.WithStack(err)
		}

		out = append(out, hydratedTemplatesForModuleLayers)
	}

	return out, nil
}

func (g *templateHydrator) hydrateModuleTemplates(framePath model.ModelLayers, templatesForLayers moduleTmpl.TemplatesForLayers) (map[label.ModelFrameResourceLabel]string, error) {
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
func (g *templateHydrator) hydrateLayerTemplates(templatesForLayer moduleTmpl.TemplatesForLayer, tmplData *BasicTemplateInput, layerLabel label.ModelFrameResourceLabel) (map[string]string, error) {
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
func (g *templateHydrator) hydrateLayerSectionTemplate(tmpl string, data *BasicTemplateInput, sectionLabel label.ModelFrameResourceLabel, layerLabel label.ModelFrameResourceLabel) (string, error) {
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
