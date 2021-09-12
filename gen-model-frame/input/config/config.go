package config

import (
	"github.com/iancoleman/strcase"
	"innovationup.stream/tools/gen-model-frame/core/label"
	"innovationup.stream/tools/gen-model-frame/core/regexp"
	"unknwon.dev/clog/v2"
)

type (
	ModelFrameGenConfig struct {
		Output         ConfigOutput
		ModelsFilePath string `json:"modelsFilePath"`
	}

	ConfigOutput struct {
		Target                                 ConfigOutputTarget
		GlobalPrefix                           ModelFilePathTemplate        `json:"globalPrefix"`
		ModuleLayerImplementationFileOverrides []ConfigOutputModuleOverride `json:"module"`
	}

	ConfigOutputTarget string

	ConfigOutputModuleOverride struct {
		Label  label.ModelFrameResourceLabel         `json:"label"`
		Prefix ModelFilePathTemplate                 `json:"prefix"`
		Files  []ConfigOutputModuleLayerFileOverride `json:"file"`
	}

	ConfigOutputModuleLayerFileOverride struct {
		Label        label.ModelFrameResourceLabel `json:"label"`
		PathTemplate ModelFilePathTemplate         `json:"pathTemplate"`
	}

	ModelFilePathTemplate string
)

const (
	ConfigOutputTargetFileSystem = ConfigOutputTarget("filesystem")

	ConfigOutputTargetStdout = ConfigOutputTarget("stdout")
)

func (p ModelFilePathTemplate) Compile(modelLabel label.ModelLabel, layerLbl label.ModelFrameResourceLabel, implLabel label.ModelFrameResourceLabel) string {
	valid := regexp.ModelFilePathTemplatePattern.MatchString(string(p))
	if !valid {
		clog.Error("Module layer path templates must match the pattern: %s", regexp.ModelFilePathTemplatePattern.String())
		clog.Fatal("Failed to compile module layer path template: %s", string(p))
	}

	ct := string(p)

	ct = regexp.FilePathTemplateModelNameSnakeMergeFieldPattern.ReplaceAllString(ct, strcase.ToSnake(modelLabel.GetName()))

	ct = regexp.FilePathTemplateModelNameKebabMergeFieldPattern.ReplaceAllString(ct, strcase.ToKebab(modelLabel.GetName()))

	ct = regexp.FilePathTemplateLayerImplementationLabelKebabMergeFieldPattern.ReplaceAllString(ct, strcase.ToKebab(implLabel.GetResourceName()))

	ct = regexp.FilePathTemplateLayerLabelKebabMergeFieldPattern.ReplaceAllString(ct, strcase.ToKebab(layerLbl.GetResourceName()))

	return ct
}

func (c ConfigOutput) GetOutputForLayerImplementation(l label.ModelFrameResourceLabel) (ModelFilePathTemplate, ConfigOutputModuleLayerFileOverride) {
	var override ConfigOutputModuleLayerFileOverride
	var modulePrefix ModelFilePathTemplate

	for _, o := range c.ModuleLayerImplementationFileOverrides {
		if o.Label.GetNamespace() == l.GetNamespace() {
			modulePrefix = o.Prefix
			for _, f := range o.Files {
				if f.Label == l {
					override = f
				}
			}
		}
	}

	return modulePrefix, override
}
