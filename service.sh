kc expose deployment visualize-cluster --port=8081 --target-port=8081 --name=visualize-cluster-service
kc port-forward service/visualize-cluster-service 8081:8081