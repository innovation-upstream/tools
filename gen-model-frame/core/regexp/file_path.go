package regexp

import (
	"fmt"
	"regexp"
)

var ModelFilePathTemplateSnakeMergeFieldPattern = regexp.MustCompile(`\[modelNameSnake\]`)

var ModelFilePathTemplateKebabMergeFieldPattern = regexp.MustCompile(`\[modelNameKebab\]`)

var ModelFilePathTemplatePattern = regexp.MustCompile(fmt.Sprintf(`^(((\w|-){1,}\/?)*((%s|%s)\/?)*)*\.?\w*$`, ModelFilePathTemplateSnakeMergeFieldPattern, ModelFilePathTemplateKebabMergeFieldPattern))
