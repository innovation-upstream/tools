package regexp

import (
	"fmt"
	"regexp"
)

var ModelFrameResourceLabelPattern = regexp.MustCompile(`^(@[A-Za-z-_]+\/)?[A-Za-z-_]+::(((function)|(section)|(layer))+\/)[A-Za-z-_]+$`)

var ModelFilePathTemplateSnakeMergeFieldPattern = regexp.MustCompile(`\[modelNameSnake\]`)

var ModelFilePathTemplateKebabMergeFieldPattern = regexp.MustCompile(`\[modelNameKebab\]`)

var ModelFilePathTemplatePattern = regexp.MustCompile(fmt.Sprintf(`^((\w|-){1,}\/?)*(%s|%s)?\.?\w*$`, ModelFilePathTemplateSnakeMergeFieldPattern, ModelFilePathTemplateKebabMergeFieldPattern))
