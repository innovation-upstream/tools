package io

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/model"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/model/model_gen"
)

type (
	ModelOut interface {
		OutputGenerated() error
	}
	modelOut struct {
		Model model.Model
	}
)

func NewModelOut(model model.Model) ModelOut {
	return &modelOut{
		model,
	}
}

func (o *modelOut) OutputGenerated() error {
	// TODO: if golang
	goGen := model_gen.NewGolangModelGenerator()
	out, err := goGen.Generate(o.Model)
	if err != nil {
		return errors.WithStack(err)
	}

	raw, err := json.Marshal(out)
	if err != nil {
		return errors.WithStack(err)
	}

	// TODO: use layer type to generate a file name and output that to a bin dir parsed from a flag
	os.Stdout.Write(raw)

	return nil
}
