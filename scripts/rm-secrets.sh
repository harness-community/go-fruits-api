#! /usr/bin/env bash

set -euo pipefail

drone secret rm --name ko_docker_repo harness-apps/go-fruits-api

drone secret rm --name image_registry_user harness-apps/go-fruits-api

drone secret rm --name image_registry_password harness-apps/go-fruits-api