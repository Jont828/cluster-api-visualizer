#!/bin/bash

ROOT=$(dirname "${BASH_SOURCE[0]}")/..
NAMESPACE=${NAMESPACE:-default}

helm install --repo https://jont828.github.io/cluster-api-visualizer/charts cluster-api-visualizer --generate-name -n ${NAMESPACE}|| exit 1
kubectl rollout status deployment -n ${NAMESPACE} capi-visualizer

echo "Running at http://localhost:8081"
kubectl port-forward -n ${NAMESPACE} service/capi-visualizer 8081:8081