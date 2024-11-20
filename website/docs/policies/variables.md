# Variables

A `AuthorizationPolicy` can define `variables` that will be made available to all authorization rules.

Variables can be used in composition of other expressions.
Each variable is defined as a named [CEL](https://github.com/google/cel-spec) expression.
The will be available under `variables` in other expressions of the policy.

The expression of a variable can refer to other variables defined earlier in the list but not those after. Thus, variables must be sorted by the order of first appearance and acyclic.

!!!info

    The incoming http request is made available to the policy under the `object` identifier.

## Variables

```yaml
apiVersion: generic.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: demo-policy.example.com
spec:
  failurePolicy: Fail
  variables:
    # `force_authorized` references the 'x-force-authorized' header
    # from the incoming http request
  - name: force_authorized
    expression: >
      object.Header("x-force-authorized")
    # `allowed` will be `true` if `variables.force_authorized` has the
    # value 'enabled' or 'true'
  - name: allowed
    expression: >
      "enabled" in variables.force_authorized || "true" in variables.force_authorized
  authorizations:
    # make an authorisation decision based on the value of `variables.allowed`
  - expression: >
      variables.allowed
        ? auth
            .Response(200)
        : auth
            .Response(401)
            .WithBody("bye")
            .WithHeader("reason", "not allowed") 
```
