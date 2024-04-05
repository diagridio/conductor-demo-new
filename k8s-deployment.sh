# create application namespace
kubectl create ns oms
kubectl apply -f ./components/k8s/rbac/dapr-secret-reader.yaml

# setup help
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# setup redis
kubectl create ns redis
helm install redis bitnami/redis -n redis
export REDIS_PASSWORD=$(kubectl get secret --namespace redis redis -o jsonpath="{.data.redis-password}" | base64 -d) 
kubectl create secret generic redis-password --from-literal=redis-password=$REDIS_PASSWORD -n oms

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
kubectl apply -f components/k8s/oms.config.yaml

# Deploy Dapr components
kubectl apply -f components/k8s/

# TODO: DEPLOY APPS



