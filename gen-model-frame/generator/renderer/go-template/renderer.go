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

type goTemplateModelLayerRenderer struct {
	transform             GoTemplateModelLayerRendererTransformer
	TemplatesForFramePath *moduleTmpl.ModuleTemplates
}

func NewGoTemplateModelLayerRenderer(transform GoTemplateModelLayerRendererTransformer, TemplatesForModules *moduleTmpl.ModuleTemplates) renderer.ModelLayerRenderer {
	return &goTemplateModelLayerRenderer{
		transform:             transform,
		TemplatesForFramePath: TemplatesForModules,
	}
}

func (g *goTemplateModelLayerRenderer) GenerateCodeLayersForFramePath(framePath model.ModelLayers) (map[label.ModelFrameResourceLabel]renderer.RenderedModelLayers, error) {
	out := make(map[label.ModelFrameResourceLabel]renderer.RenderedModelLayers)

	for layerLbl, templatesForLayers := range g.TemplatesForFramePath.Templates {
		hydratedTemplatesForModuleLayers, err := g.hydrateModuleTemplates(framePath, templatesForLayers)
		if err != nil {
			return out, errors.WithStack(err)
		}

		out[layerLbl] = hydratedTemplatesForModuleLayers
	}

	return out, nil
}

func (g *goTemplateModelLayerRenderer) hydrateModuleTemplates(layers model.ModelLayers, templatesForLayers moduleTmpl.TemplatesForLayerImplementations) (map[label.ModelFrameResourceLabel]string, error) {
	templatesForModuleLayers := make(map[label.ModelFrameResourceLabel]string)
	tmplData := g.transform.ModelFramePathToBasicTemplateInputPtr(layers)

	for k, templatesForLayer := range templatesForLayers.LayerImplementationTemplates {
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
func (g *goTemplateModelLayerRenderer) hydrateLayerTemplates(templatesForLayer moduleTmpl.TemplatesForLayerImplementation, tmplData *BasicTemplateInput, layerLabel label.ModelFrameResourceLabel) (map[string]string, error) {
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
func (g *goTemplateModelLayerRenderer) hydrateLayerSectionTemplate(
	tmpl string,
	data *BasicTemplateInput,
	sectionLabel label.ModelFrameResourceLabel,
	layerLabel label.ModelFrameResourceLabel,
) (string, error) {
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
