package http

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
	return "kyverno.http"
}

func (c *lib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		// register token type
		ext.NativeTypes(reflect.TypeFor[Request]()),
	}
}

func (*lib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}
