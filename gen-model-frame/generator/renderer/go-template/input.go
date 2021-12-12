package renderer

import "innovationup.stream/tools/gen-model-frame/core/model"

type (
	BasicTemplateInput struct {
		ModCamel      string
		ModLowerCamel string
		ModSnake      string
		ModKebab      string
		MetaData      map[string]string
		Options       model.ModelOptions
		Hooks         map[string]string
	}

	GoBasicLayoutTemplateInput struct {
		Basic    BasicTemplateInput
		Sections map[string]string
	}
)
