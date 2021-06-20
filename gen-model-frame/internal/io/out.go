package io

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/analyze"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/code_layer_generator"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/config"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/generator"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module/registry"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/transform"
)

type (
	ModelOut interface {
		OutputGenerated() error
	}

	modelOut struct {
		Model  model.Model
		Config config.ModelFrameGenConfig
	}
)

func NewModelOut(model model.Model, cfg config.ModelFrameGenConfig) ModelOut {
	return &modelOut{
		Model:  model,
		Config: cfg,
	}
}

func (o *modelOut) OutputGenerated() error {
	// TODO: if golang
	reg := registry.NewFileSystemModuleRegistry("modules")
	moduleLoader := module.NewModuleLoader(reg)
	an := analyze.NewModelAnalyzer(o.Model, moduleLoader)

	templatesForModules, err := an.GetModuleTemplates()
	if err != nil {
		return errors.WithStack(err)
	}

	tr := transform.NewModelFramePathGoTemplateTransformer(&o.Model)
	clg := code_layer_generator.NewTemplateHydrator(tr, templatesForModules)

	gen := generator.NewModelFrameGenerator(o.Model, an, clg)
	functionContent, err := gen.GenerateFrames()
	if err != nil {
		return errors.WithStack(err)
	}

	if o.Config.OutputDirectory == "" {
		raw, err := json.Marshal(functionContent)
		if err != nil {
			return errors.WithStack(err)
		}

		os.Stdout.Write(raw)
		return nil
	}

	var sb strings.Builder
	for _, v := range functionContent {
		for moduleName, moduleContent := range v {
			for layer, moduleContent := range moduleContent {
				sb.Reset()

				strLayer := layer.GetFileFriendlyName()

				baseDirOverride := o.Model.Output.Directory
				if baseDirOverride != "" {
					sb.WriteString(baseDirOverride)
					sb.WriteRune('/')
				}

				sb.WriteString(o.Config.OutputDirectory)
				sb.WriteRune('/')
				sb.WriteString(o.Model.Name)
				sb.WriteRune('/')
				sb.WriteString(strLayer)

				outDir := sb.String()
				err = os.MkdirAll(outDir, 0755)
				if err != nil {
					return errors.WithStack(err)
				}

				sb.WriteRune('/')
				sb.WriteString(fmt.Sprintf("%s.go", moduleName))

				err = ioutil.WriteFile(sb.String(), []byte(moduleContent), fs.FileMode(0644))
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}
	}

	return nil
}
