package io

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/config"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model"
)

func ParseConfigFile() (config.ModelFrameGenConfig, error) {
	var cfg config.ModelFrameGenConfig
	by, err := ioutil.ReadFile("config.json")
	if err != nil {
		return cfg, errors.WithStack(err)
	}

	data := string(by)
	err = json.Unmarshal([]byte(data), &cfg)
	if err != nil {
		return cfg, errors.WithStack(err)
	}

	return cfg, nil
}

func ParseModelsFile(cfg config.ModelFrameGenConfig) ([]model.Model, error) {
	var mods []model.Model

	modJsonPath := "models.json"
	if cfg.ModelsFilePath != "" {
		modJsonPath = cfg.ModelsFilePath
	}

	by, err := ioutil.ReadFile(modJsonPath)
	if err != nil {
		return mods, errors.WithStack(err)
	}

	data := string(by)
	err = json.Unmarshal([]byte(data), &mods)
	if err != nil {
		return mods, errors.WithStack(err)
	}

	return mods, nil
}
