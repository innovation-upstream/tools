package regexp

import (
	"regexp"
)

var ModelFrameResourceLabelPattern = regexp.MustCompile(`^(@[A-Za-z-_]+\/)?[A-Za-z-_]+::(((function)|(section)|(layer))+\/)[A-Za-z-_]+$`)
