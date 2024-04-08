# Readme## Conductor Order Management System

A sample order management system composed of 8 Dapr-enabled microservices to showcase Conductor features.

### Architecture and Service Definitions

![Logical Application Architecture Diagram](assets/demo-arch.png)


| Service          | Definition                                                                                                 | Dapr Component(s) |
|------------------|-------------------------------------------------------------------------------------------------------------| --------------------------------|
| Loyalty Service | Manages the loyalty program by modifying customer reward points based on spend | Kafka subscriber, Redis state store |
| Make Line Service | Responsible for simulating and coordinating a 'queue' of current orders. Monitors the processing and completion of each order in the 'queue' | Kafka subscriber, Redis state store |  
| Order Service | Basic CRUD API that is used to place and manage orders | Kafka publisher |
| Receipt Generation Service | Archival program that generates and stores order receipts for auditing and historical purposes  | Kafka subscriber, Redis output binding |
| Order Generator Service | Responsible for orchestrating a workflow whenever an order is created.  | Kafka subscriber, Redis output binding |
| Virtual Customer | 'Customer simulation' program that simulates customers placing orders | Order service invocation |
| Virtual Worker | 'Worker simulation' program that simulates the completion of customer orders | Cron input binding, Make-line service invocation | 

### Prerequisites

- Kubernetes cluster of your choice (3+ nodes recommended)
- Helm

### Deployment

Run the commands in the [./conductor-setup](./conductor-setup.sh) script one by one for best results. All Dapr components are deployed on Kubernetes today.

### Build and push docker images to the registry individually.

Inside each service folder there is a Makefile. Navigate to the folder and run the command below to build and push the images. 

  ```bash
  cd go/order-service
  make build_image
  ```
#### About container architecture

The containers are currently built for both ARM and AMD64 architectures. 

### Running locally

# To test the services locally, you can run the following command: 

    ```bash
    cd go/order-service
    make run
    ```