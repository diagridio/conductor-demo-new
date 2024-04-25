# Demo Workflow 

## Accessing the GKE clusters

Follow [these instructions](https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-access-for-kubectl) to connect to the GCP Clusters above to manage the environment on k8s.

## Components

All components are in the folder `components/k8s`. During the demo, you will need to apply new configurations. To do so, modify one or more yaml files and run:

```bash
kubectl apply -f components/k8s
```

Before running through the demo scenarios:
- Cycle the `virtual-customer` app pods to create load on the application using the `Rollout Application` button.
- Downgrade or revert the Dapr version of the demo cluster so that it can be upgraded using Conductor. Set `mtls` and `ha` arguments to false to enable more advisories

> Important: Some of these features may not be avaialble for the `Free` version.

## Recommended Demo Flow

- [Clusters overview](./03%20-%20Clusters.md)
- [Advisor](./04%20-%20Advisors.md)
- [Apps graph](./05%20-%20App%20graph.md)
  - [App-specific view](#app-specific-view)
- [Upgrade Dapr installation](#perform-a-dapr-upgrade-on-your-cluster)

[Other scenarios](#other-scenarios)