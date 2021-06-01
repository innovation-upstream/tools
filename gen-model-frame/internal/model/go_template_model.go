package model

type (
	BasicTemplateInput struct {
		ModCamel                string
		ModLowerCamel           string
		ModSnake                string
		ModKebab                string
		ReferenceTypeCamel      string
		ReferenceTypeLowerCamel string
		UpdateableFields        []string
	}

	BasicLayoutTemplateInput struct {
		ModCamel      string
		ModLowerCamel string
		ModSnake      string
	}

	GoBasicLayoutTemplateInput struct {
		Basic              BasicLayoutTemplateInput
		ModGoPackage       string
		ModelGoPackagePath string
		Sections           map[string]string
	}
)
