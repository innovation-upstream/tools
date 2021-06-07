package check

type (
	Check interface {
		GetPotentiallyBrokenConsumers(workspaceFilePaths []string, allConsumers []string) ([]string, error)
	}
)
