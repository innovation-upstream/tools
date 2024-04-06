package label

import (
	"fmt"
	"regexp"
	"strings"

	gmfRegexp "github.com/tools/gen-model-frame/core/regexp"
)

type (
	ModelFrameResourceLabel string
	ModelLabel              string
)

func NewModelFrameResourceLabel(l string, ns string, typ string) ModelFrameResourceLabel {
	labelIsQualified := gmfRegexp.ModelFrameResourceLabelPattern.MatchString(l)
	if labelIsQualified == false {
		return ModelFrameResourceLabel(fmt.Sprintf("%s::%s/%s", ns, typ, l))
	}

	return ModelFrameResourceLabel(l)
}

// GetNamespace returns the namespace represented in the label.
// e.g. "@innovation-updatem/golang-api" "golang-api"
func (l ModelFrameResourceLabel) GetNamespace() string {
	return strings.Split(string(l), "::")[0]
}

// GetResourceType returns the resource type represented by the label.
// e.g. "layer"
func (l ModelFrameResourceLabel) GetResourceType() string {
	return strings.Split(strings.Split(string(l), "::")[1], "/")[0]
}

// GetResourceName returns the name of the resource represented by the label.
// e.g. "data-repo"
func (l ModelFrameResourceLabel) GetResourceName() string {
	return strings.Split(strings.Split(string(l), "::")[1], "/")[1]
}

func (n ModelFrameResourceLabel) GetFileFriendlyName() string {
	moduleName := string(n)

	var reAt = regexp.MustCompile(`^@`)
	moduleName = reAt.ReplaceAllString(moduleName, "")

	var reSlash = regexp.MustCompile(`\/`)
	moduleName = reSlash.ReplaceAllString(moduleName, "_")

	var reDash = regexp.MustCompile(`-`)
	moduleName = reDash.ReplaceAllString(moduleName, "_")

	var reColon = regexp.MustCompile(`::`)
	moduleName = reColon.ReplaceAllString(moduleName, "_")

	return moduleName
}

func NameToModelFrameResourceLabel(ns string, resourceType string, name string) ModelFrameResourceLabel {
	var label ModelFrameResourceLabel

	isAlreadyQualified := gmfRegexp.ModelFrameResourceLabelPattern.MatchString(string(name))
	if isAlreadyQualified {
		return ModelFrameResourceLabel(name)
	}

	label = ModelFrameResourceLabel(fmt.Sprintf("%s::%s/%s", ns, resourceType, name))

	return label
}

func (n ModelLabel) GetNamespace() string {
	return strings.Split(string(n), "/")[0]
}

func (n ModelLabel) GetName() string {
	s := strings.Split(string(n), "/")
	if len(s) > 1 {
		return s[1]
	}

	return string(n)
}
