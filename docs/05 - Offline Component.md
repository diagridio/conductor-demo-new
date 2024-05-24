# Demo Scenario - Offline components

## Overview

This scenario showcases how Conductor reacts when a component is suddenly offline.

## How to trigger

On a terminal run the following command to set the number of kafka replicas to 0:

```
kubectl scale statefulset kafka-controller --replicas=0 -n kafka
```

## Where can you see this error?

### `Components` tab

![Screenshot 2024-05-15 at 2 00 59 PM](https://github.com/diagridio/conductor-demo-new/assets/1051195/c5fc613f-a310-4da2-aac9-07a9ab6ddc49)
![Screenshot 2024-05-15 at 2 01 13 PM](https://github.com/diagridio/conductor-demo-new/assets/1051195/7c11e13c-3004-4e15-9221-f96e8d51c82c)

### **App Summary - `order-service` Insights

![Screenshot 2024-05-15 at 2 05 33 PM](https://github.com/diagridio/conductor-demo-new/assets/1051195/966c1e8a-d755-4864-9ca3-02d0a2c0df90)

### App Summary - gRPC Request Error Rate on `order-service`

![Screenshot 2024-05-15 at 2 03 23 PM](https://github.com/diagridio/conductor-demo-new/assets/1051195/e90fa4b3-c4db-4832-9876-698813c5cf25)

### **App Notifications**

![Screenshot 2024-05-15 at 2 54 59 PM](https://github.com/diagridio/conductor-demo-new/assets/1051195/83849c12-f57d-4080-b3e9-435d5398c103)

### App Graph

TODO: ADD WORKING SCREENSHOT

## How to fix this error?

On a terminal run the following command to set the number of kafka replicas to 3:

```
kubectl scale statefulset kafka-controller --replicas=3 -n kafka
```
