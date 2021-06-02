package main

import (
	"fmt"
	"os"

	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/io"
	clog "unknwon.dev/clog/v2"
)

func main() {
	err := clog.NewConsole(100,
		clog.ConsoleConfig{
			Level: clog.LevelTrace,
		},
	)
	if err != nil {
		fmt.Printf("%+v", fmt.Errorf("%+v", err))
		os.Exit(1)
	}
	defer clog.Stop()

	cfg, err := io.ParseConfigFile()
	if err != nil {
		fmt.Printf("%+v", fmt.Errorf("%+v", err))
		os.Exit(1)
	}

	mods, err := io.ParseModelsFile(cfg)
	if err != nil {
		fmt.Printf("%+v", fmt.Errorf("%+v", err))
		os.Exit(1)
	}

	for _, m := range mods {
		outGen := io.NewModelOut(m, cfg)
		err := outGen.OutputGenerated()
		if err != nil {
			fmt.Printf("%+v", fmt.Errorf("%+v", err))
			os.Exit(1)
		}
	}
}
