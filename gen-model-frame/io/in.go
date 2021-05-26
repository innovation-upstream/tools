package io

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/model"
)

func ParseModelsFile() ([]model.Model, error) {
	var mods []model.Model
	by, err := ioutil.ReadFile("models.json")
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

//TODO: parse input into model structs and return
