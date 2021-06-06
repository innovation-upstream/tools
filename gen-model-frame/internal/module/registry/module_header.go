package registry

//go:generate mockgen -destination=../mock/module_header_mock.go -package=mock gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module/registry ModuleHeader
type (
	ModuleHeader interface {
		GetName() string
		GetLocation() string
		GetJSON() (string, error)
	}

	FileSystemModuleHeader struct {
		Name     string
		Location string
	}
)
