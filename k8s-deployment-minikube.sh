# start minikube
minikube start --cpus=4 --memory=4096
 
 # enable dashboard and ingress
minikube addons enable dashboard
minikube addons enable ingress



# create application namespaces
kubectl create ns oms
kubectl create ns dapr 
kubectl create ns redis
kubectl create ns kafka

# install dapr
dapr init -k -n dapr

kubectl apply -f ./components/k8s/rbac/dapr-secret-reader.yaml

# setting up metrics server - Conductor pre-requisites (https://docs.diagrid.io/conductor/getting-started/prereqs/#installation-prerequisites)
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
kubectl patch deployment metrics-server -n kube-system --type "json" -p '[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'

# setup help
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# setup redis
helm install redis bitnami/redis -n redis
export REDIS_PASSWORD=$(kubectl get secret --namespace redis redis -o jsonpath="{.data.redis-password}" | base64 -d) 

kubectl create secret generic redis-password --from-literal=redis-password=$REDIS_PASSWORD -n oms

##### MAYBE #####
# Install Zipkin for tracing - non-minikube
kubectl create ns zipkin
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

# TODO: DEPLOY APPS


#To connect to your database from outside the cluster execute the following commands:

 #   kubectl port-forward --namespace redis svc/redis-master 6379:6379 &
  #  REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli -h 127.0.0.1 -p 6379




gcloud compute firewall-rules list --filter="name~gke-${gke_prj-dataplane-n-demo-30534_us-west1-a_gke-n-dataplane-demo-us-west1-a2}-[0-9a-z]*-master"

# MAKE SURE TO HAVE LOCAL ENVIRONEMNT SET TO BUILD CONTAINERS FOR MULTIPLE PLATFORMS