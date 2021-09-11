package main

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	"innovationup.stream/tools/gen-model-frame/core/label"
	"innovationup.stream/tools/gen-model-frame/core/model"
	"innovationup.stream/tools/gen-model-frame/core/model/analyze"
	"innovationup.stream/tools/gen-model-frame/core/module"
	moduleRegistryClient "innovationup.stream/tools/gen-model-frame/core/registry/module/client"
	moduleRegistryHost "innovationup.stream/tools/gen-model-frame/core/registry/module/host"
	tmplRegistryClient "innovationup.stream/tools/gen-model-frame/core/registry/template/client"
	tmplRegistryHost "innovationup.stream/tools/gen-model-frame/core/registry/template/host"
	"innovationup.stream/tools/gen-model-frame/generator"
	"innovationup.stream/tools/gen-model-frame/generator/renderer"
	goTmplRenderer "innovationup.stream/tools/gen-model-frame/generator/renderer/go-template"
	"innovationup.stream/tools/gen-model-frame/input/config"
	"innovationup.stream/tools/gen-model-frame/output/target"
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
	reg := moduleRegistryHost.NewFileSystemModuleRegistry("modules")
	moduleLoader := moduleRegistryClient.NewModuleLoader(reg)
	an := analyze.NewModelAnalyzer(o.Model, moduleLoader)

	modules, err := an.GetModules()
	if err != nil {
		return errors.WithStack(err)
	}

	var content [][]renderer.RenderedCodeLayers
	// todo: move this into another struct
	tmplReg := tmplRegistryHost.NewFileSystemTemplateRegistry()
	tmplLoader := tmplRegistryClient.NewTemplateLoader(modules, tmplReg)
	for _, fp := range o.Model.FramePaths {
		moduleTemplates, err := tmplLoader.LoadModuleTemplates(fp)
		if err != nil {
			return errors.WithStack(err)
		}

		tr := goTmplRenderer.NewModelFramePathGoTemplateTransformer(&o.Model)
		clg := goTmplRenderer.NewTemplateHydrator(tr, moduleTemplates)

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
