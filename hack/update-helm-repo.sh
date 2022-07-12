#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

ROOT=$(dirname "${BASH_SOURCE[0]}")/..

helm package ${ROOT}/helm/cluster-api-visualizer -d ${ROOT}/helm/repo/new
helm repo index ${ROOT}/helm/repo/new
helm repo index --merge ${ROOT}/helm/repo/index.yaml ${ROOT}/helm/repo/new
mv ${ROOT}/helm/repo/new/index.yaml ${ROOT}/helm/repo/index.yaml
mv ${ROOT}/helm/repo/new/*.tgz ${ROOT}/helm/repo