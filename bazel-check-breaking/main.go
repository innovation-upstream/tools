package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab.innovationup.stream/innovation-upstream/tools/bazel-check-breaking/internal/cmd"
	"unknwon.dev/clog/v2"
)

var bazelTargetScopeFlag = flag.String("bazel-target-scope", "//...", "specify the bazel target to use as the scope/closure for queries")
var fromSHAFlag = flag.String("from-sha", "", "specify the base git commit to use in git diff-tree")
var toSHAFlag = flag.String("to-sha", "HEAD", "specify the current git commit to use in git diff-tree")

func main() {
	// TODO: use a flag to control the log level
	err := clog.NewConsole(100,
		clog.ConsoleConfig{
			Level: clog.LevelError,
		},
	)
	if err != nil {
		fmt.Printf("%+v", fmt.Errorf("%+v", err))
		os.Exit(1)
	}
	defer clog.Stop()

	flag.Parse()
	bazelTargetScope := *bazelTargetScopeFlag
	fromSHA := *fromSHAFlag
	toSHA := *toSHAFlag

	potentiallyBrokenConsumers, err := cmd.CheckBreaking(fromSHA, toSHA, bazelTargetScope)
	if err != nil {
		clog.Fatal("%+v", err)
	}

	for _, t := range potentiallyBrokenConsumers {
		_, err := os.Stdout.WriteString(fmt.Sprintf("%s\n", t))
		if err != nil {
			clog.Fatal("%+v", err)
		}
	}
}
