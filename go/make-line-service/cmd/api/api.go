package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dapr/go-sdk/client"
	"github.com/go-chi/chi/v5"
)

type JSONObj struct {
	PubsubName string `json:"pubsubName"`
	Topic      string `json:"topic"`
	Route      string `json:"route"`
}

type Result struct {
	Data any `json:"data"`
}

// Handles Dapr Endpoint and registers the subscription endpoint with the topic, pubsub and route
func (app *Config) HandleDaprEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Print("Received request from Dapr")
	jsonData := []JSONObj{
		{
			PubsubName: PubSubName,
			Topic:      OrderTopic,
			Route:      OrderRoute,
		},
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		log.Printf("Error Marshalling json data. Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
	}

	log.Print("Writing response to Dapr")
	_, err = w.Write(jsonBytes)
	//err = app.writeJSON(w, http.StatusOK, "{}")
	if err != nil {
		log.Printf("Error writing json response. Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
	}
}

// Handles the make line endpoint
func (app *Config) HandleAddOrderToMakeLine(w http.ResponseWriter, r *http.Request) {

	// Unmarshall customer order
	log.Printf("Unmarshalling order summary from topic.")
	var result Result
	err := app.readJSON(w, r, &result)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	s, err := json.Marshal(result.Data)
	if err != nil {
		log.Printf("Error marshalling order summary. Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	var orderSummary OrderSummary
	err = json.Unmarshal(s, &orderSummary)
	if err != nil {
		log.Printf("Error unmarshalling order summary. Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
		return
	}
	log.Printf("Received Order Summary: %v.", orderSummary.OrderID)

	log.Printf("Getting orders for store: %v.", orderSummary.StoreID)
	// Get Store Orders syncronously
	waitGroup := sync.WaitGroup{}
	orders, err := app.getStoreOrders(orderSummary.StoreID)
	if err != nil {
		log.Printf("Error getting orders from store. Error: %v", err)
		app.writeError(w, err, http.StatusInternalServerError)

		return
	}
	waitGroup.Wait()

	// adding ordersummary to order array from state store
	log.Printf("Adding order to makeline store: %v.", orderSummary.StoreID)
	orders = append(orders, orderSummary)

	err = app.updateOrderState(orders, orderSummary.StoreID)
	if err != nil {
		log.Printf("Error adding order to store. Error: %v", err)
		app.writeError(w, err, http.StatusInternalServerError)

		return
	}

	log.Printf("Order %v add to store %v", orderSummary.OrderID, orderSummary.StoreID)
	app.writeJSON(w, http.StatusOK, orderSummary.OrderID)
}

// Completes order by removing it from state store and publishing to topic
func (app *Config) HandleDeleteOrder(w http.ResponseWriter, r *http.Request) {
	storeId := chi.URLParam(r, "storeId")
	orderId := chi.URLParam(r, "orderId")

	if storeId == "" || orderId == "" {
		app.writeError(w, errors.New("no store id or order id provided"), http.StatusBadGateway)
		return
	}

	log.Printf("Trying to complete order %v from %v.", orderId, storeId)

	orders, err := app.getStoreOrders(storeId)
	if err != nil {
		log.Printf("Error getting orders from store. Error: %v", err)
		app.writeError(w, err, http.StatusInternalServerError)
		return
	}

	log.Printf("Finding order %v", orderId)

	var deletedOrder OrderSummary
	// check if there's only one order in the array

	orders, deletedOrder, err = deleteOrder(orders, orderId)
	if err != nil {
		log.Printf("Error deleting order from array. Error: %v", err)
		app.writeError(w, err, http.StatusInternalServerError)
		return
	}

	if deletedOrder.OrderID == "" {
		log.Printf("Order does not exist. Error: %v", err)
		app.writeError(w, err, http.StatusNotFound)
		return
	}

	log.Printf("Removing order from state store: %v.", storeId)
	err = app.updateOrderState(orders, storeId)
	if err != nil {
		log.Printf("Error deleting order from state store. Error: %v", err)
		app.writeError(w, err, http.StatusInternalServerError)
		return
	}

	//updating order status to completed
	deletedOrder.OrderCompletedDate = time.Now()

	// sending order to pub/sub
	ctx := context.Background()
	err = app.daprClient.PublishEvent(ctx, PubSubName, OrderCompletedTopic, deletedOrder)
	if err != nil {
		log.Printf("Error publishing the completed order summary: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
		return
	}
	app.writeJSON(w, http.StatusAccepted, orderId)
}

// Delete order from array of orders based on order id
func deleteOrder(orders []OrderSummary, orderId string) ([]OrderSummary, OrderSummary, error) {

	log.Printf("Deleting order with ID: %v", orderId)

	if len(orders) == 1 {
		// if only one order in array, check if it's the order to be deleted
		if orders[0].OrderID == orderId {
			//if only one order in array, set deleted order to that order
			//set orders as empty array

			return nil, orders[0], nil
		} else {
			//order does not exist
			return orders, OrderSummary{}, errors.New("order not found")
		}
	}

	var deletedOrder OrderSummary
	// Find the order in the array and store the index
	index := 0
	for i, order := range orders {
		if order.OrderID == orderId {
			index = i
			deletedOrder = order
			break
		}
	}

	//order not found
	if deletedOrder.OrderID == "" {
		return orders, deletedOrder, errors.New("order not found")
	}

	// recreating the orders array by appending the first part of the array before the order to be deleted
	// and the second part of the array after the order to be deleted
	orders = append(orders[:index], orders[index+1:]...)
	if len(orders) == 0 {
		return nil, deletedOrder, errors.New("error deleting order")
	}

	return orders, deletedOrder, nil
}

// Update make line state store with storeid: array of orders
func (app *Config) updateOrderState(orders []OrderSummary, storeId string) error {
	// marshall order summary to save to state store
	data, err := json.Marshal(orders)
	if err != nil {
		log.Printf("Error marshalling order summary. Error: %v", err)
		return err
	}

	// save order summary to state store
	err = app.daprClient.SaveState(context.Background(), MakeLineStateStoreName, storeId, data, nil)
	if err != nil {
		log.Printf("Error orders to state store. Error: %v", err)
		return err
	}

	log.Printf("Updated orders for store: %v.", storeId)

	return err
}

func (app *Config) HandleGetOrdersByStoreID(w http.ResponseWriter, r *http.Request) {
	storeId := chi.URLParam(r, "storeId")
	if storeId == "" {
		app.writeError(w, errors.New("no store id provided"), http.StatusBadGateway)
		return
	}

	orders, err := app.getStoreOrders(storeId)
	if err != nil {
		app.writeError(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, orders)
}

// Get all orders from state store where the key is the store ID
func (app *Config) getStoreOrders(storeId string) ([]OrderSummary, error) {
	var stateItem *client.StateItem

	// Get the current loyalty information for the customer from the loyalty state store
	stateItem, err := app.daprClient.GetState(context.Background(), MakeLineStateStoreName, storeId, nil)
	if err != nil {
		log.Printf("Error getting orders from store. Error: %v", err)
		return nil, err
	}

	var orderSummary []OrderSummary

	//if state store doesn't exist, return empty array
	if stateItem.Value == nil {
		return nil, nil
	} else {
		//if state exists populate order summary array with state store information
		err = json.Unmarshal(stateItem.Value, &orderSummary)
		if err != nil {
			log.Printf("Error unmarshalling order array. Error: %v", err)
			return nil, err
		}
		return orderSummary, nil
	}
}
