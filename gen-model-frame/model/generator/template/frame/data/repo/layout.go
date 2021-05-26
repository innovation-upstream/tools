package template

const RepoLayout = `
package repo

import (
        "context"

        "github.com/mitchellh/mapstructure"
        "github.com/pkg/errors"
        model "{{ModelGoPackagePath}}"
        "gitlab.innovationup.stream/innovation-upstream/api-frame/service/data/storage"
        "gitlab.innovationup.stream/innovation-upstream/api-frame/service/data/storage/field"
        "gitlab.innovationup.stream/innovation-upstream/api-frame/service/data/storage/query"
)

const {{ModCamel}}Collection = "{{ModSnake}}"

//go:generate mockgen -destination=../mock/{{ModSnake}}_repo_mock.go -package=mock {{ModGoPackage}} {{ModCamel}}Repo
type {{ModCamel}}Repo interface {
  {{"{{"}}InterfaceDefinition{{"}}"}}
}

type {{ModLowerCamel}}Repo struct {
  storage storage.Storage
}

func New{{ModCamel}}Repo(s storage.Storage) {{ModCamel}}Repo {
  return &{{ModLowerCamel}}Repo{
    storage: s,
  }
}

{{"{{"}}Methods{{"}}"}}
`
