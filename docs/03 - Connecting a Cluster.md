# Demo Scenario - Connecting a Cluster

To connect a cluster we need to complete a few steps:

- Setup a new cluster on your preffered k8s distribution
- Deploy the application
- Connect the cluster on Conductor

## Setup a new cluster on your prefered k8s distribution

The first step is to setup a new cluster to deploy the Order System app. This guide in the Dapr docs explains how to do it for multiple distributions: [How-to: Setup clusters](https://docs.dapr.io/operations/hosting/kubernetes/cluster/)

## Deploy the application

### Clone this repository

```git
git clone https://github.com/diagridio/conductor-demo-new.git

cd conductor-demo-new
```
### Create application namespaces

After connecting to your cluster, run the following command to create the namespaces:

```bash
# create application namespaces
kubectl create ns order-system # Our services
kubectl create ns dapr-system # Dapr
kubectl create ns redis # Redis
kubectl create ns kafka # Kafka
kubectl create ns zipkin # Zipkin
```

### Install Dapr and create RBAC

You can choose to skip this step if you want Conductor to manage your Dapr installation (recommended):

```bash
dapr init -k -n dapr-system
kubectl apply -f ./components/k8s/dapr-secret-reader.yaml
```

### Setup metrics server

If you are on **Conductor Free**, this step will be taken care automatically for you. If you are on **Conductor Enterprise**, this step is needed for some distributions. See more details at [Conductor pre-requisites](https://docs.diagrid.io/conductor/getting-started/prereqs/#installation-prerequisites).

```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
kubectl patch deployment metrics-server -n kube-system --type "json" -p '[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'
```

### Setup Helm

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

### Redis setup

Install Redis, export the password and create a secret that will be accessed from the component files.

```bash
helm install redis bitnami/redis -n redis
export REDIS_PASSWORD=$(kubectl get secret --namespace redis redis -o jsonpath="{.data.redis-password}" | base64 -d) 

kubectl create secret generic redis-password --from-literal=redis-password=$REDIS_PASSWORD -n order-system
```

### Kafka setup

Install Kafka, export the password and create a secret that will be accessed from the component files.

```bash
helm install --set persistence.enabled=false --set zookeeper.persistence.enabled=false --set auth.clientProtocol=sasl kafka bitnami/kafka -n kafka

export KAFKA_PASSWORD=$(kubectl get secret kafka-user-passwords --namespace kafka -o jsonpath='{.data.client-passwords}' | base64 -d | cut -d , -f 1)
kubectl create secret generic kafka-password --from-literal=kafka-password=$KAFKA_PASSWORD -n order-system
```

> [!IMPORTANT]
> Since we are inducing a component security advisory, update the content of `/components/k8s/oms.pubsub.yaml` with the new value for $KAFKA_PASSWORD.
>
> Run `echo $KAFKA_PASSWORD` to retrieve the value.

### Zipkin setup

Zipkin will be used for applicaiton tracing.

```bash
echo "Zipkin namespace created"
kubectl create deployment zipkin -n zipkin --image openzipkin/zipkin
echo "Zipkin deployment created"
kubectl expose deployment zipkin -n zipkin --type LoadBalancer --port 9411 
echo "Zipkin deployment exposed"
kubectl get svc -n zipkin -w
export ZIPKIN_DASHBOARD=$(kubectl get svc --namespace zipkin zipkin -o jsonpath="{.status.loadBalancer.ingress[0].ip}"):9411
echo "View tracing dashboard at $ZIPKIN_DASHBOARD"
```

### Deploy Dapr components

Now we will deploy the components that we will use throughout the demo:

```bash
kubectl apply -f ./components/k8s
```

### Deploy services

To deploy the services locally, navigate to the [order-system public registry](https://console.cloud.google.com/artifacts/docker/prj-common-d-shared-89549/us-central1/reg-d-common-docker-public), fetch the latest tag for each container starting with `order-system/` and update the image tag at each corresponding file in the `/deployment-files/k8s-dev/` folder.

Finally, run:

```bash
kubectl apply -f ./deployment-files/k8s-dev
```

> The contents of the folder `/deployment-files/k8s` are used for automated CI/CD pipelines, as described in the main [README](./../README.md) file.

Check the state of your pods by running:

```bash
kubectl get pods -n order-system
```

### Apply Dapr components

All components are in the folder `components/k8s`. During the demo, you will need to apply new configurations. To do so, modify one or more yaml files and run:

```bash
kubectl apply -f components/k8s
```

## Connect the cluster on Conductor

Follow [Cluster onboarding](https://docs.diagrid.io/conductor/getting-started/cluster-onboarding) to connect your cluster.



