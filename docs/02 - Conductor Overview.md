# Demo Scenario: Conductor overview

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



## Recommended Demo Flow

- [Clusters overview](./03%20-%20Clusters.md)
- [Advisor](./04%20-%20Advisors.md)
- [Apps graph](./05%20-%20App%20graph.md)
  - [App-specific view](./06%20-%20Application%20view.md)
- [Upgrade Dapr installation](./07%20-%20Upgrade%20Dapr.md)

## Induced errors

Errors can be seen throughout the demo workflow, for details go to [Showcasing errors.md](./08%20-%20Showcasing%20errors.md)
