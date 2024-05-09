package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	StoreId                  = "Seattle"
	MakeLineServiceAppId     = "make-line-service"
	MinSecondsToCompleteItem = 1
	MaxSecondsToCompleteItem = 3
)

var storeId = StoreId
var makeLineServiceAppId = MakeLineServiceAppId
var minSecondsToCompleteItem = MinSecondsToCompleteItem
var maxSecondsToCompleteItem = MaxSecondsToCompleteItem

func populateConstants() {
	// get store id
	if sid := os.Getenv("STORE_ID"); sid != "" {
		storeId = sid
	}

	// get max item quantity
	if miq := os.Getenv("MIN_SECONDS_TO_COMPLETE_ITEM"); miq != "" {
		minSecondsToCompleteItem, _ = strconv.Atoi(miq)
	}

	// get max unique items per order
	if mui := os.Getenv("MAX_SECONDS_TO_COMPLETE_ITEM"); mui != "" {
		maxSecondsToCompleteItem, _ = strconv.Atoi(mui)
	}
}

// Handles the update loyalty endpoint
func (app *Config) HandleHealthz(w http.ResponseWriter, r *http.Request) {

	//set app port
	daprHttpPort := "5580"
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

// Handles the update loyalty endpoint
func (app *Config) HandleCompleteOrder(w http.ResponseWriter, r *http.Request) {

	log.Printf("Starting virtual worker...")
	populateConstants()

	log.Printf("The VirtualWorker (%v) is checking orders on the make line...", storeId)

	// Get Store Orders syncronously
	waitGroup := sync.WaitGroup{}
	orders, err := app.getOrders(storeId)
	if err != nil {
		log.Printf("Error getting orders from store. Error: %v", err)
		app.writeError(w, err, http.StatusInternalServerError)
		return
	}
	waitGroup.Wait()

	if len(orders) == 0 {
		log.Printf("The make line is empty! Time to drum up some customers!")
		time.Sleep(5 * time.Second)

		app.writeJSON(w, http.StatusOK, "No orders to make.")
		return
	}

	if len(orders) == 1 {
		log.Printf("There is 1 order waiting to be made.")
	} else {
		log.Printf("There are %v orders waiting to be made.", len(orders))
	}

	completeOrdersWaitGroup := sync.WaitGroup{}
	completeOrdersWaitGroup.Add(len(orders))
	for _, order := range orders {
		defer waitGroup.Done()

		log.Printf("The VirtualWorker (%v) is making an order for %v %v...", storeId, order.FirstName, order.LastName)

		for _, orderItem := range order.OrderItems {
			//creating a slow down to simulate the time it takes to make an item
			log.Printf("The VirtualWorker (%v) is making %v %v.", storeId, orderItem.Quantity, orderItem.ProductName)

			time.Sleep(time.Duration(rand.Intn(maxSecondsToCompleteItem-minSecondsToCompleteItem)+minSecondsToCompleteItem) * time.Second)

			log.Printf("The VirtualWorker (%v) completed %v %v.", storeId, orderItem.Quantity, orderItem.ProductName)
		}

		completeWaitGroup := sync.WaitGroup{}
		err := app.completeOrder(order)
		if err != nil {
			log.Printf("%v %v, your order is ready!", order.FirstName, order.LastName)
		}
		completeWaitGroup.Wait()

	}
	completeOrdersWaitGroup.Wait()

	log.Print("The make line is empty! Time to drum up some customers!")
	app.writeJSON(w, http.StatusOK, "")
}

func (app *Config) getOrders(storeId string) ([]OrderSummary, error) {
	// Invoke make-line-service to get the orders
	methodWithParameters := "orders/" + storeId
	response, err := app.daprClient.InvokeMethod(context.Background(), makeLineServiceAppId, methodWithParameters, "GET")
	if err != nil {
		log.Printf("Error getting orders from store Error: %v", err)
		return nil, err
	}

	var orders []OrderSummary
	err = json.Unmarshal(response, &orders)
	if err != nil {
		log.Printf("Error unmarshalling  orders. Error: %v", err)
		return nil, err
	}

	return orders, nil
}

func (app *Config) completeOrder(order OrderSummary) error {
	// Invoke make-line-service to complete the order
	methodWithParameters := "orders/" + storeId + "/" + order.OrderID
	_, err := app.daprClient.InvokeMethod(context.Background(), makeLineServiceAppId, methodWithParameters, "DELETE")
	if err != nil {
		log.Printf("Error completing order. Error: %v", err)
		return err
	}

	log.Printf("Order %s completed!", order.OrderID)

	return nil
}
