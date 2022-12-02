#! /usr/bin/env bash

set -euo pipefail

drone secret rm --name KO_DOCKER_REPO harness-apps/go-fruits-api

drone secret rm --name IMAGE_REGISTRY_USER harness-apps/go-fruits-api

drone secret rm --name IMAGE_REGISTRY_PASSWORD harness-apps/go-fruits-api