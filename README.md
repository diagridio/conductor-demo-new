## Conductor Order Management System

A sample order management system composed of 8 Dapr-enabled microservices to showcase Conductor features.

### Architecture and Service Definitions

![Logical Application Architecture Diagram](assets/demo-arch.png)


| Service          | Definition                                                                                                 | Dapr Component(s) |
|------------------|-------------------------------------------------------------------------------------------------------------| --------------------------------|
| Loyalty Service | Manages the loyalty program by modifying customer reward points based on spend | Kafka subscriber, Redis state store |
| Make Line Service | Responsible for simulating and coordinating a 'queue' of current orders. Monitors the processing and completion of each order in the 'queue' | Kafka subscriber, Redis state store |  
| Order Service | Basic CRUD API that is used to place and manage orders | Kafka publisher |
| Receipt Generation Service | Archival program that generates and stores order receipts for auditing and historical purposes  | Kafka subscriber, Redis output binding |
| Order Processor Workflow  | Responsible for orchestrating a workflow(Loyalty, Receipt, Make-Line) whenever an order is created.   | Kafka subscriber, Redis output binding, Dapr Workflow Client |
| Virtual Customer | 'Customer simulation' program that simulates customers placing orders | Order service invocation |
| Virtual Worker | 'Worker simulation' program that simulates the completion of customer orders | Cron input binding, Make-line service invocation | 

### Prerequisites

- Kubernetes cluster of your choice (3+ nodes recommended)
- Helm

### Deployment

Run the commands in the [./conductor-setup](./conductor-setup.sh) script one by one for best results. All Dapr components are deployed on Kubernetes today.

#### Manual deployments to k8s during development

For manual deployments, use the deployment files in `deployment-files/k8s-dev`. The folder `deployment-files/k8s` is used exclusively by deployments triggered by pushes to the main branch.

#### GitHub Actions

Every push to the main branch will trigger a process that builds and publishes a new image tagged with the commit SHA. The image is deployed to GKE after that. [Kustomize](https://kustomize.io/) is used to create the deployments dynamically.

#### Important

Since we are inducing a compinent security advisory, update the content of _oms.pubsub.yaml_ with the new value for $KAFKA_PASSWORD. Redeploy the component.

### Build and push docker images to the registry individually.

Inside each service folder there is a Makefile. Navigate to the folder and run the command below to build and push the images. 

  ```bash
  cd go/order-service
  make build_image
  ```
#### About container architecture

The containers are currently built for both ARM and AMD64 architectures. 

### Running locally

#### To test the services locally, you can run the following command: 

```bash
cd go/order-service
make run
``` 
