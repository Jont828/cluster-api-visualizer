#!/bin/bash

ROOT=$(dirname "${BASH_SOURCE[0]}")/..
NAMESPACE=${NAMESPACE:-default}

helm install -name cluster-api-visualizer ${ROOT}/helm/cluster-api-visualizer -n ${NAMESPACE} || exit 1
kubectl rollout status deployment -n ${NAMESPACE} capi-visualizer

echo "Running at http://localhost:8081"
kubectl port-forward -n ${NAMESPACE} service/capi-visualizer 8081:8081