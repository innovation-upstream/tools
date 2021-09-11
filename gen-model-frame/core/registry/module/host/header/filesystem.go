package header

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
