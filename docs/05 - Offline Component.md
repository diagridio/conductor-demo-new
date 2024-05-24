# Demo Scenario - Offline components

## Overview

This scenario showcases how Conductor reacts when a component is suddenly offline.

## How to trigger

On a terminal run the following command to set the number of kafka replicas to 0:

```
kubectl scale statefulset kafka-controller --replicas=0 -n kafka
```

## Where can you see this error?

### Components tab

Within the components tab, highligth the error icon, scope affected, and see raw error by clicking _View Raw_.

![Screenshot 2024-05-15 at 2 00 59 PM](https://github.com/diagridio/conductor-demo-new/assets/1051195/c5fc613f-a310-4da2-aac9-07a9ab6ddc49)

<p align="center">
<img src="https://github.com/diagridio/conductor-demo-new/assets/1051195/7c11e13c-3004-4e15-9221-f96e8d51c82c" border="10"/>
</p>

### App - order-service

Open the applications tab and select `the order-service` app. Highlight the `Insights` error at the top. 

![Screenshot 2024-05-15 at 2 05 33 PM](https://github.com/diagridio/conductor-demo-new/assets/1051195/966c1e8a-d755-4864-9ca3-02d0a2c0df90)

on the _Metrics_ section, select _Throughputs_ and look at the `gRPC Request Error Rate` graph. Showcase the number of error requests growing.

![Screenshot 2024-05-15 at 2 03 23 PM](https://github.com/diagridio/conductor-demo-new/assets/1051195/e90fa4b3-c4db-4832-9876-698813c5cf25)

Lastly, navigate to the _Notifications_ tab and display the **Component rule** and the **Application rule** errors.

![330971108-83849c12-f57d-4080-b3e9-435d5398c103](https://github.com/diagridio/conductor-demo-new/assets/1051195/2c448432-c22b-453c-b582-0a15cd2c63f8)


### App Graph

**TODO: ADD WORKING SCREENSHOT AFTER CONDUCTOR BUG IS FIXED**

## How to fix this error?

On a terminal run the following command to set the number of kafka replicas to 3:

```
kubectl scale statefulset kafka-controller --replicas=3 -n kafka
```
