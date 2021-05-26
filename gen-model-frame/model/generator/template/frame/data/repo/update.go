package template

const InterfaceRepoUpdate = `
UpdateBy{{ReferenceTypeCamel}}Reference(ctx context.Context, data map[string]*model.{{ModCamel}}) error
`

const MethodRepoUpdate = `
func (r *{{ModLowerCamel}}Repo) UpdateBy{{ReferenceTypeCamel}}Reference(ctx context.Context, data map[string]*model.{{ModCamel}}) error {
  var opts []query.QueryOption
  {{if Attributes.Update.UpdatableFields}}
  var updatableFields []string
  {{range Attributes.Update.UpdatableFields}}
  updatableFields = append(updatableFields, {{.}})
  {{end}}
  opts = append(opts, query.WithFields(updatableFields)
  {{end}}

  for ref, m := range data {
    err := r.storage.UpdateFirst(ctx, field.PurposeReference{{ReferenceTypeCamel}}, ref, m, opts...)
    if err != nil {
      return errors.WithStack(err)
    }
  }

  return nil
}
`
