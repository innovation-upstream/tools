package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab.innovationup.stream/innovation-upstream/tools/bazel-check-breaking/internal/cmd"
)

var bazelTargetScopeFlag = flag.String("bazel-target-scope", "//...", "specify the bazel target to use as the scope/closure for queries")
var fromSHAFlag = flag.String("from-sha", "", "specify the base git commit to use in git diff-tree")
var toSHAFlag = flag.String("to-sha", "", "specify the current git commit to use in git diff-tree")

func main() {
	flag.Parse()
	bazelTargetScope := *bazelTargetScopeFlag
	fromSHA := *fromSHAFlag
	toSHA := *toSHAFlag

	potentiallyBrokenConsumers, err := cmd.CheckBreaking(fromSHA, toSHA, bazelTargetScope)
	if err != nil {
		panic(err)
	}

	for _, t := range potentiallyBrokenConsumers {
		_, err := os.Stdout.WriteString(fmt.Sprintf("%s\n", t))
		if err != nil {
			panic(err)
		}
	}
}
