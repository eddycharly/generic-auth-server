package http

import (
	"reflect"

	"github.com/eddycharly/generic-auth-server/pkg/auth/cel/utils"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
)

var (
	stringListType = types.NewListType(types.StringType)
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
	options := []cel.EnvOption{}
	// register request type
	options = append(options, ext.NativeTypes(reflect.TypeFor[Request]()))
	// create env options corresponding to our function overloads
	libraryDecls := map[string][]cel.FunctionOpt{
		"Header": {
			cel.MemberOverload("request_header_string", []*cel.Type{RequestType, types.StringType}, stringListType, cel.BinaryBinding(request_header_string)),
		},
	}
	for name, overloads := range libraryDecls {
		options = append(options, cel.Function(name, overloads...))
	}
	return options
}

func (*lib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

func request_header_string(request ref.Val, key ref.Val) ref.Val {
	if request, err := utils.ConvertToNative[*Request](request); err != nil {
		return types.WrapErr(err)
	} else if key, err := utils.ConvertToNative[string](key); err != nil {
		return types.WrapErr(err)
	} else {
		return types.DefaultTypeAdapter.NativeToValue(request.Header.Values(key))
	}
}
