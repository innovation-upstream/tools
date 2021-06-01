package io

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"regexp"
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
	functionContent, err := goGen.Generate(o.Model)
	if err != nil {
		return errors.WithStack(err)
	}

	if cfg.OutputDirectory == "" {
		raw, err := json.Marshal(functionContent)
		if err != nil {
			return errors.WithStack(err)
		}

		os.Stdout.Write(raw)
		return nil
	}

	var sb strings.Builder
	for _, v := range functionContent {
		for moduleName, moduleContent := range v {
			for layer, moduleContent := range moduleContent {
				sb.Reset()

				strLayer := string(layer)
				var reDash = regexp.MustCompile(`-`)
				strLayer = reDash.ReplaceAllString(strLayer, "_")

				baseDirOverride := o.Model.Metadata[model.ModelMetadataOutputBaseDirectory]
				if baseDirOverride != "" {
					sb.WriteString(baseDirOverride)
					sb.WriteRune('/')
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
				sb.WriteString(fmt.Sprintf("%s_%s.go", strLayer, moduleName))

				err = ioutil.WriteFile(sb.String(), []byte(moduleContent), fs.FileMode(0644))
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}
	}

	return nil
}
