# Conductor Demo Scenarios




##  









## App-specific view

OMS Demo Cluster Applications -> Receipt Generation Service App Summary:
Dapr specific properties for an app on your cluster including name (app-id), active components (used in the past 5 minutes), health (including both app and Dapr sidecar containers) and the annotations (by clicking on `More Info`).

Near realtime resource information like CPU and memory resources including the percentage of memory and CPU used compared to the CPU request and memory limit values. You can also see a breakdown of Dapr and App resource usage over time in the graphs.

-> Throughputs page shows the golden metrics graphs over time. For instance the latency, throughput and error rate for both Dapr HTTP and gRPC calls broken down by API calls.

Notice that digging into the `receipt-generation-service` you will see a high error rate for the `HTTP Request Error Rate` and `gRPC Request Error Rate` graphs. Specifically the call to the Redis output binding from the app seems to be failing 100% of the time.

![image](assets/receipt-generation-svc-errorrate.png)

- Looking at the Dapr Components enabled for the Receipt Generation Service, I have a Redis output binding and a Kafka pubsub broker scoped to this app.
    - Scroll down on the components page to notice that there is a high component error rate (~100%).

![image](assets/receipt-generation-svc-componenterror.png)

-> Navigating to the Notifications tab within the app shows the logs error many times from the Dapr sidecar logs. All Dapr sidecar logs of levels error, warning and fatal are collected and shown in this page.

![image](assets/receipt-generation-svc-logs.png)

Filter on `metrics` to see what default metrics rules are being triggered for this application. A combination of the following default metrics rules should be firing:
- App Component Error Rate has been more than 30 percent for more than 2m
- App to Dapr GRPC Error Rate has been more than 30 percent for more than 2m
- App to Dapr HTTP Error Rate has been more than 30 percent for more than 2m

-> Navigate to the Notifications blade and click on `Rules` in order to see what alerts are created in the organization by default. Can also create your own rules based on the criteria of your apps and send to channels using webhooks, email and of course, the Conductor dashboard.

## Perform a Dapr upgrade on your cluster

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


## Other scenarios:

### View your organizations audit log

### Monitor resiliency policies

### Add a user to your organization

To access the Conductor UI and API you need to have a User Account associated with your organizational email and a Conductor RBAC role. Navigate to the Users page to create a new user.
- Name, email and the role are required.
  - For role, choose `Admin` for managing all Conductor resources including RBAC control, `Editor` for managing everything except for RBAC and `Viewer` for only having access to read your Condcutor resources.

### Configure SSO for your organization

### Get Conductor support

In the top right-hand corner there is a `?` icon with links to the Conductor Documentation, Release notes, Live Chat and Support email. Searching the documentation first is recommended.

Using the chat window in the bottom right hand corner of the Conductor UI, you can access us via Slack during our business hours. If you are outside of our business hours, please email us at `help@diagrid.io`.
