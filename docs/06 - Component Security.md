# Demo Scenario - Component Security

## Overview 
This scenario showcases how Conductor reacts when a component has a security issue.

## Where can you see this error?

### Advisor tab
![Screenshot 2024-05-15 at 1 06 29 PM](https://github.com/diagridio/conductor-demo-new/assets/1051195/2e34f8cb-832e-44c4-b821-35f048201e5a)


![330959227-c5fc613f-a310-4da2-aac9-07a9ab6ddc49](https://github.com/diagridio/conductor-demo-new/assets/1051195/d18df5d2-e76a-4281-ae1f-f8e9518d81c5)


### Components insights on Cluster Summary
![Screenshot 2024-05-15 at 1 07 13 PM](https://github.com/diagridio/conductor-demo-new/assets/1051195/8bbadd99-36c0-4d16-984b-8fa1bd6d891e)

## How to fix this issue?

To mitigate, modify the file `components/k8s/oms.pubsub.yaml` replacing the hardcoded password with:

```yaml
- name: saslPassword
    secretKeyRef:
      name: kafka-password
      key: kafka-password
```

Apply the configuration with `kubectl apply -f components/k8s` and wait a few minutes for the issue to be resolved in conductor.

**Next:** [Demo scenario - Component Invocation](./07%20-%20Component%20Invocation.md) 
