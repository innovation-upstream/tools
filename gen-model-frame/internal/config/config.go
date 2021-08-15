package config

import (
	"github.com/iancoleman/strcase"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/regexp"
	"unknwon.dev/clog/v2"
)

type (
	ModelFrameGenConfig struct {
		OutputDirectory ModelFilePathTemplate `json:"outputDirectory"`
		ModelsFilePath  string                `json:"modelsFilePath"`
	}

	ModelFilePathTemplate string
)

func (p ModelFilePathTemplate) Compile(modelLabel model.ModelLabel) string {
	valid := regexp.ModelFilePathTemplatePattern.MatchString(string(p))
	if !valid {
		clog.Error("Module layer path templates must match the pattern: %s", regexp.ModelFilePathTemplatePattern.String())
		clog.Fatal("Failed to compile module layer path template: %s", string(p))
	}

	ct := string(p)

	ct = regexp.ModelFilePathTemplateSnakeMergeFieldPattern.ReplaceAllString(ct, strcase.ToSnake(modelLabel.GetName()))

	ct = regexp.ModelFilePathTemplateKebabMergeFieldPattern.ReplaceAllString(ct, strcase.ToKebab(modelLabel.GetName()))

	return ct
}
