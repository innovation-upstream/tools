package registry

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"
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

// TODO: extract common blocks into reusable functions so there is less redundency with QueryAllModuleHeaders
func (r *fileSystemModuleRegistry) QueryModuleHeaders(lbls []label.ModelFrameResourceLabel) ([]ModuleHeader, error) {
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

			for _, l := range lbls {
				modNS := l.GetNamespace()
				fileIsTargetModuleDir := file.Name() == modNS

				// if module is namespaced check if dir matches name with namespace
				// removed
				if fileIsTargetModuleDir == false && strings.Contains(modNS, "/") == true {
					modNameNoNS := strings.Split(modNS, "/")[1]
					fileIsTargetModuleDir = file.Name() == modNameNoNS
				}

				if fileIsTargetModuleDir {
					headers = append(headers, &FileSystemModuleHeader{
						Name:     file.Name(),
						Location: r.ModulesDir,
					})
				}
			}
		}
	}

	if headers == nil {
		// Recursivly load headers from nested modules
		for _, nsModuleDir := range nsModuleDirs {
			reg := NewFileSystemModuleRegistry(nsModuleDir)
			nsModule, err := reg.QueryModuleHeaders(lbls)
			if err != nil {
				return headers, errors.WithStack(err)
			}

			headers = append(headers, nsModule...)
		}
	}

	return headers, nil
}
