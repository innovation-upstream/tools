package main

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/tools/gen-model-frame/core/label"
	"github.com/tools/gen-model-frame/core/model"
	"github.com/tools/gen-model-frame/core/model/analyze"
	"github.com/tools/gen-model-frame/core/module"
	moduleRegistryClient "github.com/tools/gen-model-frame/core/registry/module/client"
	moduleRegistryHost "github.com/tools/gen-model-frame/core/registry/module/host"
	tmplRegistryClient "github.com/tools/gen-model-frame/core/registry/template/client"
	tmplRegistryHost "github.com/tools/gen-model-frame/core/registry/template/host"
	"github.com/tools/gen-model-frame/generator"
	"github.com/tools/gen-model-frame/generator/renderer"
	goTmplRenderer "github.com/tools/gen-model-frame/generator/renderer/go-template"
	"github.com/tools/gen-model-frame/input/config"
	"github.com/tools/gen-model-frame/output/target"
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

	var content []map[label.ModelFrameResourceLabel]renderer.RenderedModelLayers
	// todo: move this into another struct
	tmplReg := tmplRegistryHost.NewFileSystemTemplateRegistry()
	tmplLoader := tmplRegistryClient.NewTemplateLoader(modules, tmplReg)
	for _, fp := range o.Model.FramePaths {
		moduleTemplates, err := tmplLoader.LoadModuleTemplates(fp)
		if err != nil {
			return errors.WithStack(err)
		}

		tr := goTmplRenderer.NewModelFramePathGoTemplateTransformer(&o.Model)
		clg := goTmplRenderer.NewGoTemplateModelLayerRenderer(tr, moduleTemplates)

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
		for layerLbl, v := range c {
			for implLabel, codeLayerContent := range v {
				layer := searchForModelLayer(modules, layerLbl)
				impl := searchForModelLayerImplementation(modules, implLabel)

				layerFilePath := outTarget.GetLayerImplementationOutputPath(layer, impl)
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

func searchForModelLayerImplementation(modules []*module.ModelFrameModule, lbl label.ModelFrameResourceLabel) *module.ModelLayerImplementation {
	if len(modules) == 0 {
		return nil
	}

	head := modules[0]
	tail := modules[1:]

	layer := head.GetLayerImplementationByLabel(lbl)
	if layer == nil {
		layer = searchForModelLayerImplementation(tail, lbl)
	}

	return layer
}
