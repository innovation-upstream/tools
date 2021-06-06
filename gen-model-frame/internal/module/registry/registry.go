package registry

//go:generate mockgen -destination=../mock/module_registry_mock.go -package=mock gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module/registry ModuleRegistry
type ModuleRegistry interface {
	QueryAllModuleHeaders() ([]ModuleHeader, error)
}
