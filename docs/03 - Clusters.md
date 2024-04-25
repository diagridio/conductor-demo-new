
# Clusters 

## View all clusters in one place

View all your Kubernetes clusters connected to Conductor.
- Conductor runs as an agent on your cluster that reports its health status back to the UI periodically.
- You can see if Dapr is installed, and the control plane status and version.
- You can see the number of Dapr-enabled apps on your cluster.

> Important: On the `Free` version only one cluster is visible.

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

> Important: Follow the instructions in [Local-Setup.md](./Local-Setup.md) to setup a local cluster manually.

## View the cluster summary

Summary:
After you've created your cluster connection in Conductor then you can view the cluster summary page to look at an overview of what is happening in your cluster from a Dapr perspective.
- A ton of Dapr details shown including whether `mtls` is enabled, the `root certificate expiry`, `dapr version` etc. 
    - For additional details on the Dapr control plane pods, view the Dapr Status page (click on `Healthy` link). Polling every ~5 seconds.
- Diagrid Agent details like the status, version and agent manifests.
- Dapr control plane uptime data is shown along with CPU and memory usage of the control plane pods broken down by component. If the control plane health is degraded or the agent has lost connection, these show up in red.
    - Conductor agent is polling every 15 seconds for health data.
 
> Important: Some of these features may not be avaialble for the `Free` version.

**Next:** [Advisors](./04%20-%20Advisors.md)