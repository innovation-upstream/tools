package filesystem

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	"innovationup.stream/tools/gen-model-frame/core/label"
	"innovationup.stream/tools/gen-model-frame/input/config"
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

	for _, m := range cfg.Output.ModuleLayerImplementationFileOverrides {
		for _, f := range m.Files {
			f.Label = label.NewModelFrameResourceLabel(string(f.Label), m.Label.GetNamespace(), "implementation")
		}
	}

	return cfg, nil
}
