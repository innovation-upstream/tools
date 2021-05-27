package template_hydrator

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/model_frame_path"
	modelTemplate "gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/template"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/transform"
)

type goTemplateHydrator struct {
	transform transform.ModelFramePathGoTemplateTransformer
}

func NewGoTemplateHydrator(transform transform.ModelFramePathGoTemplateTransformer) Generator {
	return &goTemplateHydrator{
		transform,
	}
}

func (g *goTemplateHydrator) Generate(mod model_frame_path.ModelFramePath) (map[model_frame_path.ModelFrameLayerType]string, error) {
	out := make(map[model_frame_path.ModelFrameLayerType]string)
	out, err := g.generateDataRepo(mod)
	if err != nil {
		return out, errors.WithStack(err)
	}

	return out, nil
}

func (g *goTemplateHydrator) generateDataRepo(framePath model_frame_path.ModelFramePath) (map[model_frame_path.ModelFrameLayerType]string, error) {
	out := make(map[model_frame_path.ModelFrameLayerType]string)
	layers := []model_frame_path.ModelFrameLayerType{
		model_frame_path.ModelFrameLayerTypeIO,
	}
	//TODO: implement the rest of the dataframe types + handle different model frame types in template inline logic
	switch framePath.DataFrameType {
	case model_frame_path.DataFrameTypeRepo:
		layers = append(layers, []model_frame_path.ModelFrameLayerType{
			model_frame_path.ModelFrameLayerTypeDataRepo,
			model_frame_path.ModelFrameLayerTypeLogicRepo,
		}...)
		break
	case model_frame_path.DataFrameTypeRelay:
		layers = append(layers, []model_frame_path.ModelFrameLayerType{
			model_frame_path.ModelFrameLayerTypeDataRelay,
			model_frame_path.ModelFrameLayerTypeLogicRelay,
		}...)
		break
	case model_frame_path.DataFrameTypeClient:
		layers = append(layers, []model_frame_path.ModelFrameLayerType{
			model_frame_path.ModelFrameLayerTypeDataClient,
			model_frame_path.ModelFrameLayerTypeLogicClient,
		}...)
		break
	}

	tmplReg := modelTemplate.NewGolangTemplateRegistry(framePath)
	tmplData := g.transform.ModelFramePathToBasicTemplateInputPtr(framePath)
	for _, l := range layers {
		layerTmplSections := make(map[modelTemplate.TemplateSection]string)
		templatesForLayer, layoutTemplate := tmplReg.GetTemplatesForModelFrameLayer(l)
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

		out[l] = buff.String()
	}

	return out, nil
}
