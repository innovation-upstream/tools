package main

import "gitlab.innovationup.stream/innovation-upstream/gen-model-frame/io"

func main() {
	mods, err := io.ParseModelsFile()
	if err != nil {
		panic(err)
	}

	for _, m := range mods {
		outGen := io.NewModelOut(m)
		err := outGen.OutputGenerated()
		if err != nil {
			panic(err)
		}
	}
}
