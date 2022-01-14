package check

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"unknwon.dev/clog/v2"
)

type (
	bazelCheck struct {
		targetScope string
	}
)

var isBazelLabel = regexp.MustCompile(`^\/\/.*:.*$`)
var removeBazelLabelName = regexp.MustCompile(`:.*$`)

func NewBazelCheck(targetScope string) Check {
	return &bazelCheck{
		targetScope,
	}
}

func (c *bazelCheck) GetPotentiallyBrokenConsumers(
	workspaceFilePaths []string,
	allConsumers []string,
) ([]string, error) {
	var dependantBins []string
	// TODO: extract into a func and use go routines to speed up (may need to use different bazel cache patsh to allow concurrency
	for _, f := range workspaceFilePaths {
		if f == "" {
			continue
		}

		var removeFile = regexp.MustCompile(`.[^/]*$`)
		var getBazelLabelPath = regexp.MustCompile(`^\/\/`)
		var inTargetPath = regexp.MustCompile(c.targetScope)

		path := removeFile.ReplaceAllString(f, ":all")

		if !inTargetPath.Match([]byte(getBazelLabelPath.ReplaceAllString(path, ""))) {
			clog.Trace("Skipping %s", f)
			continue
		}

		rDepsCmd := exec.Command(
			"bazel",
			"query",
			"--output",
			"label",
			fmt.Sprintf("rdeps(%s, //%s)", c.targetScope, path),
		)
		rDepsOut, err := rDepsCmd.CombinedOutput()
		if err != nil {
			clog.Warn(
				"Bazel query failed on path: %s with error: %+v",
				path,
				errors.New(string(rDepsOut)),
			)
			continue
		}

		dependentConsumers := strings.Split(string(rDepsOut), "\n")
		for _, dt := range dependentConsumers {
			if isBazelLabel.MatchString(dt) {
				d := removeBazelLabelName.ReplaceAllString(dt, "")
			consumers:
				for _, b := range allConsumers {
					if b == d {
						binLabelPath := getBazelLabelPath.ReplaceAllString(b, "")
						// de-dupe bin label paths and append them to dependantBins array
						for _, dupe := range dependantBins {
							if dupe == binLabelPath {
								continue consumers
							}
						}

						dependantBins = append(dependantBins, binLabelPath)
					}
				}
			}
		}
	}

	return dependantBins, nil
}
