package template

const InterfaceRepoDelete = `
DeleteBy{{.ReferenceTypeCamel}}Reference(ctx context.Context, {{.ReferenceTypeLowerCamel}}IDs []string) error
`

const MethodRepoDelete = `
DeleteBy{{.ReferenceTypeCamel}}Reference(ctx context.Context, {{.ReferenceTypeLowerCamel}}IDs []string) error {
  err := r.repo.Delete(ctx, field.PurposeReference{{.ReferenceTypeCamel}}, ReferenceTypeLowerCamel}}IDs)
  if err != nil {
    return err
  }

  return nil
}
`
