package model

type (
	BasicTemplateInput struct {
		ModCamel                string
		ModLowerCamel           string
		ModSnake                string
		ModKebab                string
		ReferenceTypeCamel      string
		ReferenceTypeLowerCamel string
		MetaData                map[string]string
	}

	GoBasicLayoutTemplateInput struct {
		Basic    BasicTemplateInput
		Sections map[string]string
	}
)
