package target

import (
	"strings"

	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/config"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module"
)

type (
	FileSystemOutputTarget interface {
		GetLayerOutputPath(file module.ModelLayerFile) string
	}

	fileSystemOutputTarget struct {
		BaseOutputDirectory config.ModelFilePathTemplate
		ModelLabel          model.ModelLabel
	}

	FileSystemOutputTargetFactory func(BaseOutputDirectory config.ModelFilePathTemplate, ModelLabel model.ModelLabel) FileSystemOutputTarget
)

func NewFileSystemOutputTarget(BaseOutputDirectory config.ModelFilePathTemplate, ModelLabel model.ModelLabel) FileSystemOutputTarget {
	return &fileSystemOutputTarget{
		BaseOutputDirectory: BaseOutputDirectory,
		ModelLabel:          ModelLabel,
	}
}

func (o *fileSystemOutputTarget) GetLayerOutputPath(file module.ModelLayerFile) string {
	var sb strings.Builder

	sb.WriteString(o.BaseOutputDirectory.Compile(o.ModelLabel))
	sb.WriteRune('/')
	sb.WriteString(file.PathTemplate.Compile(o.ModelLabel))

	outDir := sb.String()

	return outDir
}
