kc expose deployment test-deployment --port=8081 --target-port=8081 --name=test-service
kc port-forward service/test-service 8081:8081