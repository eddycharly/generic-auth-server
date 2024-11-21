package policy

import (
	"github.com/eddycharly/generic-auth-server/apis/v1alpha1"
	engine "github.com/eddycharly/generic-auth-server/pkg/auth/cel"
	"github.com/eddycharly/generic-auth-server/pkg/auth/cel/libs/auth"
	"github.com/eddycharly/generic-auth-server/pkg/auth/cel/libs/http"
	"github.com/eddycharly/generic-auth-server/pkg/auth/cel/utils"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/cel/lazy"
)

const (
	VariablesKey = "variables"
	ObjectKey    = "object"
)

type PolicyFunc func(*http.Request) (*auth.Response, error)

type Compiler interface {
	Compile(*v1alpha1.AuthorizationPolicy) (PolicyFunc, field.ErrorList)
}

func NewCompiler() Compiler {
	return &compiler{}
}

type compiler struct{}

func (c *compiler) Compile(policy *v1alpha1.AuthorizationPolicy) (PolicyFunc, field.ErrorList) {
	var allErrs field.ErrorList
	variables := map[string]cel.Program{}
	var authorizations []cel.Program
	base, err := engine.NewEnv()
	if err != nil {
		return nil, append(allErrs, field.InternalError(nil, err))
	}
	provider := engine.NewVariablesProvider(base.CELTypeProvider())
	env, err := base.Extend(
		cel.Variable(ObjectKey, http.RequestType),
		cel.Variable(VariablesKey, engine.VariablesType),
		cel.CustomTypeProvider(provider),
	)
	if err != nil {
		return nil, append(allErrs, field.InternalError(nil, err))
	}
	path := field.NewPath("spec")
	{
		path := path.Child("variables")
		for i, variable := range policy.Spec.Variables {
			path := path.Index(i)
			ast, issues := env.Compile(variable.Expression)
			if err := issues.Err(); err != nil {
				return nil, append(allErrs, field.Invalid(path.Child("expression"), variable.Expression, err.Error()))
			}
			provider.RegisterField(variable.Name, ast.OutputType())
			prog, err := env.Program(ast)
			if err != nil {
				return nil, append(allErrs, field.Invalid(path.Child("expression"), variable.Expression, err.Error()))
			}
			variables[variable.Name] = prog
		}
	}
	{
		path := path.Child("authorizations")
		for i, rule := range policy.Spec.Authorizations {
			path := path.Index(i)
			ast, issues := env.Compile(rule.Expression)
			if err := issues.Err(); err != nil {
				return nil, append(allErrs, field.Invalid(path.Child("expression"), rule.Expression, err.Error()))
			}
			if !ast.OutputType().IsExactType(auth.ResponseType) {
				return nil, append(allErrs, field.Invalid(path.Child("expression"), rule.Expression, "rule output is expected to be of type auth.Response"))
			}
			prog, err := env.Program(ast)
			if err != nil {
				return nil, append(allErrs, field.Invalid(path.Child("expression"), rule.Expression, err.Error()))
			}
			authorizations = append(authorizations, prog)
		}
	}
	eval := func(r *http.Request) (*auth.Response, error) {
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
			response, err := utils.ConvertToNative[*auth.Response](out)
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
	return func(r *http.Request) (*auth.Response, error) {
		response, err := eval(r)
		if err != nil && policy.Spec.GetFailurePolicy() == admissionregistrationv1.Fail {
			return nil, err
		}
		return response, nil
	}, nil
}
