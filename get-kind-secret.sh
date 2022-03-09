kind get kubeconfig --name capi-test --internal > kind-kc
kubectl create secret generic kind-kubeconfig --from-file=kind-kc