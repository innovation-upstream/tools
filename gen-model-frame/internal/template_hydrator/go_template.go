package template_hydrator

import (
	"bytes"
	"html/template"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/model_frame_path"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/module"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/transform"
)

type templateHydrator struct {
	transform           transform.ModelFramePathGoTemplateTransformer
	TemplatesForModules []*module.ModuleTemplates
}

func NewTemplateHydrator(transform transform.ModelFramePathGoTemplateTransformer, TemplatesForModules []*module.ModuleTemplates) Generator {
	return &templateHydrator{
		transform:           transform,
		TemplatesForModules: TemplatesForModules,
	}
}

func (g *templateHydrator) Generate(framePath model_frame_path.ModelFramePath) (map[string]map[module.ModelFrameLayerLabel]string, error) {
	out := make(map[string]map[module.ModelFrameLayerLabel]string)

	for _, m := range g.TemplatesForModules {
		templatesForModuleLayers := make(map[module.ModelFrameLayerLabel]string)
		templatesForFunction := m.GetTemplatesForFunctionLabel(framePath.FunctionType)
		tmplData := g.transform.ModelFramePathToBasicTemplateInputPtr(framePath)

		for _, l := range framePath.Layers {
			layerName := module.ModelFrameLayerLabel(strings.Split(string(l), "::")[1])
			layerTmplSections := make(map[string]string)
			templatesForLayer := templatesForFunction.LayerTemplates[layerName]
			layoutTemplate := templatesForFunction.LayoutTemplates[layerName]

			for section, tmpl := range templatesForLayer.SectionTemplates {
				t := template.Must(template.New("mod_fp_tmpl").Parse(tmpl))
				var buff bytes.Buffer
				err := t.Execute(&buff, tmplData)
				if err != nil {
					return out, errors.WithStack(err)
				}

				trimmedOut := strings.Trim(buff.String(), "\n")
				layerTmplSections[section] = trimmedOut
			}
			t := template.Must(template.New("mod_fp_layout_tmpl").Parse(layoutTemplate))
			var buff bytes.Buffer
			layoutTmplData := g.transform.LayerSectionsToGoBasicLayoutTemplateInputPtr(layerTmplSections)
			err := t.Execute(&buff, layoutTmplData)
			if err != nil {
				return out, errors.WithStack(err)
			}

			templatesForModuleLayers[layerName] = buff.String()
		}

		moduleName := m.Module.Name
		var reAt = regexp.MustCompile(`^@`)
		moduleName = reAt.ReplaceAllString(moduleName, "")
		var reSlash = regexp.MustCompile(`\/`)
		moduleName = reSlash.ReplaceAllString(moduleName, "_")
		var reDash = regexp.MustCompile(`-`)
		moduleName = reDash.ReplaceAllString(moduleName, "_")

		out[moduleName] = templatesForModuleLayers
	}

	return out, nil
}
