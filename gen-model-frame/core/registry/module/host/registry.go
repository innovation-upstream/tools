package registry

import (
	"github.com/tools/gen-model-frame/core/label"
	"github.com/tools/gen-model-frame/core/registry/module/host/header"
)

type ModuleRegistry interface {
	QueryAllModuleHeaders() ([]header.ModuleHeader, error)
	QueryModuleHeaders(l []label.ModelFrameResourceLabel) ([]header.ModuleHeader, error)
}
