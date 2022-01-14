package cmd

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/tools/bazel-check-breaking/internal/changes"
	"gitlab.innovationup.stream/innovation-upstream/tools/bazel-check-breaking/internal/check"
	"unknwon.dev/clog/v2"
)

var isBazelLabel = regexp.MustCompile(`^\/\/.*:.*$`)
var removeBazelLabelName = regexp.MustCompile(`:.*$`)

func CheckBreaking(
	fromSHA string,
	toSHA string,
	bazelTargetScope string,
) ([]string, error) {
	// Get the files we changed in to SHA commit compared to from SHA commit
	g := changes.NewGitChanges(fromSHA, toSHA)
	changedFiles, err := g.GetChangedFiles()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	clog.Trace("Changed files: %+v", changedFiles)

	//TODO: de-dupe git files so we just check pkgs and not every file

	// TODO: move this to a struct
	// Get all go binaries
	bazelBinsCmd := exec.Command(
		"bazel",
		"query",
		fmt.Sprintf("kind(go_binary.*, %s)", bazelTargetScope),
		"--output",
		"label",
	)
	bazelBinsOut, err := bazelBinsCmd.CombinedOutput()
	if err != nil {
		clog.Error(
			"Bazel go_binary query error: %+v",
			errors.New(string(bazelBinsOut)),
		)
		return nil, errors.WithStack(err)
	}

	allBinsTargets := strings.Split(string(bazelBinsOut), "\n")
	var allBins []string
	for _, t := range allBinsTargets {
		if isBazelLabel.MatchString(t) {
			s := removeBazelLabelName.ReplaceAllString(t, "")
			// TODO: de-dupe allBins because this query can match go image in addition to the go binary target
			allBins = append(allBins, s)
		}
	}

	clog.Trace("Found project binaries: %+v", allBins)

	c := check.NewBazelCheck(bazelTargetScope)
	potentiallyBrokenConsumers, err :=
		c.GetPotentiallyBrokenConsumers(changedFiles, allBins)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	clog.Trace(
		"Potentially broken dependant labels: %+v",
		potentiallyBrokenConsumers,
	)

	return potentiallyBrokenConsumers, nil
}
