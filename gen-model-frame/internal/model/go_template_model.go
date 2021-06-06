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

	BasicLayoutTemplateInput struct {
		ModCamel      string
		ModLowerCamel string
		ModSnake      string
	}

	GoBasicLayoutTemplateInput struct {
		Basic    BasicLayoutTemplateInput
		Sections map[string]string
	}
)
