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
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module/registry"
	tmplRegistry "gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/template/registry"
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

	modules, err := an.GetModules()
	if err != nil {
		return errors.WithStack(err)
	}

	var content []map[label.ModelFrameResourceLabel]code_layer_generator.ModuleCodeLayers
	// todo: move this into another struct
	tmplReg := tmplRegistry.NewFileSystemTemplateRegistry()
	tmplLoader := module.NewTemplateLoader(modules, tmplReg)
	for _, fp := range o.Model.FramePaths {
		moduleTemplates, err := tmplLoader.LoadModuleTemplates(fp)
		if err != nil {
			return errors.WithStack(err)
		}

		tr := transform.NewModelFramePathGoTemplateTransformer(&o.Model)
		clg := code_layer_generator.NewTemplateHydrator(tr, moduleTemplates)

		gen := generator.NewModelFrameGenerator(clg)
		framePathContent, err := gen.GenerateFrames(fp)
		if err != nil {
			return errors.WithStack(err)
		}

		content = append(content, framePathContent)
	}

	for _, c := range content {
		// if there is no out dir, dump to stdout
		if o.Config.OutputDirectory == "" {
			raw, err := json.Marshal(c)
			if err != nil {
				return errors.WithStack(err)
			}

			os.Stdout.Write(raw)
			return nil
		}

		var sb strings.Builder
		for _, v := range c {
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
	}

	return nil
}
