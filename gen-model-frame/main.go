package main

import "gitlab.innovationup.stream/innovation-upstream/gen-model-frame/io"

func main() {
	cfg, err := io.ParseConfigFile()
	if err != nil {
		panic(err)
	}

	mods, err := io.ParseModelsFile(cfg)
	if err != nil {
		panic(err)
	}

	for _, m := range mods {
		outGen := io.NewModelOut(m)
		err := outGen.OutputGenerated(cfg)
		if err != nil {
			panic(err)
		}
	}
}
