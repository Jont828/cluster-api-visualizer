#!/bin/bash

ROOT=$(dirname "${BASH_SOURCE[0]}")/..
NAMESPACE=${NAMESPACE:-default}

helm install -name cluster-api-visualizer ${ROOT}/helm/cluster-api-visualizer -n ${NAMESPACE} || exit 1
kubectl rollout status deployment -n ${NAMESPACE} capi-visualizer

./${ROOT}/hack/port-forward.sh