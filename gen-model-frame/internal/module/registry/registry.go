package registry

import "gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"

//go:generate mockgen -destination=../mock/module_registry_mock.go -package=mock gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module/registry ModuleRegistry
type ModuleRegistry interface {
	QueryAllModuleHeaders() ([]ModuleHeader, error)
	QueryModuleHeaders(l []label.ModelFrameResourceLabel) ([]ModuleHeader, error)
}
