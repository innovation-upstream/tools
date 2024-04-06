package host

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/tools/gen-model-frame/core/label"
	"unknwon.dev/clog/v2"
)

type fsTemplateRegistry struct {
}

func NewFileSystemTemplateRegistry() TemplateRegistry {
	return &fsTemplateRegistry{}
}

func (l *fsTemplateRegistry) LoadTemplateForAllSections(layerLabel label.ModelFrameResourceLabel, implementationLabel label.ModelFrameResourceLabel) (map[label.ModelFrameResourceLabel]string, error) {
	tmpls := make(map[label.ModelFrameResourceLabel]string)

	sectionsDirPath := fmt.Sprintf("modules/%s/templates/%s/%s/sections", layerLabel.GetNamespace(), layerLabel.GetResourceName(), implementationLabel.GetResourceName())
	f, err := os.Open(sectionsDirPath)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			clog.Trace("Skipping loading section templates for %s (no sections dir found)", layerLabel)
			return tmpls, nil
		}
		return tmpls, errors.WithStack(err)
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		return tmpls, errors.WithStack(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".gotemplate") {
			fName := strings.Split(f.Name(), ".gotemplate")[0]
			sectionName := label.NameToModelFrameResourceLabel(layerLabel.GetNamespace(), "section", fName)
			sectionTemplatePath := fmt.Sprintf("%s/%s.gotemplate", sectionsDirPath, sectionName.GetResourceName())
			by, err := ioutil.ReadFile(sectionTemplatePath)
			if err != nil {
				return tmpls, errors.WithStack(err)
			}

			tmpls[sectionName] = string(by)
		}
	}

	return tmpls, nil
}

func (l *fsTemplateRegistry) LoadLayerLayoutTemplate(layerLabel label.ModelFrameResourceLabel, implementationLabel label.ModelFrameResourceLabel) (string, error) {
	var tmpl string

	sectionTemplatePath := fmt.Sprintf("modules/%s/templates/%s/%s/layout.gotemplate", layerLabel.GetNamespace(), layerLabel.GetResourceName(), implementationLabel.GetResourceName())
	by, err := ioutil.ReadFile(sectionTemplatePath)
	if err != nil {
		return tmpl, errors.WithStack(err)
	}

	tmpl = string(by)

	return tmpl, nil
}
