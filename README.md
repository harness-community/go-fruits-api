# Fruits API

The demo application that can be used to demonstrate on how to do CI with Java Applications with Drone CI.

## Pre-requisites

* [Docker Desktop](https://docs.docker.com/desktop/)
* [kubectl](https://kubernetes.io/docs/tasks/tools)
* [httpie](https://httpie.io)
* [k3d](https://k3d.io)
* [kustomize](https://kustomize.io/)
  
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

## Test

### Local

```shell
mkdir work
docker-compose up -d
go run server.go -dbType pg 
```

Swagger UI <http://localhost:8080/swagger/index.html>

### Kubernetes

### Create Cluster

Lets keep kube config local to the folder,

```shell
mkdir -p "$PWD/.kube"
```

Create a local container registry to use,

```shell
k3d registry create myregistry.localhost --port 5001
```

Use the registry with the cluster,

```shell
k3d cluster create fruits-api \
  --registry-use k3d-myregistry.localhost:5001
```

## Deploy Database

```shell
kubectl apply -k k8s/db
```

Wait for pods to be in running state:

```shell
kubectl rollout status -n db deploy/postgresql --timeout=60s
```

## Deploy application

```shell
kubectl apply -f k8s/app
```

Wait for pods to be in running state:

```shell
kubectl rollout status -n fruits-app deploy/fruits-api --timeout=60s
```

```shell
kubectl port-forward -n fruits-app svc/fruits-api 8080:8080

```

Check if you are able to access the API, the following call should return you list of fruits as JSON,

```shell
curl localhost:8080/v1/api/fruits
```

```json
[
  { "id": 1, "name": "Mango", "season": "Spring", "emoji": "U+1F96D" },
  { "id": 2, "name": "Strawberry", "season": "Spring", "emoji": "U+1F96D" },
  { "id": 3, "name": "Orange", "season": "Winter", "emoji": "U+1F34B" },
  { "id": 4, "name": "Lemon", "season": "Winter", "emoji": "U+1F34A" },
  { "id": 5, "name": "Blueberry", "season": "Summer", "emoji": "U+1FAD0" },
  { "id": 6, "name": "Banana", "season": "Summer", "emoji": "U+1F34C" },
  { "id": 7, "name": "Watermelon", "season": "Summer", "emoji": "U+1F349" },
  { "id": 8, "name": "Apple", "season": "Fall", "emoji": "U+1F34E" },
  { "id": 9, "name": "Pear", "season": "Fall", "emoji": "U+1F350" }
]
```

## Local Development

Download/install [ko](https://github.com/ko-build/ko) and run,

```shell
kustomize build config/app | ko resolve -f - | k apply -f -
```

> **NOTE:**
> If you are using arm64 machines then add option --platform=linux/arm64 to ko command
>
> ```shell
> kustomize build config/app | ko resolve --platform=linux/arm64 -f - | k apply -f -
> ```

## Clean

```shell
k3d cluster delete fruits-api
```
