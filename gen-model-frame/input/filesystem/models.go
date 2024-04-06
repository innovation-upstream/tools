package filesystem

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/tools/gen-model-frame/core/model"
	"github.com/tools/gen-model-frame/input/config"
)

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
