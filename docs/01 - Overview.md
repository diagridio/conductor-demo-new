# Overview 

## App description:

The Order System is composed of five Go microservices, one .NET 8 Actor, and one .NET 8 Dapr Workflow. All services are Dapr enabled. The system simulates a pharmacy's order processing system where customers place orders for their prescriptions onto a central order application. The Order Service then publishes the order summary onto a message bus (on the /orders topic) where it is picked up by multiple subscribers that all perform different actions. The Make-line Service holds the state of the orders to be processed and is invoked by the Virtual Worker service when prescriptions have been completed. The Loyalty Service tracks customer order points based on the number of orders they've placed and the Receipt Service acts as an archival service for auditing all prescription receipts. Finally, the Order Service Workflow manages the chained process that completes the order.

The demo is capable of showcasing features for the Enteprise and Free tiers for Conductor.

For more information on the app services refer to the main [README.md](./../README.md) file.

## Clusters

The demo is deployed in two separate Conductor clusters.

| **Tier**          |**Org**| **Cluster** |  **k8s Cluster** |
|------------------|-------------------|-----------------|--------------------------------|
| Enterprise |diagrid-cloud-conductor-stg| Order-System-Demo | [gke-n-dataplane-demo-us-west1-a0](https://console.cloud.google.com/kubernetes/clusters/details/us-west1-a/gke-n-dataplane-demo-us-west1-a0/details?project=prj-dataplane-n-demo-30534) |
| Free |Demo Free Org| Order-System-Demo | [gke-n-dataplane-demo-us-west1-a1](https://console.cloud.google.com/kubernetes/clusters/details/us-west1-a/gke-n-dataplane-demo-us-west1-a1/details?project=prj-dataplane-n-demo-30534) |
| Enterprise |Demo Org| Order-System-Demo |  [gke-n-dataplane-demo-us-west1-a2](https://console.cloud.google.com/kubernetes/clusters/details/us-west1-a/gke-n-dataplane-demo-us-west1-a2/details?project=prj-dataplane-n-demo-30534) |  

## Connecting to the clusters and modifying Dapr components

### Accessing the GKE clusters

Run the command below to login to Google Cloud:

```
gcloud auth login
```

Get the credentials, replacing `CLUSTER-NAME` with the cluster that you want to connect:

```
gcloud container clusters get-credentials [CLUSTER-NAME] --zone us-west1-a --project prj-dataplane-n-demo-30534
```

For more details, see [these instructions](https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-access-for-kubectl) to connect to the GCP Clusters above to manage the environment on k8s.

### Clone this repository

```git
git clone https://github.com/diagridio/conductor-demo-new.git

cd conductor-demo-new
```

### Components

All components are in the folder `components/k8s`. During the demo, you will need to apply new configurations. To do so, modify one or more yaml files and run:

```bash
kubectl apply -f components/k8s
```

**Next:** [Demo scenario - Conductor overview](./02%20-%20Conductor.md)

