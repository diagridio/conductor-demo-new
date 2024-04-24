# Conductor Demo Scenarios

## App description:

The Order Management System (OMS) is composed of five G0 microservices, one .NET 8 Dapr Workflow, and a Python API. All services are Dapr enabled. The system simulates a pharmacy's order processing system where customers place orders for their prescriptions onto a central order application. The Order Service then publishes the order summary onto a message bus (on the /orders topic) where it is picked up by multiple subscribers that all perform different actions. The Make-line Service holds the state of the orders to be processed and is invoked by the Virtual Worker service when prescriptions have been completed. The Loyalty Service tracks customer order points based on the number of orders they've placed and the Receipt Service acts as an archival service for auditing all prescription receipts. Finally, the Order Service Workflow manages the chained process that completes the order.

> Important: There is [failing code](https://github.com/diagridio/conductor-demo-new/blob/main/OrderManagementSystem.ReceiptGenerationService/Controllers/ReceiptGenerationController.cs#L36) in the receipt-generation-service and [lagging code](https://github.com/diagridio/conductor-demo/blob/main/OrderManagementSystem.LoyaltyService/Controllers/LoyaltyController.cs#L64) in the loyalty-service to show how Conductor deals with failures. These are the only **expected** failures in the app.

For more information on the app services read [./README.md](./README.md).

## Prerequisites

Before running through the demo scenarios:
- Cycle the [Virtual Customers app pods](https://conductor.diagrid.io/clusters/160bd018-32a9-477c-ab8a-58d664d55829/app/order-management-system:virtual-customers/detail/summary) to create load on the application using the `Rollout Application` button
- Create an empty cluster ready for the agent to be deployed to on your local machine (eg. Kind, Docker K8s, etc)
- Downgrade or revert the Dapr version of the demo cluster so that it can be upgraded using Conductor. Set `mtls` and `ha` arguments to false to enable more advisories

## Recommended Demo Flow

- [Clusters list](#view-all-clusters-in-one-place)
- [Connect a new cluster](#connect-a-new-cluster-to-conductor)
- [OMS Demo Cluster: Cluster summary](#view-the-cluster-summary)
- [OMS Demo Cluster: Advisor](#view-and-implement-advisor-recommendations)
- [OMS Demo Cluster: Apps graph](#view-the-application-graph)
  - [OMS Demo Cluster: App-specific view](#app-specific-view)
- [OMS Demo Cluster: Upgrade Dapr installation](#perform-a-dapr-upgrade-on-your-cluster)

[Other scenarios](#other-scenarios)

## View all clusters in one place

View all your Kubernetes clusters connected to Conductor.
- Conductor runs as an agent on your cluster that reports its health status back to the UI periodically.
- You can see if Dapr is installed, and the control plane status and version.
- You can see the number of Dapr-enabled apps on your cluster.

## Connect a new cluster to Conductor

First thing that you might want to do is connect a new cluster from the Clusters Dashboard.
- Conductor supports EKS, AKS, GKE and RHOS. These have all been tested but you may try additional distributions at your own risk. Local distros work as well like Kind and Docker Desktop and should use the `native` option on cluster creation.
- Conductor can install Dapr for you, using the version, Helm values and namespace of your choice. If you choose this option while already having Dapr installed on your cluster, Conductor upgrades the version for you.
    - Option to automatically enforce Dapr configuration from Conductor - essentially any changes that are made outside of Conductor will be brought back to the desired state set stored in Conductor (via UI, CLI or operator). Conductor becomes the single source of truth and continuously reconciles the desired state.
    - A number of supported Dapr versions on the list. Conductor lists ~3 minor versions of Dapr, the latest version n, and n-2 as well as various patch versions. Upstream release candidates are loaded into Conductor within hours of becoming available for users to test.
    - Option to perform a rolling update on all Dapr-enabled applications after the Dapr control plane has been installed/upgraded. (This runs `kubectl rollout restart` on the k8s deployment)
    - Two cluster profiles for recommended Development and Production configurations with editable Helm values
- Option to automatically upgrade the Conductor agent. This is recommended as releases are weekly.
- Automatic certificate rotation replaces the Dapr self-signed certificates with Diagrid certs and rotates them before they expire with zero downtime. Option to choose the renewal frequency and renewal window (best effort) of the certificate rotation. This is not possible in upstream unless you bring your own certs and manually rotate them yourself.

After configuring the cluster installation options you will receive a kubectl command to install the Diagrid agent on your cluster. You can optionally download the manifest file to inspect what resources are installed. All of this can also be done using a [K8s operator](https://docs.diagrid.io/operator-overview) for GitOps deployments.

Apply the manifests to a test cluster and navigate to the [dashboard](https://conductor.diagrid.io/clusters) to watch the agent come online and install/upgrade Dapr.

## View the cluster summary

OMS Demo Cluster Summary:
After you've created your cluster connection in Conductor then you can view the cluster summary page to look at an overview of what is happening in your cluster from a Dapr perspective.
- A ton of Dapr details shown including whether `mtls` is enabled, the `root certificate expiry`, `dapr version` etc. 
    - For additional details on the Dapr control plane pods, view the Dapr Status page (click on `Healthy` link). Polling every ~5 seconds.
- Diagrid Agent details like the status, version and agent manifests.
- Dapr control plane uptime data is shown along with CPU and memory usage of the control plane pods broken down by component. If the control plane health is degraded or the agent has lost connection, these show up in red.
    - Conductor agent is polling every 15 seconds for health data.

## View and implement Advisor recommendations

OMS Demo Cluster Advisor:
Looking at the Cluster Advisor page you can see there are a number of recommended actions to take on Conductor for production Dapr deployments. These recommendations are automatically generated based on your cluster configuration and Dapr-enabled applications so just need to install the agent to see the advisories. Advisor recommendations are evaluated every 15 mins but can be synced when desired as well as dismissed if they are not relevant to your deployment.

Over 30 advisories including:

- **Security**:
    - The recommendation for the component metadata containing sensitive information as plaintext is directly linked to the resource it affects. You can see on the component manifest that the pubsub component for Kafka has the `saslPassword` value fully readable by anyone in the cluster.
        - Now you can add a secretKeyRef to the Oms.pubsub.yaml manifest for storing this value and a kubernetes secret to the cluster to get rid of this recommendation.
    - This shows that neither App to Dapr API or Dapr API to App authentication is enabled. Note: Dapr supports token-based auth to increase the security between the app container and the Dapr sidecar. 

    - Other recommendations that could be displayed here are:
      - MTLS between services (enabled on this cluster)
      - Dapr component without scopes
      - T-10 day Dapr certificate expiry
    
- **Reliability**:
    - Since the Dapr control plane is not deployed with a highly available (HA) configuration, there is only one replica of each dapr pod giving a high impact alert. This advisory can be `auto-fixed` by clicking on the `Fix Advisory` button which will update the control plane with the following:

  ```
  global:
  ha:
    enabled: true
  ```
    - Other recommendations that could be displayed here are:
      - If the Dapr Injector watchdog is not configured on this cluster then when new applications come up that are Dapr enabled and the Dapr control plane is not already up, there could be app failures. This is useful in production environments and when you are stopping your Kubernetes cluster.

- **Performance**:
    - Performance optimization advisories: *"Application x resource settings can be optimized"*
      - Advisor will recommend CPU and memory resource optimization values for the top 10 most resource-intensive apps running on your cluster. These will recommend a new value for the Dapr sidecar and app container based on the performance of the containers over the past 10-15 days. Note: if there are no values shown then the apps have not been deployed for at least 10 days.
    - Other recommendations that could be displayed here are:
      - There are no Dapr sidecar resource requests and limits set for these applications so all of the app manifests should be updated with the recommended values.
    
- **Observability**:
    - Many log collectors and search engines require JSON formatted logs for better parsing, so without that configuration turned on these applications are logging to stdout as plaintext. The control plane advisory for JSON-formatted logs can be `auto-fixed` by clicking on the `Fix Advisory` button which will update the control plane with the following:
 ```
  global:
  logAsJson: true
  ```
  - Other recommendations that could be displayed here are:
    - Sidecar log level should not be debug when deployed in production

## View the application graph

OMS Demo cluster:
Applications list shows all the Dapr-enabled apps and infrastructure components running on the cluster. Important Dapr properties such as app-id, health status (combining both app and sidecar container), app port and resource data can be seen here. Option to rollout each of these. Viewing the Dapr version can be handy here to ensure that all apps haven rolled out and are the same as control plane version.

Clicking on `Apps Graph` shows a graphical view of the Dapr applications running and how they communicate in the cluster. Each app can be `isolated` for additional details and the app icons can also be dragged around for easier viewing.
- Solid lines denote service invocation between apps
- Dotted lines denote pubsub
- Green apps are healthy, yellow: degraded, red: unhealthy, grey: unknown
- Isolating apps shows metrics data and infrastructure information of the connected Dapr components
- All orange boxes are Dapr components.

Isolate on the `order-service` to view metrics from it publishing orders to the `oms.pubsub` message broker to a number of subscribers. Isolate on `oms.pubsub` to see all metrics of connected apps that are communicating with the broker. The metrics shown are in near-realtime and are *not* using Dapr tracing configurations but instead just the Dapr metrics data scraped from the Prometheus endpoint on each Dapr sidecar.

-> Clicking on `receipt-generation-service` shows the metrics from the message broker to the subscribing receipt service and it failing to output the receipt to `oms.binding.receipt` by drawing the edge as red similar to the below.

![image](assets/receipt-generation-svc-appgraph.png)

-> Click on the link button next to the `receipt-generation-service` to navigate to the app-specific view.

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