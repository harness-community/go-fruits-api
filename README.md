# Fruits API

The demo application that can be used to demonstrate on how to do CI with Java Applications with Drone CI.

## Pre-requisites

* [Docker Desktop](https://docs.docker.com/desktop/)
* [kubectl](https://kubernetes.io/docs/tasks/tools)
* [httpie](https://httpie.io)
* [k3d](https://k3d.io)
  
## Download Sources

This application has two microservices

1. [Fruits API](https://github.com/kameshsampath/go-fruits-api)  -  that provides the REST API for managing the Fruits

2. [Fruits App UI](https://github.com/kameshsampath/fruits-app-ui) - the front end to the application

```shell
git clone https://github.com/harness-apps/go-fruits-api go-fruits-api && cd $_
```

We will refer to the cloned folder as `$PROJECT_HOME`:

```shell
export PROJECT_HOME="$(pwd)"
```

## Deploy Database

```shell
k apply -k config/db
```

Wait for pods to be in running state:

```shell
kubectl rollout status -n db deploy/postgresql --timeout=60s
```

## Deploy application

```shell
kustomize build config/app | ko resolve -f - | k apply -f -
```

Wait for pods to be in running state:

```shell
kubectl rollout status -n fruits-app deploy/fruits-app-api --timeout=60s
```

## Test

### Local

Swagger UI <http://localhost:8080/swagger/index.html>

### Kubernetes

```shell
kubectl port-forward 
```

Check if you are able to access the API, the following call should return you list of fruits as JSON,

```shell
http localhost:8080/v1/api/fruits
```
