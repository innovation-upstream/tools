package template

const RepoLayout = `
package repo

import (
        "context"

        "github.com/mitchellh/mapstructure"
        "github.com/pkg/errors"
        model "{{.ModelGoPackagePath}}"
        "gitlab.innovationup.stream/innovation-upstream/api-frame/service/data/storage"
        "gitlab.innovationup.stream/innovation-upstream/api-frame/service/data/storage/field"
        "gitlab.innovationup.stream/innovation-upstream/api-frame/service/data/storage/query"
)

const {{.Basic.ModCamel}}Collection = "{{.Basic.ModSnake}}"

//go:generate mockgen -destination=../mock/{{.Basic.ModSnake}}_repo_mock.go -package=mock {{.ModGoPackage}} {{.Basic.ModCamel}}Repo
type {{.Basic.ModCamel}}Repo interface {
  {{.InterfaceDefinition}}
}

type {{.Basic.ModLowerCamel}}Repo struct {
  storage storage.Storage
}

func New{{.Basic.ModCamel}}Repo(s storage.Storage) {{.Basic.ModCamel}}Repo {
  return &{{.Basic.ModLowerCamel}}Repo{
    storage: s,
  }
}

{{.Methods}}
`
