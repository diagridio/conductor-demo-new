
# create application namespaces
kubectl create ns order-system # Our services
kubectl create ns dapr-system # Dapr
kubectl create ns redis # Redis
kubectl create ns kafka # Kafka
kubectl create ns zipkin # Zipkin

# install dapr
dapr init -k -n dapr-system

kubectl apply -f ./components/k8s/dapr-secret-reader.yaml

# setting up metrics server - Conductor pre-requisites (https://docs.diagrid.io/conductor/getting-started/prereqs/#installation-prerequisites)
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
kubectl patch deployment metrics-server -n kube-system --type "json" -p '[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'

# setup help
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# setup redis
helm install redis bitnami/redis -n redis
export REDIS_PASSWORD=$(kubectl get secret --namespace redis redis -o jsonpath="{.data.redis-password}" | base64 -d) 

kubectl create secret generic redis-password --from-literal=redis-password=$REDIS_PASSWORD -n order-system

# Install Kafka for pub/sub
helm install --set persistence.enabled=false --set zookeeper.persistence.enabled=false --set auth.clientProtocol=sasl kafka bitnami/kafka -n kafka

# Modifying the kafka-jaas secret to get the password, as the kafka-jaas secret is not available in the kafka namespace
#export KAFKA_PASSWORD=$(kubectl get secret kafka-jaas --namespace kafka -o jsonpath='{.data.client-passwords}' | base64 -d | cut -d , -f 1)
export KAFKA_PASSWORD=$(kubectl get secret kafka-user-passwords --namespace kafka -o jsonpath='{.data.client-passwords}' | base64 -d | cut -d , -f 1)
kubectl create secret generic kafka-password --from-literal=kafka-password=$KAFKA_PASSWORD -n order-system

# Install Zipkin for tracing
echo "Zipkin namespace created"
kubectl create deployment zipkin -n zipkin --image openzipkin/zipkin
echo "Zipkin deployment created"
kubectl expose deployment zipkin -n zipkin --type LoadBalancer --port 9411 
echo "Zipkin deployment exposed"
kubectl get svc -n zipkin -w
export ZIPKIN_DASHBOARD=$(kubectl get svc --namespace zipkin zipkin -o jsonpath="{.status.loadBalancer.ingress[0].ip}"):9411
echo "View tracing dashboard at $ZIPKIN_DASHBOARD"

# apply zipkin configuration
kubectl apply -f ./components/k8s/oms.config.yaml
##### MAYBE #####

# Deploy Dapr components
kubectl apply -f ./components/k8s

#Deploy services
kubectl apply -f ./deployment-files/k8s


#To connect to your database from outside the cluster execute the following commands:

  #  kubectl port-forward --namespace redis svc/redis-master 6379:6379 &
  #  REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli -h 127.0.0.1 -p 6379

gcloud compute firewall-rules list --filter="name~gke-${gke_prj-dataplane-n-demo-30534_us-west1-a_gke-n-dataplane-demo-us-west1-a2}-[0-9a-z]*-master"

# The containers are currently built for both ARM and AMD64 architectures. Every service folder has a Makefile that can be used to build the container images. 
# The build and push targets can be run with the following command: 
  #   make build_image
# To test the services locally, you can run the following command: 
  #   make run


# Flushall keys in redis
k exec -it redis-master-0 -n redis -- /bin/bash

redis-cli -a $REDIS_PASSWORD --scan --pattern '*' | xargs redis-cli -a $REDIS_PASSWORD DEL
