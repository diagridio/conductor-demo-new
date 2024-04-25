## App description:

The Order Management System (OMS) is composed of five Go microservices, one .NET 8 Dapr Workflow, and a Python API. All services are Dapr enabled. The system simulates a pharmacy's order processing system where customers place orders for their prescriptions onto a central order application. The Order Service then publishes the order summary onto a message bus (on the /orders topic) where it is picked up by multiple subscribers that all perform different actions. The Make-line Service holds the state of the orders to be processed and is invoked by the Virtual Worker service when prescriptions have been completed. The Loyalty Service tracks customer order points based on the number of orders they've placed and the Receipt Service acts as an archival service for auditing all prescription receipts. Finally, the Order Service Workflow manages the chained process that completes the order.

The demo is capable of showcasing features for the Enteprise and Free tiers for Conductor.

> Important: There is [failing code](https://github.com/diagridio/conductor-demo-new/blob/main/OrderManagementSystem.ReceiptGenerationService/Controllers/ReceiptGenerationController.cs#L36) in the receipt-generation-service and [lagging code](https://github.com/diagridio/conductor-demo/blob/main/OrderManagementSystem.LoyaltyService/Controllers/LoyaltyController.cs#L64) in the loyalty-service to show how Conductor deals with failures. These are the only **expected** failures in the app.

For more information on the app services read [./README.md](./README.md).

## Clusters

The demo is deployed in two separate Conductor clusters.

| **Tier**          |**Org**| **Cluster** |  **k8s Cluster** |
|------------------|-------------------|-----------------|--------------------------------|
| Free |Demo Free Org| Order-System-Demo | [gke-n-dataplane-demo-us-west1-a1](https://console.cloud.google.com/kubernetes/clusters/details/us-west1-a/gke-n-dataplane-demo-us-west1-a1/details?project=prj-dataplane-n-demo-30534) |
| Enterprise |Demo Org| Order-System-Demo |  [gke-n-dataplane-demo-us-west1-a2](https://console.cloud.google.com/kubernetes/clusters/details/us-west1-a/gke-n-dataplane-demo-us-west1-a2/details?project=prj-dataplane-n-demo-30534) |  
