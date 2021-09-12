package target

import (
	"strings"

	"innovationup.stream/tools/gen-model-frame/core/label"
	"innovationup.stream/tools/gen-model-frame/core/module"
	"innovationup.stream/tools/gen-model-frame/input/config"
)

type (
	FileSystemOutputTarget interface {
		GetLayerImplementationOutputPath(layer *module.ModelLayer, impl *module.ModelLayerImplementation) string
	}

	fileSystemOutputTarget struct {
		ModelLabel   label.ModelLabel
		ConfigOutput config.ConfigOutput
	}

	FileSystemOutputTargetFactory func(BaseOutputDirectory config.ModelFilePathTemplate, ModelLabel label.ModelLabel) FileSystemOutputTarget
)

func NewFileSystemOutputTarget(ModelLabel label.ModelLabel, configOutput config.ConfigOutput) FileSystemOutputTarget {
	return &fileSystemOutputTarget{
		ModelLabel:   ModelLabel,
		ConfigOutput: configOutput,
	}
}

func (o *fileSystemOutputTarget) GetLayerImplementationOutputPath(layer *module.ModelLayer, impl *module.ModelLayerImplementation) string {
	var sb strings.Builder
	globalOutPrefix := o.ConfigOutput.GlobalPrefix
	moduleOutPrefix, override := o.ConfigOutput.GetOutputForLayerImplementation(layer.Label)
	layerPathTpl := layer.PathTemplate
	implPathTpl := impl.File.PathTemplate

	if override.PathTemplate != "" {
		layerPathTpl = override.PathTemplate
	} else if implPathTpl != "" {
		layerPathTpl = implPathTpl
	}

	if globalOutPrefix != "" {
		sb.WriteString(o.ConfigOutput.GlobalPrefix.Compile(o.ModelLabel, layer.Label, impl.Label))
		sb.WriteRune('/')
	}

	if moduleOutPrefix != "" {
		sb.WriteString(moduleOutPrefix.Compile(o.ModelLabel, layer.Label, impl.Label))
		sb.WriteRune('/')
	}

	sb.WriteString(layerPathTpl.Compile(o.ModelLabel, layer.Label, impl.Label))

	outDir := sb.String()

	return outDir
}
