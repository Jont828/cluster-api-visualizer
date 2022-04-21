#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <name-of-kind-management-cluster>"
  echo "Note: existing kind clusters are [ $(kind get clusters) ]"
  exit 1
fi

ROOT=$(dirname "${BASH_SOURCE[0]}")/..
NAMESPACE=${NAMESPACE:-default}

KUBECONFIG_DATA=$(kind get kubeconfig --name ${1} --internal)

kubectl delete secret -n ${NAMESPACE} management-kubeconfig --ignore-not-found
helm install --repo https://raw.githubusercontent.com/Jont828/cluster-api-visualizer/main/helm/repo cluster-api-visualizer --generate-name --set kubeconfig="$KUBECONFIG_DATA" -n ${NAMESPACE}|| exit 1
kubectl rollout status deployment -n ${NAMESPACE} capi-visualizer

echo "Running at http://localhost:8081"
kubectl port-forward -n ${NAMESPACE} service/capi-visualizer 8081:8081