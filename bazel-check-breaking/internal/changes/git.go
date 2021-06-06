package changes

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

type (
	gitChanges struct {
		fromSHA string
		toSHA   string
	}
)

func NewGitChanges(fromSHA string, toSHA string) Changes {
	return &gitChanges{
		fromSHA,
		toSHA,
	}
}

func (c *gitChanges) GetChangedFiles() ([]string, error) {
	var changedFiles []string
	gDiffCmd := exec.Command("git", "diff-tree", "-r", "--name-only", c.fromSHA, c.toSHA)
	gOut, err := gDiffCmd.CombinedOutput()
	if err != nil {
		return changedFiles, errors.Wrap(err, string(gOut))
	}

	changedFiles = strings.Split(string(gOut), "\n")

	return changedFiles, nil
}
