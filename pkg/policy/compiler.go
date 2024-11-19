package policy

import (
	"errors"

	"github.com/eddycharly/generic-auth-server/apis/v1alpha1"
	engine "github.com/eddycharly/generic-auth-server/pkg/authz/cel"
	"github.com/eddycharly/generic-auth-server/pkg/authz/cel/libs/http"
	"github.com/eddycharly/generic-auth-server/pkg/authz/cel/libs/model"
	"github.com/eddycharly/generic-auth-server/pkg/authz/cel/utils"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apiserver/pkg/cel/lazy"
)

const (
	VariablesKey = "variables"
	ObjectKey    = "object"
)

type PolicyFunc func(*http.Request) (*model.Response, error)

type Compiler interface {
	Compile(v1alpha1.AuthorizationPolicy) (PolicyFunc, error)
}

func NewCompiler() Compiler {
	return &compiler{}
}

type compiler struct{}

func (c *compiler) Compile(policy v1alpha1.AuthorizationPolicy) (PolicyFunc, error) {
	variables := map[string]cel.Program{}
	var authorizations []cel.Program
	base, err := engine.NewEnv()
	if err != nil {
		return nil, err
	}
	provider := engine.NewVariablesProvider(base.CELTypeProvider())
	env, err := base.Extend(
		cel.Variable(ObjectKey, http.RequestType),
		cel.Variable(VariablesKey, engine.VariablesType),
		cel.CustomTypeProvider(provider),
	)
	if err != nil {
		return nil, err
	}
	for _, variable := range policy.Spec.Variables {
		ast, issues := env.Compile(variable.Expression)
		if err := issues.Err(); err != nil {
			return nil, err
		}
		provider.RegisterField(variable.Name, ast.OutputType())
		prog, err := env.Program(ast)
		if err != nil {
			return nil, err
		}
		variables[variable.Name] = prog
	}
	for _, rule := range policy.Spec.Authorizations {
		ast, issues := env.Compile(rule.Expression)
		if err := issues.Err(); err != nil {
			return nil, err
		}
		if !ast.OutputType().IsExactType(model.ResponseType) {
			return nil, errors.New("rule output is expected to be of type model.Response")
		}
		prog, err := env.Program(ast)
		if err != nil {
			return nil, err
		}
		authorizations = append(authorizations, prog)
	}
	eval := func(r *http.Request) (*model.Response, error) {
		vars := lazy.NewMapValue(engine.VariablesType)
		data := map[string]any{
			ObjectKey:    r,
			VariablesKey: vars,
		}
		for name, variable := range variables {
			vars.Append(name, func(*lazy.MapValue) ref.Val {
				out, _, err := variable.Eval(data)
				if out != nil {
					return out
				}
				if err != nil {
					return types.WrapErr(err)
				}
				return nil
			})
		}
		for _, rule := range authorizations {
			// evaluate the rule
			out, _, err := rule.Eval(data)
			// check error
			if err != nil {
				return nil, err
			}
			// evaluation result is nil, continue
			if _, ok := out.(types.Null); ok {
				continue
			}
			// try to convert to a check response
			response, err := utils.ConvertToNative[*model.Response](out)
			// check error
			if err != nil {
				return nil, err
			}
			// evaluation result is nil, continue
			if response == nil {
				continue
			}
			// no error and evaluation result is not nil, return
			return response, nil
		}
		return nil, nil
	}
	return func(r *http.Request) (*model.Response, error) {
		response, err := eval(r)
		if err != nil && policy.Spec.GetFailurePolicy() == admissionregistrationv1.Fail {
			return nil, err
		}
		return response, nil
	}, nil
}
