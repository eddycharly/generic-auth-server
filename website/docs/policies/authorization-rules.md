# Authorization rules

An `AuthorizationPolicy` main element is the authorization rules defined in `authorizations`.

Every authorization rule must contain a [CEL](https://github.com/google/cel-spec) `expression`. It is expected to return a `Response` describing the decision made by the rule (or nothing if no decision is made).

Creating the `Response` can be a tedious task, you need to remember the different types names and format.

The CEL engine used to evaluate the authorization rules has been extended with a library to make the creation of `Response` easier. Browse the [available libraries documentation](../cel-extensions/index.md) for details.

## Authorization rules

TODO