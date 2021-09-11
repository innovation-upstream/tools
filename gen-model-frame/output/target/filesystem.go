package target

import (
	"strings"

	"innovationup.stream/tools/gen-model-frame/core/model"
	"innovationup.stream/tools/gen-model-frame/core/module"
	"innovationup.stream/tools/gen-model-frame/input/config"
)

type (
	FileSystemOutputTarget interface {
		GetLayerOutputPath(layer *module.ModelLayer) string
	}

	fileSystemOutputTarget struct {
		ModelLabel   model.ModelLabel
		ConfigOutput config.ConfigOutput
	}

	FileSystemOutputTargetFactory func(BaseOutputDirectory config.ModelFilePathTemplate, ModelLabel model.ModelLabel) FileSystemOutputTarget
)

func NewFileSystemOutputTarget(ModelLabel model.ModelLabel, configOutput config.ConfigOutput) FileSystemOutputTarget {
	return &fileSystemOutputTarget{
		ModelLabel:   ModelLabel,
		ConfigOutput: configOutput,
	}
}

func (o *fileSystemOutputTarget) GetLayerOutputPath(layer *module.ModelLayer) string {
	var sb strings.Builder
	globalOutPrefix := o.ConfigOutput.GlobalPrefix
	moduleOutPrefix, override := o.ConfigOutput.GetOutputForLayer(layer.Label)
	layerPathTpl := layer.File.PathTemplate

	if override.PathTemplate != "" {
		layerPathTpl = override.PathTemplate
	}

	if globalOutPrefix != "" {
		sb.WriteString(o.ConfigOutput.GlobalPrefix.Compile(o.ModelLabel))
		sb.WriteRune('/')
	}

	if moduleOutPrefix != "" {
		sb.WriteString(moduleOutPrefix.Compile(o.ModelLabel))
		sb.WriteRune('/')
	}

	sb.WriteString(layerPathTpl.Compile(o.ModelLabel))

	outDir := sb.String()

	return outDir
}
