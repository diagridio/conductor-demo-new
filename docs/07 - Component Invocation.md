# Demo Scenario - Component invocation issue

## Overview

The [Redis binding spec](https://docs.dapr.io/reference/components-reference/supported-bindings/redis/) states that every _create_ request requires a key _key_ as metadata. 

Within the Receipt Service [api.go](https://github.com/diagridio/conductor-demo-new/blob/demo-scenarios/go/receipt-generation-service/app/main.py) file, the following snippet induces a malformed metadata key to be sent to Redis ~35% of the time for all invocations:

```go
// ~35% of the time, induce error when using Redis. No key "key" in the request
if rand.Float64() < 0.35 {
  metadata = make(map[string]string)
  metadata["orderId"] = orderSummary.OrderID
}
```

This error showcases how Conductor can help you identify underlying issues with the applications code.

## Where can you see this error?

### App Insights - Critical metrics issue

First as a critical metric within App Insights. If the error is not showing, you can increase the `Time Span` in the _Notifications_ page to find it.

![metrics-alert](./images/metrics-alert.png)

### Apps Graph

By isolating `receipt-generation-service` you can see the error rate, RPS and latency:

![receipt-generation-service-isolate](./images/receipt-generation-service-isolate.png)

**Next:** [Demo scenario - Component Misconfiguration](./08%20-%20Component%20Misconfiguration.md) 
