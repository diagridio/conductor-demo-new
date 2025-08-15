package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	dapr "github.com/dapr/go-sdk/client"
)

// Handles the Healthz endpoint
func (app *Config) HandleHealthz(w http.ResponseWriter, r *http.Request) {

	//set app port
	daprHttpPort := "5380"
	if value, ok := os.LookupEnv("DAPR_HTTP_PORT"); ok {
		daprHttpPort = value
	}
	_, err := http.Get("http://localhost:" + daprHttpPort + "/v1.0/healthz")

	if err != nil {
		app.writeError(w, err, http.StatusInternalServerError)
		os.Exit(1)
	}

	app.writeJSON(w, http.StatusOK, "Healthy")
}

// Handles the generate receipt endpoint
func (app *Config) HandleGenerateReceipt(w http.ResponseWriter, r *http.Request) {

	var orderSummary OrderSummary
	err := app.readJSON(w, r, &orderSummary)
	if err != nil {
		log.Printf("Error unmarshall ing order summary. Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("Received Order Summary : %v.", orderSummary.OrderID)

	//marshall order summary to save to state store
	data, err := json.Marshal(map[string]string{
		"orderID":   orderSummary.OrderID,
		"storeID":   orderSummary.StoreID,
		"firstName": orderSummary.FirstName,
		"lastName":  orderSummary.LastName,
		"loyaltyID": orderSummary.LoyaltyID,
	})
	if err != nil {
		log.Printf("Error marshalling order. Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	log.Println("Binding with Redis")

	//create metadata map
	var metadata map[string]string = make(map[string]string)

	var keyValue = "Receipt:" + orderSummary.OrderID
	//append metadata with order id
	metadata["key"] = keyValue

	failRateStr := os.Getenv("FAIL_RATE")
	if failRateStr == "" {
		failRateStr = "0.35"
	}
	failRate, err := strconv.ParseFloat(failRateStr, 64)
	if err != nil {
		log.Printf("Invalid FAIL_RATE value '%s'. Defaulting to 0.35. Error: %v", failRateStr, err)
		failRate = 0.35
	}

	if rand.Float64() < failRate {
		log.Printf("Inducing failure based on fail rate: %f", failRate)
		metadata = make(map[string]string)
		metadata["orderId"] = orderSummary.OrderID
	}

	log.Println("Metadata created")
	// Redis output binding
	// Insert order using Dapr output binding via Dapr SDK
	in := &dapr.InvokeBindingRequest{
		Name:      ReceiptBindingName,
		Operation: "create",
		Data:      data,
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

// Handles the generate receipt endpoint
func (app *Config) HandleCronBinding(w http.ResponseWriter, r *http.Request) {

	app.writeJSON(w, http.StatusOK, "Cron job invoked successfully")

}
