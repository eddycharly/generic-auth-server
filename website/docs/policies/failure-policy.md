# Failure policy

FailurePolicy defines how to handle failures for the policy.

Failures can occur from CEL expression parse errors, type check errors, runtime errors and invalid or mis-configured policy definitions.

Allowed values are:

- `Ignore`
- `Fail`

If not set, the failure policy defaults to `Fail`.

!!!info

    FailurePolicy does not define how validations that evaluate to `false` are handled.

## Fail

```yaml
apiVersion: generic.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: demo-policy.example.com
spec:
  # if something fails the request will be denied
  failurePolicy: Fail
  variables:
  - name: force_authorized
    expression: >
      object.Header("x-force-authorized")
  - name: allowed
    expression: >
      "enabled" in variables.force_authorized || "true" in variables.force_authorized
  authorizations:
  - expression: >
      variables.allowed
        ? auth
            .Response(200)
        : auth
            .Response(401)
            .WithBody("bye")
            .WithHeader("reason", "not allowed") 
```

## Ignore

```yaml
apiVersion: generic.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: demo-policy.example.com
spec:
  # if something fails the failure will be ignored and the request will be allowed
  failurePolicy: Ignore
  variables:
  - name: force_authorized
    expression: >
      object.Header("x-force-authorized")
  - name: allowed
    expression: >
      "enabled" in variables.force_authorized || "true" in variables.force_authorized
  authorizations:
  - expression: >
      variables.allowed
        ? auth
            .Response(200)
        : auth
            .Response(401)
            .WithBody("bye")
            .WithHeader("reason", "not allowed") 
```
