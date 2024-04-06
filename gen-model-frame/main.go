package main

import (
	"github.com/pkg/errors"
	"github.com/tools/gen-model-frame/input/filesystem"
	clog "unknwon.dev/clog/v2"
)

func main() {
	err := clog.NewConsole(100,
		clog.ConsoleConfig{
			Level: clog.LevelTrace,
		},
	)
	if err != nil {
		clog.Fatal("%+v", errors.WithStack(err))
	}
	defer clog.Stop()

	cfg, err := filesystem.ParseConfigFile()
	if err != nil {
		clog.Fatal("%+v", errors.WithStack(err))
	}

	mods, err := filesystem.ParseModelsFile(cfg)
	if err != nil {
		clog.Fatal("%+v", errors.WithStack(err))
	}

	for _, m := range mods {
		outGen := NewModelOut(m, cfg)
		err := outGen.OutputGenerated()
		if err != nil {
			clog.Fatal("%+v", errors.WithStack(err))
		}
	}
}
