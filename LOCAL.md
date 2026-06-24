# Run the demo in your own cluster

## Prerequisites

- Kubernetes cluster of your choice (3+ nodes recommended)
- Helm

## Clone this repository

```git
git clone https://github.com/diagridio/conductor-demo-new.git

cd conductor-demo-new
```

## Create application namespaces

After connecting to your cluster, run the following command to create the namespaces:

```bash
# create application namespaces
kubectl create ns order-system # Our services
kubectl create ns dapr-system # Dapr
kubectl create ns redis # Redis
kubectl create ns kafka # Kafka
kubectl create ns zipkin # Zipkin
```

## Install Dapr and create RBAC

You can choose to skip this step if you want Conductor to manage your Dapr installation (recommended):

```bash
dapr init -k -n dapr-system
kubectl apply -f ./components/k8s/dapr-secret-reader.yaml
```

## Setup metrics server

If you are on **Conductor Free**, this step will be taken care automatically for you. If you are on **Conductor Enterprise**, this step is needed for some distributions. See more details at [Conductor pre-requisites](https://docs.diagrid.io/conductor/getting-started/prereqs/#installation-prerequisites).

```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
kubectl patch deployment metrics-server -n kube-system --type "json" -p '[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'
```

## Setup Helm

```bash
helm repo add valkey https://valkey.io/valkey-helm/
helm repo update
```

## Redis setup

Install Valkey (Redis-compatible).

```bash
helm install valkey valkey/valkey -n redis
```

## Kafka setup

Install Strimzi and create the single-node Kafka cluster used by the demo components.

```bash
# Install the operator
helm install strimzi-kafka-operator oci://quay.io/strimzi-helm/strimzi-kafka-operator -n kafka

# Create a Kafka cluster named "my-cluster"
kubectl apply -f https://strimzi.io/examples/latest/kafka/kafka-single-node.yaml -n kafka
```

## Install Zipkin

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

## Optional: Jaeger/OpenTelemetry backend

If you skip this section, services still run, but logs may include periodic trace export timeout warnings.

```bash
kubectl create ns cert-manager
kubectl create ns observability

kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.16.1/cert-manager.yaml -n cert-manager
kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/latest/download/opentelemetry-operator.yaml
kubectl apply -f ./deployment-files/k8s/jaeger.yaml -n observability
```

## Deploy Dapr components

Now we will deploy the components that we will use throughout the demo:

```bash
kubectl apply -f ./components/k8s
```

## Deploy services

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

## Known issues

There is a current issue with Redis where constantly becomes full. To prevent that we have introduced a cronjob called`redis-full-wipe-cronjob` which flushes redis every 24hrs at 05:00 UTC.

```sh
# View recent executions
kubectl get jobs -n redis

# Check latest logs
kubectl logs -n redis jobs/redis-full-wipe-cronjob-xxxx
```
**Sample execution**
```
Starting Redis wipe at Tue Jul 15 05:00:01 UTC 2025 for valkey.redis.svc.cluster.local:6379
Keys before wipe: 243906
Executing FLUSHALL command...
OK
FLUSHALL completed successfully
Keys after wipe: 0
```

```bash
VALKEY_POD=$(kubectl get pods -n redis -l app.kubernetes.io/name=valkey -o jsonpath='{.items[0].metadata.name}')
kubectl exec -it "$VALKEY_POD" -n redis -- /bin/sh
```

When the Redis prompt is available, run:

```bash
valkey-cli --scan --pattern '*' | xargs valkey-cli DEL
```

Finally, run `exit` to leave the Redis prompt.

