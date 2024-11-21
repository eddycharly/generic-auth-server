package auth

import (
	"reflect"

	"github.com/eddycharly/generic-auth-server/pkg/auth/cel/utils"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
)

type lib struct{}

func Lib() cel.EnvOption {
	// create the cel lib env option
	return cel.Lib(&lib{})
}

func (*lib) LibraryName() string {
	return "kyverno.auth"
}

func (c *lib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		// register response type
		ext.NativeTypes(reflect.TypeFor[Response]()),
		// extend environment with function overloads
		c.extendEnv,
	}
}

func (*lib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

func (*lib) extendEnv(env *cel.Env) (*cel.Env, error) {
	// get env type adapter
	adapter := env.CELTypeAdapter()
	// build our function overloads
	libraryDecls := map[string][]cel.FunctionOpt{
		"auth.Response": {
			cel.Overload("response_int", []*cel.Type{types.IntType}, ResponseType, cel.UnaryBinding(response_int(adapter))),
			cel.Overload("response_int_string", []*cel.Type{types.IntType, types.StringType}, ResponseType, cel.BinaryBinding(response_int_string(adapter))),
			cel.Overload("response_int_bytes", []*cel.Type{types.IntType, types.BytesType}, ResponseType, cel.BinaryBinding(response_int_bytes(adapter))),
		},
		"WithBody": {
			cel.MemberOverload("response_withbody_string", []*cel.Type{ResponseType, types.StringType}, ResponseType, cel.BinaryBinding(response_withbody_string(adapter))),
			cel.MemberOverload("response_withbody_bytes", []*cel.Type{ResponseType, types.BytesType}, ResponseType, cel.BinaryBinding(response_withbody_bytes(adapter))),
		},
		"WithHeader": {
			cel.MemberOverload("response_withheader_string_string", []*cel.Type{ResponseType, types.StringType, types.StringType}, ResponseType, cel.FunctionBinding(response_withheader_string_string(adapter))),
		},
	}
	// create env options corresponding to our function overloads
	options := []cel.EnvOption{}
	for name, overloads := range libraryDecls {
		options = append(options, cel.Function(name, overloads...))
	}
	// extend environment with our function overloads
	return env.Extend(options...)
}

func response_int(adapter types.Adapter) func(ref.Val) ref.Val {
	return func(code ref.Val) ref.Val {
		if code, err := utils.ConvertToNative[int](code); err != nil {
			return types.WrapErr(err)
		} else {
			return adapter.NativeToValue(makeResponse(code, nil))
		}
	}
}

func response_int_string(adapter types.Adapter) func(ref.Val, ref.Val) ref.Val {
	return func(code ref.Val, body ref.Val) ref.Val {
		if code, err := utils.ConvertToNative[int](code); err != nil {
			return types.WrapErr(err)
		} else if body, err := utils.ConvertToNative[string](body); err != nil {
			return types.WrapErr(err)
		} else {
			return adapter.NativeToValue(makeResponse(code, []byte(body)))
		}
	}
}

func response_int_bytes(adapter types.Adapter) func(ref.Val, ref.Val) ref.Val {
	return func(code ref.Val, body ref.Val) ref.Val {
		if code, err := utils.ConvertToNative[int](code); err != nil {
			return types.WrapErr(err)
		} else if body, err := utils.ConvertToNative[[]byte](body); err != nil {
			return types.WrapErr(err)
		} else {
			return adapter.NativeToValue(makeResponse(code, body))
		}
	}
}

func response_withbody_string(adapter types.Adapter) func(ref.Val, ref.Val) ref.Val {
	return func(response ref.Val, body ref.Val) ref.Val {
		if response, err := utils.ConvertToNative[Response](response); err != nil {
			return types.WrapErr(err)
		} else if body, err := utils.ConvertToNative[string](body); err != nil {
			return types.WrapErr(err)
		} else {
			response.Body = []byte(body)
			return adapter.NativeToValue(response)
		}
	}
}

func response_withbody_bytes(adapter types.Adapter) func(ref.Val, ref.Val) ref.Val {
	return func(response ref.Val, body ref.Val) ref.Val {
		if response, err := utils.ConvertToNative[Response](response); err != nil {
			return types.WrapErr(err)
		} else if body, err := utils.ConvertToNative[[]byte](body); err != nil {
			return types.WrapErr(err)
		} else {
			response.Body = body
			return adapter.NativeToValue(response)
		}
	}
}

func response_withheader_string_string(adapter types.Adapter) func(...ref.Val) ref.Val {
	return func(values ...ref.Val) ref.Val {
		if response, err := utils.ConvertToNative[Response](values[0]); err != nil {
			return types.WrapErr(err)
		} else if key, err := utils.ConvertToNative[string](values[1]); err != nil {
			return types.WrapErr(err)
		} else if value, err := utils.ConvertToNative[string](values[2]); err != nil {
			return types.WrapErr(err)
		} else {
			response.Header.Add(key, value)
			return adapter.NativeToValue(response)
		}
	}
}
