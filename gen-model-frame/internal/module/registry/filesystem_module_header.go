package registry

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

func (h *FileSystemModuleHeader) GetLocation() string {
	return h.Location
}

func (h *FileSystemModuleHeader) GetName() string {
	return h.Name
}

func (h *FileSystemModuleHeader) GetJSON() (string, error) {
	var moduleJSON string

	moduleJSONPath := fmt.Sprintf("%s/%s/module.json", h.Location, h.Name)
	by, err := ioutil.ReadFile(moduleJSONPath)
	if err != nil {
		return moduleJSON, errors.WithStack(err)
	}

	moduleJSON = string(by)

	return moduleJSON, nil
}
