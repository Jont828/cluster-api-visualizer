#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <name-of-kind-management-cluster>"
  echo "Note: existing kind clusters are [ $(kind get clusters) ]"
  exit 1
fi

ROOT=$(dirname "${BASH_SOURCE[0]}")/..

KUBECONFIG_DATA=$(kind get kubeconfig --name ${1} --internal)

helm install --generate-name ${ROOT}/helm/capi-visualization --set kubeconfig="$KUBECONFIG_DATA" || exit 1
kubectl rollout status deployment capi-visualizer

echo "Running at http://localhost:8081"
kubectl port-forward service/capi-visualizer 8081:8081