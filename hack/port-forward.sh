#!/bin/bash

NAMESPACE=${NAMESPACE:-default}

echo "Running at http://localhost:8081"
kubectl port-forward -n ${NAMESPACE} service/capi-visualizer 8081:8081