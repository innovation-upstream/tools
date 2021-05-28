package io

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/config"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/generator"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/model"
)

type (
	ModelOut interface {
		OutputGenerated(cfg config.ModelFrameGenConfig) error
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

func (o *modelOut) OutputGenerated(cfg config.ModelFrameGenConfig) error {
	// TODO: if golang
	goGen := generator.NewGolangModelGenerator()
	out, err := goGen.Generate(o.Model)
	if err != nil {
		return errors.WithStack(err)
	}

	if cfg.OutputDirectory == "" {
		raw, err := json.Marshal(out)
		if err != nil {
			return errors.WithStack(err)
		}

		os.Stdout.Write(raw)
		return nil
	}

	var sb strings.Builder
	for _, v := range out {
		for layer, content := range v {
			strLayer := string(layer)
			sb.Reset()

			baseDirOverride := o.Model.Metadata[model.ModelMetadataOutputBaseDirectory]
			if baseDirOverride != "" {
				sb.WriteString(baseDirOverride)
			}

			sb.WriteString(cfg.OutputDirectory)
			sb.WriteRune('/')
			sb.WriteString(o.Model.Name)
			sb.WriteRune('/')
			sb.WriteString(strLayer)

			outDir := sb.String()
			err = os.MkdirAll(outDir, 0755)
			if err != nil {
				return errors.WithStack(err)
			}

			sb.WriteRune('/')
			sb.WriteString(fmt.Sprintf("%s.go", strLayer))

			err = ioutil.WriteFile(sb.String(), []byte(content), fs.FileMode(0444))
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return nil
}
