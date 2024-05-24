# Connecting a Cluster to Conductor

To connect a cluster we need to complete a few steps:

- Access the cluster via `kubectl`
- Deploy the application
- Connect the cluster on Conductor

## Access the cluster via `kubectl`

Follow [these instructions](https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-access-for-kubectl) to connect to the GCP Clusters above to manage the environment on k8s.

## Clone this repository

```git
git clone https://github.com/diagridio/conductor-demo-new.git

cd conductor-demo-new
```

## Components

All components are in the folder `components/k8s`. During the demo, you will need to apply new configurations. To do so, modify one or more yaml files and run:

```bash
kubectl apply -f components/k8s
```

Before running through the demo scenarios:
- Downgrade or revert the Dapr version of the demo cluster so that it can be upgraded using Conductor. Set `mtls` and `ha` arguments to false to enable more advisories

> Important: Some of these features may not be avaialble for the `Free` version.
