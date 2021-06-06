package registry

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

type fileSystemModuleRegistry struct {
	ModulesDir string
}

func NewFileSystemModuleRegistry(modulesDir string) ModuleRegistry {
	return &fileSystemModuleRegistry{
		ModulesDir: modulesDir,
	}
}

func (r *fileSystemModuleRegistry) QueryAllModuleHeaders() ([]ModuleHeader, error) {
	var headers []ModuleHeader
	f, err := os.Open(r.ModulesDir)
	if err != nil {
		return headers, errors.WithStack(err)
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		return headers, errors.WithStack(err)
	}

	var nsModuleDirs []string

	for _, file := range files {
		if file.IsDir() {
			moduleJSONFilePath := fmt.Sprintf("%s/%s/module.json", r.ModulesDir, file.Name())
			// if its a namespaced module dir, load nested modules
			if _, err := os.Stat(moduleJSONFilePath); errors.Is(err, os.ErrNotExist) {
				nsModuleDir := fmt.Sprintf("%s/%s", r.ModulesDir, file.Name())
				nsModuleDirs = append(nsModuleDirs, nsModuleDir)
				continue
			}

			headers = append(headers, &FileSystemModuleHeader{
				Name:     file.Name(),
				Location: r.ModulesDir,
			})
		}
	}

	// Recursivly load headers for nested modules
	for _, nsModuleDir := range nsModuleDirs {
		reg := NewFileSystemModuleRegistry(nsModuleDir)
		nsModules, err := reg.QueryAllModuleHeaders()
		if err != nil {
			return headers, errors.WithStack(err)
		}

		headers = append(headers, nsModules...)
	}

	return headers, nil
}
