package generator

import (
	"bytes"
	"html/template"

	"github.com/pkg/errors"
	modelTemplate "gitlab.innovationup.stream/innovation-upstream/gen-model-frame/model/generator/template"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/model/model_function"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/model/transform"
)

type goTemplateGenerator struct {
	transform transform.ModelFramePathGoTemplateTransformer
}

func NewGoTemplateGenerator(transform transform.ModelFramePathGoTemplateTransformer) Generator {
	return &goTemplateGenerator{
		transform,
	}
}

func (g *goTemplateGenerator) Generate(mod model_function.ModelFramePath) (map[model_function.ModelFrameLayerType]string, error) {
	out := make(map[model_function.ModelFrameLayerType]string)
	out, err := g.generateDataRepo(mod)
	if err != nil {
		return out, errors.WithStack(err)
	}

	return out, nil
}

func (g *goTemplateGenerator) generateDataRepo(framePath model_function.ModelFramePath) (map[model_function.ModelFrameLayerType]string, error) {
	out := make(map[model_function.ModelFrameLayerType]string)
	layers := []model_function.ModelFrameLayerType{
		model_function.ModelFrameLayerTypeIO,
	}
	//TODO: implement the rest of the dataframe types + handle different model frame types in template inline logic
	switch framePath.DataFrameType {
	case model_function.DataFrameTypeRepo:
		layers = append(layers, []model_function.ModelFrameLayerType{
			model_function.ModelFrameLayerTypeDataRepo,
			model_function.ModelFrameLayerTypeLogicRepo,
		}...)
		break
	case model_function.DataFrameTypeRelay:
		layers = append(layers, []model_function.ModelFrameLayerType{
			model_function.ModelFrameLayerTypeDataRelay,
			model_function.ModelFrameLayerTypeLogicRelay,
		}...)
		break
	case model_function.DataFrameTypeClient:
		layers = append(layers, []model_function.ModelFrameLayerType{
			model_function.ModelFrameLayerTypeDataClient,
			model_function.ModelFrameLayerTypeLogicClient,
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

			layerTmplSections[section] = buff.String()
		}
		t := template.Must(template.New("mod_fp_tmpl").Parse(layoutTemplate))
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
