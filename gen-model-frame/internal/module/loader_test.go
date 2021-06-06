package module

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module/mock"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module/registry"
)

func TestNewModuleLoader(t *testing.T) {
	type args struct {
		registry registry.ModuleRegistry
	}

	ctrl := gomock.NewController(t)
	reg := mock.NewMockModuleRegistry(ctrl)
	want := NewModuleLoader(reg)

	tests := []struct {
		name string
		args args
		want ModuleLoader
	}{
		{
			name: "happy",
			want: want,
			args: args{
				registry: reg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewModuleLoader(tt.args.registry); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewModuleLoader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_moduleLoader_LoadSectionTemplate(t *testing.T) {
	type fields struct {
		Registry registry.ModuleRegistry
	}
	type args struct {
		moduleName    string
		functionLabel ModelFunctionLabel
		layerLabel    ModelFrameLayerLabel
		sectionLabel  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &moduleLoader{
				Registry: tt.fields.Registry,
			}
			got, err := l.LoadSectionTemplate(tt.args.moduleName, tt.args.functionLabel, tt.args.layerLabel, tt.args.sectionLabel)
			if (err != nil) != tt.wantErr {
				t.Errorf("moduleLoader.LoadSectionTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("moduleLoader.LoadSectionTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_moduleLoader_LoadLayerLayoutTemplate(t *testing.T) {
	type fields struct {
		Registry registry.ModuleRegistry
	}
	type args struct {
		moduleName string
		layerLabel ModelFrameLayerLabel
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &moduleLoader{
				Registry: tt.fields.Registry,
			}
			got, err := l.LoadLayerLayoutTemplate(tt.args.moduleName, tt.args.layerLabel)
			if (err != nil) != tt.wantErr {
				t.Errorf("moduleLoader.LoadLayerLayoutTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("moduleLoader.LoadLayerLayoutTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_moduleLoader_LoadAllModulesFromDirectory(t *testing.T) {
	type fields struct {
		Registry registry.ModuleRegistry
	}
	type args struct {
		modulesDir string
	}

	ctrl := gomock.NewController(t)
	mReg := mock.NewMockModuleRegistry(ctrl)

	mModule := &ModelFrameModule{
		Name: "@builtin/golang-api",
		Functions: []ModelFunction{
			{
				Label: "create",
			},
		},
		Layers: []ModelFrameLayer{
			{
				Label:     "data-logic",
				Functions: []ModelFunctionLabel{"@builtin/golang-api::create"},
				Sections: []ModelSection{
					{
						Label: "interface-definition",
					},
					{
						Label: "method",
					},
				},
			},
		},
	}
	mModules := []*ModelFrameModule{mModule}

	rawModule, err := json.Marshal(mModule)
	if err != nil {
		panic(err)
	}

	mHeader := mock.NewMockModuleHeader(ctrl)
	mHeader.EXPECT().GetJSON().Return(string(rawModule), nil)
	mModuleHeaders := []registry.ModuleHeader{mHeader}
	mReg.EXPECT().QueryAllModuleHeaders().Return(mModuleHeaders, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*ModelFrameModule
		wantErr bool
	}{
		{
			name: "happy",
			args: args{
				modulesDir: "modules",
			},
			want:    mModules,
			wantErr: false,
			fields: fields{
				Registry: mReg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &moduleLoader{
				Registry: tt.fields.Registry,
			}
			got, err := l.LoadAllModulesFromDirectory(tt.args.modulesDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("moduleLoader.LoadAllModulesFromDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("moduleLoader.LoadAllModulesFromDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_moduleLoader_loadModule(t *testing.T) {
	type fields struct {
		Registry registry.ModuleRegistry
	}
	type args struct {
		moduleDir  string
		moduleName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ModelFrameModule
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &moduleLoader{
				Registry: tt.fields.Registry,
			}
			got, err := l.loadModule(tt.args.moduleDir, tt.args.moduleName)
			if (err != nil) != tt.wantErr {
				t.Errorf("moduleLoader.loadModule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("moduleLoader.loadModule() = %v, want %v", got, tt.want)
			}
		})
	}
}
