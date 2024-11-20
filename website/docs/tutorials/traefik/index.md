# Traefik 

[Traefik](https://doc.traefik.io/traefik/) is an open-source Application Proxy that makes publishing your services a fun and easy experience. It receives requests on behalf of your system and identifies which components are responsible for handling them, and routes them securely.

This tutorial shows how Traefik's [ForwardAuth middleware](https://doc.traefik.io/traefik/middlewares/http/forwardauth/) can be configured to delegate authorization decisions to the Generic Auth Server.

## Setup

### Prerequisites

- A Kubernetes cluster
- [Helm](https://helm.sh/) to install Traefik and the Generic Auth Server
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) to interact with the cluster

### Setup a cluster

Create a local cluster with [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation).

```bash
KIND_IMAGE=kindest/node:v1.31.1

# create cluster
kind create cluster --image $KIND_IMAGE --wait 1m
```

### Install Traefik

Install Traefik on the cluster.

```bash
# install traefik
helm install traefik --namespace traefik --create-namespace \
  --wait --repo https://traefik.github.io/charts traefik \
  --values - <<EOF
service:
  type: ClusterIP
tolerations:
  - key: node-role.kubernetes.io/control-plane
    operator: Equal
    effect: NoSchedule
EOF
```

### Deploy the Generic Auth Server

Now deploy the Generic Auth Server.

```bash
# deploy the generic auth server
helm install generic-auth-server --namespace kyverno --create-namespace \
  --wait --repo https://eddycharly.github.io/generic-auth-server \
  generic-auth-server
```

### Configure the ForwardAuth middleware

The [ForwardAuth middleware](https://doc.traefik.io/traefik/middlewares/http/forwardauth/) will delegate authorization decisions to the Generic Auth Server.

```bash
# configure the forward auth middleware
kubectl apply -n traefik -f - <<EOF
apiVersion: traefik.io/v1alpha1
kind: Middleware
metadata:
  name: auth
spec:
  forwardAuth:
    address: http://generic-auth-server.kyverno.svc.cluster.local:9081/auth
EOF
```

Notice that the middleware uses the Generic Auth Server url we deployed earlier:

```yaml
[...]
  forwardAuth:
    address: http://generic-auth-server.kyverno.svc.cluster.local:9081/auth
[...]
```

### Deploy a sample application

We will use the `whoami` sample application that can be used to test HTTP requests and helps to show quickly how we can play with the request and response attributes.

```bash
# create the demo namespace
kubectl create ns demo

# deploy the httpbin application
kubectl apply -n demo -f - <<EOF
kind: Deployment
apiVersion: apps/v1
metadata:
  name: whoami
  labels:
    app: whoami
spec:
  replicas: 1
  selector:
    matchLabels:
      app: whoami
  template:
    metadata:
      labels:
        app: whoami
    spec:
      containers:
        - name: whoami
          image: traefik/whoami
          ports:
            - name: web
              containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: whoami
spec:
  ports:
    - name: web
      port: 80
      targetPort: web

  selector:
    app: whoami
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: whoami-ingress
  annotations:
    traefik.ingress.kubernetes.io/router.middlewares: traefik-auth@kubernetescrd
spec:
  rules:
  - host: foo.bar.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: whoami
            port:
              name: web
EOF
```

Notice that the ingress references our middleware using annotations:

```yaml
[...]
  annotations:
    traefik.ingress.kubernetes.io/router.middlewares: traefik-auth@kubernetescrd
[...]
```

## Create an AuthorizationPolicy

In summary the policy below does the following:

- Checks if the incoming request contains the header `x-force-authorized` with the value `enabled` or `true`
- Allows the request if it has the header or denies it if not

```yaml
# deploy authorization policy
kubectl apply -f - <<EOF
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
EOF
```

## Testing

At this point we have deployed and configured Traefik, the Generic Auth Server, a sample application, and an authorization policies.

### Start an in-cluster shell

Let's start a pod in the cluster with a shell to call into the sample application.

```bash
# run an in-cluster shell
kubectl run -i -t busybox --image=alpine --restart=Never -n demo
```

### Install curl

We will use curl to call into the sample application but it's not installed in our shell, let's install it in the pod.

```bash
# install curl
apk add curl
```

### Call into the sample application

Now we can send requests to the sample application and verify the result.

The following request will return `401` (denied by our policy):

```bash
curl -s -w "\nhttp_code=%{http_code}" http://traefik.traefik.svc.cluster.local \
  -H "Host: foo.bar.com"
```

The following request will return `200` (allowed by our policy):

```bash
curl -s -w "\nhttp_code=%{http_code}" http://traefik.traefik.svc.cluster.local \
  -H "Host: foo.bar.com" \
  -H "x-force-authorized: true"
```

## Wrap Up

Congratulations on completing the tutorial!

This tutorial demonstrated how to configure Traefik’s ForwardAuth middleware to utilize the Generic Auth Server as an external authorization service.
