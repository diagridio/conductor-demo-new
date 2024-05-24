# Demo workflow 

## Conductor Features Overview

Glance through Conductor's tabs providing a quick overview of each one.

### Cluster Summary

Here you can point out:

- Cluster control plane health and uptime
- Conductor agent status logs
- Critical Dapr control plane metrics issues
- Proactive recommendations for optimizing Dapr operations within the cluster via advisor

### Advisor

- Talk about the 4 pillars and the impact of the advisories.
- Autofix an advisory if available.
- Dismiss an advisory.

#### List advisories using the CLI 

1. [Install and configure the CLI](https://docs.diagrid.io/conductor/cli/introduction#installing-the-cli)
2. On a terminal, run the command `diagrid clusters list` and copy your cluster id.
3. Run the command `diagrid advisories list <clusterId>`

Talk about how you can add this process to a GitHub Actions workflow, a GitLab pipeline, or an Azure DevOps pipeline to capture the current state of your cluster advisories during deployment. 

### Applications

- Mention app namespaces, insights, and pod status.
- Show how to rollout specific apps.
- Open App Graph and show the relationship between the apps and components by isolating one of them. A good option is either `order-service` or `order-processor-workflow` as they both have links with components and multiple apps via service invocation or as workflow activities.

### Components

- Mention initialization status, type, scope, and visualize the YAML file.
- Showcase Component Configurator.

### Configurations

- Mention configurations and their scopes.
- Show raw YAML for one of them.

### Subscriptions

- Talk about _declarative_ and _programatic_ subscriptions and how Conductor displays both.

### Resiliency Policies

- Talk about resiliency polices, their scopes.
- Show raw YAML for one of them.

### Actors

## Pieces of information to share

- Conductor polls every 15 seconds for health data.
- Dapr component initialization runs on a minute-by-minute basis.


## Accessing the GKE clusters

View the cluster summary
Summary: After you've created your cluster connection in Conductor then you can view the cluster summary page to look at an overview of what is happening in your cluster from a Dapr perspective.

A ton of Dapr details shown including whether mtls is enabled, the root certificate expiry, dapr version etc.
For additional details on the Dapr control plane pods, view the Dapr Status page (click on Healthy link). Polling every ~5 seconds.
Diagrid Agent details like the status, version and agent manifests.
Dapr control plane uptime data is shown along with CPU and memory usage of the control plane pods broken down by component. If the control plane health is degraded or the agent has lost connection, these show up in red.
Conductor agent is polling every 15 seconds for health data.


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

## Recommended Demo Flow

- [Clusters overview](./03%20-%20Clusters.md)
- [Advisor](./04%20-%20Advisors.md)
- [Apps graph](./05%20-%20App%20graph.md)
  - [App-specific view](./06%20-%20Application%20view.md)
- [Upgrade Dapr installation](./07%20-%20Upgrade%20Dapr.md)

## Induced errors

Errors can be seen throughout the demo workflow, for details go to [Showcasing errors.md](./08%20-%20Showcasing%20errors.md)
