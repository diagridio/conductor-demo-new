# Demo workflow 

## Accessing the GKE clusters

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
- Cycle the `virtual-customer` app pods to create load on the application using the `Rollout Application` button.
- Downgrade or revert the Dapr version of the demo cluster so that it can be upgraded using Conductor. Set `mtls` and `ha` arguments to false to enable more advisories

> Important: Some of these features may not be avaialble for the `Free` version.

## Recommended Demo Flow

- [Clusters overview](./03%20-%20Clusters.md)
- [Advisor](./04%20-%20Advisors.md)
- [Apps graph](./05%20-%20App%20graph.md)
  - [App-specific view](./06%20-%20Application%20view.md)
- [Upgrade Dapr installation](./07%20-%20Upgrade%20Dapr.md)

## Induced errors

Errors can be seen throughout the demo workflow, for details go to [Induced errors.md](./08%20-%20Induced%20errors.md)
