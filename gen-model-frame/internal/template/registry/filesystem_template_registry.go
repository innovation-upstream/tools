package registry

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"
)

type fsTemplateRegistry struct {
}

func NewFileSystemTemplateRegistry() TemplateRegistry {
	return &fsTemplateRegistry{}
}

func (l *fsTemplateRegistry) LoadSectionTemplate(functionLabel label.ModelFrameResourceLabel, layerLabel label.ModelFrameResourceLabel, sectionLabel label.ModelFrameResourceLabel) (string, error) {
	var tmpl string

	sectionTemplatePath := fmt.Sprintf("modules/%s/templates/%s/%s/%s.gotemplate", sectionLabel.GetNamespace(), layerLabel.GetResourceName(), functionLabel.GetResourceName(), sectionLabel.GetResourceName())
	by, err := ioutil.ReadFile(sectionTemplatePath)
	if err != nil {
		return tmpl, errors.WithStack(err)
	}

	tmpl = string(by)

	return tmpl, nil
}

func (l *fsTemplateRegistry) LoadLayerLayoutTemplate(layerLabel label.ModelFrameResourceLabel) (string, error) {
	var tmpl string

	sectionTemplatePath := fmt.Sprintf("modules/%s/templates/%s/layout.gotemplate", layerLabel.GetNamespace(), layerLabel.GetResourceName())
	by, err := ioutil.ReadFile(sectionTemplatePath)
	if err != nil {
		return tmpl, errors.WithStack(err)
	}

	tmpl = string(by)

	return tmpl, nil
}
