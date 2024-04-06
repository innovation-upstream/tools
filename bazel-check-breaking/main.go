package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/innovation-upstream/tools/bazel-check-breaking/internal/cmd"
	"unknwon.dev/clog/v2"
)

var bazelTargetScopeFlag = flag.String(
	"bazel-target-scope",
	"//...",
	"specify the bazel target to use as the scope for queries",
)
var bazelKindFlag = flag.String(
	"bazel-kind",
	"go_binary",
	"specify the kind of bazel target check for breaking changes",
)
var fromSHAFlag = flag.String(
	"from-sha",
	"",
	"specify the base git commit to use in git diff-tree",
)
var toSHAFlag = flag.String(
	"to-sha",
	"HEAD",
	"specify the current git commit to use in git diff-tree",
)
var verboseFlag = flag.Bool(
	"verbose",
	false,
	"verbose log level",
)

func main() {
	flag.Parse()
	level := getLogLevel(getIsVerbose())
	err := clog.NewConsole(
		100,
		clog.ConsoleConfig{
			Level: level,
		},
	)
	if err != nil {
		fmt.Printf("%+v", fmt.Errorf("%+v", err))
		os.Exit(1)
	}
	defer clog.Stop()

	bazelKind := *bazelKindFlag
	bazelTargetScope := *bazelTargetScopeFlag
	toSHA := getStringFlag(toSHAFlag)
	fromSHA := getStringFlag(fromSHAFlag)

	if fromSHA == "" {
		clog.Fatal("Missing --from-sha flag")
	}

	potentiallyBrokenConsumers, err :=
		cmd.CheckBreaking(fromSHA, toSHA, bazelKind, bazelTargetScope)
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

func getStringFlag(f *string) string {
	if f == nil {
		return ""
	}

	return *f
}

func getIsVerbose() bool {
	if verboseFlag == nil {
		return false
	}

	return *verboseFlag
}

func getLogLevel(isVerbose bool) clog.Level {
	if isVerbose {
		return clog.LevelTrace
	}

	return clog.LevelError
}
