#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <name-of-kind-management-cluster>"
  exit 1
fi

kind get kubeconfig --name "$1" --internal > kind-kc
kubectl delete secret kind-kubeconfig --ignore-not-found
kubectl create secret generic kind-kubeconfig --from-file=kind-kc
rm kind-kc

kubectl apply -f ./hack/deployments/visualize.yaml
kubectl rollout status deployment visualize-cluster

echo "Running at http://localhost:8081"
kubectl port-forward service/visualize-cluster 8081:8081