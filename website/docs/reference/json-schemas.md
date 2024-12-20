# JSON schemas

JSON schemas for policies are available:

- [AuthorizationPolicy (v1alpha1)](https://github.com/eddycharly/generic-auth-server/blob/main/.schemas/json/authorizationpolicy-generic-v1alpha1.json)

They can be used to enable validation and autocompletion in your IDE.

## VS code

In VS code, simply add a comment on top of your YAML resources.

### AuthorizationPolicy

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/eddycharly/generic-auth-server/main/.schemas/json/authorizationpolicy-generic-v1alpha1.json
apiVersion: generic.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: demo-policy.example.com
spec:
  authorizations:
  - expression: >
      "bar" in object.Header("foo")
        ? auth
            .Response(401)
            .WithBody("bye")
            .WithHeader("xxx", "yyy")
        : auth
            .Response(200)
```
