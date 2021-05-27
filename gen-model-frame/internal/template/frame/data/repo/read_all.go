package template

const InterfaceRepoReadAll = `
GetBy{{.ReferenceTypeCamel}}Reference(ctx context.Context, {{.ReferenceTypeLowerCamel}}ID string) ([]*model.{{.ModCamel}}, error)
`

const MethodRepoReadAll = `
func (r *{{.ModLowerCamel}}Repo) GetBy{{.ReferenceTypeCamel}}Reference(ctx context.Context, {{.ReferenceTypeLowerCamel}}IDs []string) ([]*model.{{.ModCamel}}, error) {
  var data []*model.{{.ModCamel}}
  var genericData []interface{}
  _, err := r.repo.Get(ctx, field.PurposeReference{{.ReferenceTypeCamel}}, {{.ReferenceTypeLowerCamel}}IDs, &genericData)
  if err != nil {
          return data, errors.WithStack(err)
  }

  for _, d := range genericData {
    var d model.{{.ModCamel}}
    decode, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
      TagName: "firestore",
      Result:  &d,
    })
    if err != nil {
      return data, errors.WithStack(err)
    }

    err = decode.Decode(d)
    if err != nil {
      return data, errors.WithStack(err)
    }

    data = append(data, &data)
  }

  return data, nil
}
`
