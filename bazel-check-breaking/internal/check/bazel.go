package check

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"
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

func (c *bazelCheck) GetPotentiallyBrokenConsumers(workspaceFilePaths []string, allConsumers []string) []string {
	var dependantBins []string
	for _, f := range workspaceFilePaths {
		if f == "" {
			return dependantBins
		}

		var removeFile = regexp.MustCompile(`.[^/]*$`)
		var getBazelLabelPath = regexp.MustCompile(`^\/\/`)

		path := removeFile.ReplaceAllString(f, ":all")
		rDepsCmd := exec.Command("bazel", "query", "--output", "label", fmt.Sprintf("rdeps(%s, //%s)", c.targetScope, path))
		rDepsOut, err := rDepsCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("err: %+v", string(rDepsOut))
			panic(errors.WithStack(err))
		}

		dependentConsumers := strings.Split(string(rDepsOut), "\n")
		//fmt.Printf("%+v\n", dependentConsumers)
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

	return dependantBins
}
