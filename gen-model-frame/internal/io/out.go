package io

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/analyze"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/code_layer_generator"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/code_layer_generator/gotemplate"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/config"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/generator"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/io/out/target"
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

	var content [][]code_layer_generator.CodeLayers
	// todo: move this into another struct
	tmplReg := tmplRegistry.NewFileSystemTemplateRegistry()
	tmplLoader := module.NewTemplateLoader(modules, tmplReg)
	for _, fp := range o.Model.FramePaths {
		moduleTemplates, err := tmplLoader.LoadModuleTemplates(fp)
		if err != nil {
			return errors.WithStack(err)
		}

		tr := transform.NewModelFramePathGoTemplateTransformer(&o.Model)
		clg := gotemplate.NewTemplateHydrator(tr, moduleTemplates)

		gen := generator.NewModelFrameGenerator(clg)
		framePathContent, err := gen.GenerateFrames(fp)
		if err != nil {
			return errors.WithStack(err)
		}

		content = append(content, framePathContent)
	}

	for _, c := range content {
		if o.Config.Output.Target != config.ConfigOutputTargetFileSystem {
			raw, err := json.Marshal(c)
			if err != nil {
				return errors.WithStack(err)
			}

			os.Stdout.Write(raw)
			return nil
		}

		outTarget := target.NewFileSystemOutputTarget(o.Model.Label, o.Config.Output)
		for _, v := range c {
			for layerLbl, codeLayerContent := range v {
				layer := searchForModelLayer(modules, layerLbl)

				layerFilePath := outTarget.GetLayerOutputPath(layer)
				dir, _ := path.Split(layerFilePath)
				err := os.MkdirAll(dir, 0755)
				if err != nil {
					return errors.WithStack(err)
				}

				err = ioutil.WriteFile(layerFilePath, []byte(codeLayerContent), fs.FileMode(0644))
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}
	}

	return nil
}

func searchForModelLayer(modules []*module.ModelFrameModule, lbl label.ModelFrameResourceLabel) *module.ModelLayer {
	if len(modules) == 0 {
		return nil
	}

	head := modules[0]
	tail := modules[1:]

	layer := head.GetLayerByLabel(lbl)
	if layer == nil {
		layer = searchForModelLayer(tail, lbl)
	}

	return layer
}
