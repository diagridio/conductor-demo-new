# Demo Scenario - Service invocation issue

## Overview

The `virtual-customer` main path sends orders to `order-service`, but there's an induced error that calls a non-existing service from `receipt-generation-service` ~45% of the time.

```golang
// ~45% of the time, let's mock a call to a non-existing method from an existing app-id to demonstrate error handling in Conductor
if rand.Float64() < 0.45 {
    log.Println("Mocking  call to non-existing service!")
    _, err = client.InvokeMethod(context.Background(), "receipt-generation-service", "non-existing-method", "GET")
    if err != nil {
        log.Printf("Error invoking method. Error: %v", err)
    }
}
```

Use this error to showcase how Conductor displays service invoction issues between apps.

## Where can you see this error?

This error is displayed within the App Graph by isolating the `virtual-customer` app and verifying the error rate to `receipt-generation-service`.

<p align="center">
<img src="./images/virtual-customer-isolate.png" border="10"/>
</p>
