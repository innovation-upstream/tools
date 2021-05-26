package model

import (
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/model/model_function"
)

type ModelMetadata string

const ModelMetadataUpdatableFields = ModelMetadata("updatable_fields")
const ModelMetadataGolangModelPackagePath = ModelMetadata("golang_model_pkg_path")
const ModelMetadataGolangModelPackage = ModelMetadata("golang_model_pkg")

type Model struct {
	Name       string
	Metadata   map[ModelMetadata]string
	FramePaths []model_function.ModelFramePath
}
