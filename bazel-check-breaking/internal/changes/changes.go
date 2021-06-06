package changes

type (
	Changes interface {
		GetChangedFiles() ([]string, error)
	}
)
