# Conductor Order Management System

A sample order management system composed of 8 Dapr-enabled microservices to showcase Conductor features.

## Architecture and Service Definitions

![Logical Application Architecture Diagram](assets/demo-arch.png)


| Service          | Definition                                                                                                 | Dapr Component(s) |
|------------------|-------------------------------------------------------------------------------------------------------------| --------------------------------|
| Loyalty Service | Manages the loyalty program by modifying customer reward points based on spend | Kafka subscriber, Redis state store |
| Make Line Service | Responsible for simulating and coordinating a 'queue' of current orders. Monitors the processing and completion of each order in the 'queue' | Kafka subscriber, Redis state store |  
| Order Service | Basic CRUD API that is used to place and manage orders | Kafka publisher |
| Receipt Generation Service | Archival program that generates and stores order receipts for auditing and historical purposes  | Kafka subscriber, Redis output binding |
| Order Processor Workflow  | Responsible for orchestrating a workflow(Loyalty, Receipt, Make-Line) whenever an order is created.   | Kafka subscriber, Redis output binding, Dapr Workflow Client |
| Virtual Customer | 'Customer simulation' program that simulates customers placing orders | Order service invocation |
| Worker Actor | 'Worker simulation' program that simulates the completion of customer orders | Actor, Redis State |

## Demo setup

See [Demo overview](./docs/01%20-%20Overview.md).

> If you plan to deploy the services locally in your own cluster, follow the instructions in [LOCAL.md](./LOCAL.md). Skip the rest of the instructions below.

If you are a developer working on improving these services, keep following the instructions here.

## Development instructions


### Clone this repository

```git
git clone https://github.com/diagridio/conductor-demo-new.git

cd conductor-demo-new
```

### Choose your path

There are two development processes you can follow: 

 - [Development and deployment to local k8s cluster.](https://github.com/diagridio/conductor-demo-new/tree/demo-scenarios?tab=readme-ov-file#deployment-to-local-cluster-kindminikubeetc) 
 - [Development and run the services individually in your local machine (non-containerized)](https://github.com/diagridio/conductor-demo-new/tree/demo-scenarios?tab=readme-ov-file#running-dapr-services-locally-non-containerized).

## Deployment to local cluster (Kind/Minikube/etc)

### Prerequisites

- Kubernetes cluster of your choice (3+ nodes recommended)
- Helm

### Create application namespaces

After connecting to your cluster, run the following command to create the namespaces we will use for the app, Dapr, and supporting components:

```bash
kubectl create ns order-system 
kubectl create ns dapr-system 
kubectl create ns redis 
kubectl create ns kafka 
kubectl create ns zipkin 
```

### Install Dapr and create RBAC

```bash
dapr init -k -n dapr-system
kubectl apply -f ./components/k8s/dapr-secret-reader.yaml
```

### Setup metrics server

The metrics server is a Conductor pre-requisite. Some k8s distributions already cover this requirement.

See more details at [Conductor pre-requisites](https://docs.diagrid.io/conductor/getting-started/prereqs/#installation-prerequisites).

```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
kubectl patch deployment metrics-server -n kube-system --type "json" -p '[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'
```

### Redis(Valkey) setup

Install Valkey

```bash
helm repo add valkey https://valkey.io/valkey-helm/ 

helm repo update

helm install valkey valkey/valkey --set replicaCount=3 -n redis
```

### Kafka setup

Install Kafka

```bash
# set up the strimzi operator
helm install strimzi-kafka-operator oci://quay.io/strimzi-helm/strimzi-kafka-operator

# Add CRD to spin up a single node kafka
kubectl apply -f https://strimzi.io/examples/latest/kafka/kafka-single-node.yaml -n kafka 
```

### Install Jaeger

```bash
#install cert manager 
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.16.1/cert-manager.yaml -n cert-manager

# install OpenTelemetry Operator
kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/latest/download/opentelemetry-operator.yaml

# install instance of jaeger 
kubectl apply -f ./deployment-files/k8s/jaeger.yaml -n observability
```

### Install Zipkin

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

### Manual deployments to k8s during development

For manual deployments, use the deployment files in `deployment-files/k8s-dev`. Modify the images within the deployment files with your own IMAGE:TAG.

> The folder `deployment-files/k8s` is used exclusively by deployments triggered by pushes to the main branch.

#### Build and push docker images to the development registry individually.

Inside each service folder there is a Makefile. Navigate to the folder, modify the contents of the Makefile to deploy the images to your own registry and run the command below to build and push the images. 

For example, to build and push the `order-service` image, run:

  ```bash
  cd go/order-service
  make build_image
  ```

## Running Dapr services locally (non-containerized)

Running the services in your computer will be simpler, as we won't depend on k8s, Helm, and Kafka.

### Pre-requisites

- Go version 1.20
- .NET Core 8
- Docker desktop
- Dapr 1.13.2

### Deploy Dapr components

First we will deploy the components that we will use throughout the demo:

```bash
kubectl apply -f ./components/local
```

### To test the services locally, you can leverage Makefile.

Inside each service folder there is a Makefile. Navigate to the folder. run the command below to tun the service locally. 

For example, to run the `order-service` app, run:

```bash
cd go/order-service
make run
``` 

## CI/CD Information

## GitHub Actions

Every push to the main branch will trigger a process that builds and publishes a new image tagged with the commit SHA. The image is deployed to two GKE clusters after that. For details on the clustes, see [Demo overview](./docs/01%20-%20Overview.md).

[Kustomize](https://kustomize.io/) is used to create the deployments dynamically.

## About container architecture

The containers are currently built for both ARM and AMD64 architectures. 

