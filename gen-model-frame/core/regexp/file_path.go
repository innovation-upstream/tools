package regexp

import (
	"fmt"
	"regexp"
)

var FilePathTemplateModelNameSnakeMergeFieldPattern = regexp.MustCompile(`\[modelNameSnake\]`)

var FilePathTemplateModelNameKebabMergeFieldPattern = regexp.MustCompile(`\[modelNameKebab\]`)

var FilePathTemplateLayerImplementationLabelKebabMergeFieldPattern = regexp.MustCompile(`\[layerImplementationLabelKebab\]`)

var FilePathTemplateLayerLabelKebabMergeFieldPattern = regexp.MustCompile(`\[layerLabelKebab\]`)

var ModelFilePathTemplatePattern = regexp.MustCompile(fmt.Sprintf(`^(((\w|-){1,}\/?)*((%s|%s|%s|%s)\/?)*)*\.?\w*$`, FilePathTemplateModelNameSnakeMergeFieldPattern, FilePathTemplateModelNameKebabMergeFieldPattern, FilePathTemplateLayerImplementationLabelKebabMergeFieldPattern, FilePathTemplateLayerLabelKebabMergeFieldPattern))
