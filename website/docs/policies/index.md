# Policies

A Kyverno `AuthorizationPolicy` is a custom [Kubernetes resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) and can be easily managed via Kubernetes APIs, GitOps workflows, and other existing tools.

## Resource Scope

An `AuthorizationPolicy` is a cluster-wide resource.

## API Group and Kind

An `AuthorizationPolicy` belongs to the `generic.kyverno.io/v1alpha1` group and can only be of kind `AuthorizationPolicy`.

```yaml
apiVersion: generic.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: demo-policy.example.com
spec:
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

## Generic External Auth

The Generic Auth Server analyses an incoming HTTP request and can make a decision by returning a `Response` object (or nothing if no decision is made).

## CEL language

An `AuthorizationPolicy` uses the [CEL language](https://github.com/google/cel-spec) to process the incoming HTTP request.

CEL is an expression language that’s fast, portable, and safe to execute in performance-critical applications.

## Policy structure

A `AuthorizationPolicy` is made of:

- A [failure policy](./failure-policy.md)
- Eventually some [variables](./variables.md)
- The [authorization rules](./authorization-rules.md)
