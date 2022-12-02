#! /usr/bin/env bash

set -euo pipefail

drone secret add --name ko_docker_repo --data "${KO_DOCKER_REPO}" harness-apps/go-fruits-api

drone secret add --name image_registry_user --data "${IMAGE_REGISTRY_USER}" harness-apps/go-fruits-api

drone secret add --name image_registry_password --data "${IMAGE_REGISTRY_PASSWORD}" harness-apps/go-fruits-api
