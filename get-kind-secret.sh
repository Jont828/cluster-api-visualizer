kind get kubeconfig --name capz --internal > kind-kc
kubectl create secret generic kind-kubeconfig --from-file=kind-kc