#! /usr/bin/env bash

set -euo pipefail

drone secret update --name KO_DOCKER_REPO --data "${KO_DOCKER_REPO}" harness-apps/go-fruits-api

drone secret update --name IMAGE_REGISTRY_USER --data "${IMAGE_REGISTRY_USER}" harness-apps/go-fruits-api

drone secret update --name IMAGE_REGISTRY_PASSWORD --data "${IMAGE_REGISTRY_PASSWORD}" harness-apps/go-fruits-api