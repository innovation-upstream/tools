package template

// template fragments should return a hydrated model frame template that can
// then be passed to a final generator for compilable source code generation

const InterfaceRepoCreate = `
Create(ctx context.Context, data *model.{{ModCamel}}) error
`

const MethodRepoCreate = `
func (r *{{ModLowerCamel}}Repo) Create(ctx context.Context, data []*model.{{ModCamel}}) error {
  err := r.storage.CreateOne(ctx, data)
  if err != nil {
          return errors.WithStack(err)
  }

  return nil
}
`
