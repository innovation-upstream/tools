package model

import "gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/model_frame_path"

type ModelMetadata string

const ModelMetadataUpdatableFields = ModelMetadata("updatable_fields")
const ModelMetadataGolangModelPackagePath = ModelMetadata("golang_model_pkg_path")
const ModelMetadataGolangModelPackage = ModelMetadata("golang_model_pkg")

type Model struct {
	Name       string                            `json:"name"`
	Metadata   map[ModelMetadata]string          `json:"metadata"`
	FramePaths []model_frame_path.ModelFramePath `json:"frame_paths"`
}
