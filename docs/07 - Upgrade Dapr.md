# Perform a Dapr upgrade on your cluster

Upgrade your Dapr installation from the Cluster Summary page by clicking on the top right-hand action menu `Update Dapr`. You see many of the same options that were on the initial cluster connect page. You can upgrade or downgrade your Dapr version to any other Conductor supported Dapr versions which follows the upstream Dapr version support policy. Note: Dapr only lets you upgrade or downgrade one *minor* version at a time. Choose the latest version available to upgrade to.

Add the recommended resource settings for the control plane using the following Helm chart. Conductor always pulls the current Helm values that are installed on your cluster.
- Notice you have the option to choose specific application Deployments to rollout or choose to restart all your dapr-enabled apps after the control plane upgrades. This ensures that all your applications are running the same version of Dapr.
- You also have the option of scheduling the upgrade to a time when there is less usage on your cluster. Dapr upgrades should not typically cause any downtime but your apps to have to restart to get the latest sidecar version.

``` yaml
dapr_dashboard:
  resources:
    limits:
      cpu: 200m
      memory: 200Mi
    requests:
      cpu: 50m
      memory: 20Mi
dapr_operator:
  resources:
    limits:
      cpu: '1'
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 100Mi
  watchInterval: 1m
dapr_placement:
  cluster:
    forceInMemoryLog: true
  resources:
    limits:
      cpu: '1'
      memory: 150Mi
    requests:
      cpu: 250m
      memory: 75Mi
dapr_sentry:
  resources:
    limits:
      cpu: '1'
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 30Mi
dapr_sidecar_injector:
  resources:
    limits:
      cpu: '1'
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 30Mi
global:
  ha:
    enabled: false
    replicaCount: 3
  logAsJson: false
  mtls:
    enabled: false
  registry: docker.io/daprio
```

Confirm your upgrade and then watch the rollout of your Dapr control plane by clicking the label under "Dapr Status" on the cluster summary page and then by looking at the Agent Status Logs and filtering on `dapr/version:upgrading`. You can also check to see that all the apps (or your selected ones) were rolled out successfully when their sidecars are all updated to the correct version by looking at the Applications List.