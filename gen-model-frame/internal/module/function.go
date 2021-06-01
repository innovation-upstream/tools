package module

import "strings"

type ModelFunctionLabel string

func (l ModelFunctionLabel) GetName() string {
	if strings.Contains(string(l), "::") {
		return strings.Split(string(l), "::")[1]
	}
	return string(l)
}
