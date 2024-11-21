package validation

import (
	"context"
	"fmt"

	"github.com/eddycharly/generic-auth-server/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func NewValidator(compile func(*v1alpha1.AuthorizationPolicy) error) *validator {
	return &validator{
		compile: compile,
	}
}

type validator struct {
	compile func(*v1alpha1.AuthorizationPolicy) error
}

func (v *validator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	policy, ok := obj.(*v1alpha1.AuthorizationPolicy)
	if !ok {
		return nil, fmt.Errorf("expected an AuthorizationPolicy object but got %T", obj)
	}
	return nil, v.validate(policy)
}

func (v *validator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	policy, ok := newObj.(*v1alpha1.AuthorizationPolicy)
	if !ok {
		return nil, fmt.Errorf("expected an AuthorizationPolicy object but got %T", newObj)
	}
	return nil, v.validate(policy)
}

func (*validator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func (v *validator) validate(policy *v1alpha1.AuthorizationPolicy) error {
	return v.compile(policy)
	// var allErrs field.ErrorList
	// return apierrors.NewInvalid(
	// 	v1alpha1.SchemeGroupVersion.WithKind("AuthorizationPolicy").GroupKind(),
	// 	policy.Name,
	// 	allErrs,
	// )
}
