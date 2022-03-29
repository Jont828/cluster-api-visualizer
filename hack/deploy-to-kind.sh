#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <name-of-kind-management-cluster>"
  exit 1
fi

KUBECONFIG_DATA=$(kind get kubeconfig --name ${1} --internal)

helm install --generate-name ./helm/capi-visualization --set kubeconfig="$KUBECONFIG_DATA"  || exit 1
kubectl rollout status deployment visualize-cluster

echo "Running at http://localhost:8081"
kubectl port-forward service/visualize-cluster 8081:8081