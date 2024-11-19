#!/bin/bash

ROOT=$(dirname "${BASH_SOURCE[0]}")/..
NAMESPACE=${NAMESPACE:-default}

helm upgrade --install --repo https://jont828.github.io/cluster-api-visualizer/charts cluster-api-visualizer cluster-api-visualizer -n ${NAMESPACE} || exit 1
kubectl rollout status deployment -n ${NAMESPACE} capi-visualizer

./${ROOT}/hack/port-forward.sh