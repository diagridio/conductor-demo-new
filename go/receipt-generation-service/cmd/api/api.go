package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	// "fmt"
	// "os"
	// "time"

	dapr "github.com/dapr/go-sdk/client"
)

// Handles the update loyalty endpoint
func (app *Config) HandleGenerateReceipt(w http.ResponseWriter, r *http.Request) {

	// Unmarshall customer order
	var orderSummary OrderSummary
	err := app.readJSON(w, r, &orderSummary)
	if err != nil {
		log.Printf("Error unmarshalling order summary. Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("Received Order Summary: %v.", orderSummary.OrderID)
	log.Println("Binding with Redis")

	orderContent, err := json.Marshal(orderSummary)
	if err != nil {
		log.Printf("Error marshalling order summary. Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	//Mocking call to non-existing method and app-id to demonstrate error handling in Conductor
	_, err = app.daprClient.InvokeMethod(context.Background(), "receipt-generation-service", "non-existing-method", "GET")
	if err != nil {
		log.Printf("Error invoking method. Error: %v", err)
	}

	// TODO: Swap this out for the above method and get it to show in Conductor.
	// daprHttpPort := os.Getenv("DAPR_HTTP_PORT")
	// if daprHttpPort == "" {
	// 	daprHttpPort = "3500"
	// }

	// client := &http.Client{
	// 	Timeout: 15 * time.Second,
	// }

	// url := fmt.Sprintf("http://localhost:%s/non-existing-method/100", daprHttpPort)
	// req, err := http.NewRequest(http.MethodGet, url, nil)
	// if err != nil {
	// 	panic(err)
	// 	log.Printf("Error: %v", err)
	// }

	// // Adding target app id as part of the header
	// req.Header.Add("dapr-app-id", "loyalty-service")

	// // Invoking a service
	// resp, err := client.Do(req)
	// log.Printf("Service invocation to loyalty service received response: %v %v", resp.StatusCode, http.StatusText(resp.StatusCode))
	// if err != nil {
	// 	panic(err)
	// 	log.Printf("Error: %v", err)
	// }

	//create metadata map
	var metadata map[string]string = make(map[string]string)

	//append metadata with order id
	//metadata["receiptName"] = orderSummary.OrderID
	metadata["key"] = orderSummary.OrderID

	// Redis output binding
	// Insert order using Dapr output binding via Dapr SDK
	in := &dapr.InvokeBindingRequest{
		Name:      ReceiptBindingName,
		Operation: "create",
		Data:      orderContent,
		Metadata:  metadata,
	}
	err = app.daprClient.InvokeOutputBinding(context.Background(), in)
	if err != nil {
		log.Printf("Error invoking output binding. Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)

		return
	}

	log.Printf("Receipt for order %v generated successfully.", orderSummary.OrderID)
	app.writeJSON(w, http.StatusOK, orderSummary.OrderID)
}
