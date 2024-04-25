# Advisors

## View and implement Advisor recommendations

Advisor:
Looking at the Cluster Advisor page you can see there are a number of recommended actions to take on Conductor for production Dapr deployments. These recommendations are automatically generated based on your cluster configuration and Dapr-enabled applications so just need to install the agent to see the advisories. Advisor recommendations are evaluated every 15 mins but can be synced when desired as well as dismissed if they are not relevant to your deployment.

Over 30 advisories including:

- **Security**:
    - The recommendation for the component metadata containing sensitive information as plaintext is directly linked to the resource it affects. You can see on the component manifest that the pubsub component for Kafka has the `saslPassword` value fully readable by anyone in the cluster.
        - Now you can add a secretKeyRef to the Oms.pubsub.yaml manifest for storing this value and a kubernetes secret to the cluster to get rid of this recommendation. See details in [induced errors](./08%20-%20Induced%20errors.md)
    - This shows that neither App to Dapr API or Dapr API to App authentication is enabled. Note: Dapr supports token-based auth to increase the security between the app container and the Dapr sidecar.
    - automountServiceAccountToken is enabled. It’s recommended deploying your apps with `automountServiceAccountToken: false` to improve the security posture of your pods, unless your apps depend on having a Service Account token.

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
    
    - Only 2 out of 5 services have health checks enabled. This advises you to add the following to your annotations within your deployment files:

  ```
    dapr.io/enable-app-health-check: "true"
    dapr.io/app-health-check-path: "/healthz"
    dapr.io/app-health-probe-interval: "3"
    dapr.io/app-health-probe-timeout: "200"
    dapr.io/app-health-threshold: "2"
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