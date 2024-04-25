# Conductor Demo Scenarios

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

## Running from the deployed environment

Follow [these instructions](https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-access-for-kubectl) to connect to the GCP Clusters above to manage the environment on k8s.

#### Components

All components are in the folder `components/k8s`. To apply new configurations, modify one or more files and run:

```bash
kubectl apply -f components/k8s
```

## Failures

The solution contains multiple induced errors that can be displayed and fixed in real-time.

### Component metadata contains sensitive information as plain text

This will show as an unresolved Security recommendation within the **Advisor** section in Conductor.

![Screenshot 2024-04-24 at 6 55 27 PM](https://github.com/diagridio/conductor-demo-new/assets/1051195/39fd85c9-f523-4c24-bcb6-8762f5eddda8)

To mitigate, modify the file `components/k8s/oms.pubsub.yaml` replacing eh hardcoded password with:

```yaml
- name: saslPassword
    secretKeyRef:
      name: kafka-password
      key: kafka-password
```

Apply the configuration with `kubectl apply -f components/k8s` and wait a few minutes for the issue to be resolved in conductor.

### Redis binding error

The [Redis binding spec](https://docs.dapr.io/reference/components-reference/supported-bindings/redis/) states that every _create_ request requires a key _key_ as metadata. 

Within the Receipt Service [main.py](https://github.com/diagridio/conductor-demo-new/blob/demo-scenarios/python/receipt-generation-service/app/main.py) file, the following snippet induces a malformed metadata key to be sent to Redis ~30% of the time for all invocations:

```python
# 30% of the time, induce error when using Redis. No key "key" in the request
if random.random() < 0.3:
  binding_key = {'receiptName': order.orderId}
```

This issue can be observed in multiple places:

#### App Insights - Critical metrics issue

First as a critical metric within App Insights. If the error is not showing, you can increase the `Time Span` in the _Notifications_ page to find it.

![Screenshot 2024-04-24 at 7 40 08 PM](https://github.com/diagridio/conductor-demo-new/assets/1051195/4202191a-4ee0-42a0-94a2-f0ee2e11d0e3)

#### Apps Graph

By isolating `receipt-generation-service`:

![Screenshot 2024-04-24 at 7 42 41 PM](https://github.com/diagridio/conductor-demo-new/assets/1051195/f92d70bb-ed90-4081-b4be-09f141781bfc)

### Component initialization error

The component `components/oms.state.loyalty-fail.yaml` contains a malformed `redisHost` value. This can be detected in the _Component Insights_ section:

![Screenshot 2024-04-24 at 8 04 23 PM](https://github.com/diagridio/conductor-demo-new/assets/1051195/680e18eb-ecd2-47b0-b566-29f572297b63)

You can check the component yaml file to easily detect the error:

![Screenshot 2024-04-24 at 8 04 48 PM](https://github.com/diagridio/conductor-demo-new/assets/1051195/71a1ca42-492b-4e96-ac26-f48e19e4ceb7)

To fix this issue, comment out the file `components/oms.state.loyalty-fail.yaml` and uncomment the content in `components/oms.state.loyalty.yaml`. Then run the command below to apply the configurations and wait a few seconds for the error to disappear. 

```bash
kubectl apply -f components/k8s
```

### Service invocation error

The `virtual-customer` main path sends orders to `order-service`, but there's an induced error that calls a non-existing service from `receipt-generation-service` ~20% of the time.

```golang
// ~20% of the time, let's mock a call to a non-existing method from an existing app-id to demonstrate error handling in Conductor
if rand.Float64() < 0.2 {
    log.Println("Mocking  call to non-existing service!")
    _, err = client.InvokeMethod(context.Background(), "receipt-generation-service", "non-existing-method", "GET")
    if err != nil {
        log.Printf("Error invoking method. Error: %v", err)
    }
}
```

This error is displayed within the App Graph by isolating the `virtual-customer` app and verifying the error rate to `receipt-generation-service`.

![Screenshot 2024-04-24 at 10 43 09 PM](https://github.com/diagridio/conductor-demo-new/assets/1051195/2dfbaa13-83d9-4d69-a6cc-ae000587e057)

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

Summary:
After you've created your cluster connection in Conductor then you can view the cluster summary page to look at an overview of what is happening in your cluster from a Dapr perspective.
- A ton of Dapr details shown including whether `mtls` is enabled, the `root certificate expiry`, `dapr version` etc. 
    - For additional details on the Dapr control plane pods, view the Dapr Status page (click on `Healthy` link). Polling every ~5 seconds.
- Diagrid Agent details like the status, version and agent manifests.
- Dapr control plane uptime data is shown along with CPU and memory usage of the control plane pods broken down by component. If the control plane health is degraded or the agent has lost connection, these show up in red.
    - Conductor agent is polling every 15 seconds for health data.

## View and implement Advisor recommendations

Advisor:
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

-> Clicking on `receipt-generation-service` shows the metrics from the message broker to the subscribing receipt service and it failing to output the receipt to `oms.binding.receipt` by drawing the edge as red. This error is detailed at [Redis Binding Error](https://github.com/diagridio/conductor-demo-new/edit/demo-scenarios/Conductor-Demo-Scenarios.md#redis-binding-error).

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
