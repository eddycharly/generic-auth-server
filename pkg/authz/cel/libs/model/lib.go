package model

import (
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"
)

type lib struct{}

func Lib() cel.EnvOption {
	// create the cel lib env option
	return cel.Lib(&lib{})
}

func (*lib) LibraryName() string {
	return "kyverno.model"
}

func (c *lib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		// register token type
		ext.NativeTypes(reflect.TypeFor[Response]()),
	}
}

func (*lib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}
