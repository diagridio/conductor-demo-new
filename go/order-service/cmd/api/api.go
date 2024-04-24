package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

const (
	OrderTopic = "orders"
	PubSubName = "oms.pubsub"
)

// Handles the update loyalty endpoint
func (app *Config) HandleHealthz(w http.ResponseWriter, r *http.Request) {

	//set app port
	daprHttpPort := "5180"
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

func (app *Config) HandleCreateOrderSummary(w http.ResponseWriter, r *http.Request) {

	// Unmarshall customer order
	var customerOrder CustomerOrder
	err := app.readJSON(w, r, &customerOrder)
	if err != nil {
		log.Printf("Error unmarshalling order summary! Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("Customer %v %v is placing an order in %v.", customerOrder.FirstName, customerOrder.LastName, customerOrder.StoreID)

	// Create order summary
	summary, err := app.createOrderSummary(customerOrder)
	if err != nil {
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	// Send order to pub/sub
	ctx := context.Background()
	err = app.daprClient.PublishEvent(ctx, PubSubName, OrderTopic, summary)
	if err != nil {
		log.Printf("Error publishing the order summary: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("Order  %v created for customer %v %v", summary.OrderID, summary.FirstName, summary.LastName)

	app.writeJSON(w, http.StatusCreated, summary.OrderID)
}

func findProduct(products []Product, productId int) *Product {
	// Iterate through the list of products to find the one that cointains the id specified
	for _, p := range products {
		if p.ProductID == productId {
			return &p
		}
	}

	return nil
}

func (app *Config) createOrderSummary(customerOrder CustomerOrder) (OrderSummary, error) {
	// Iterate through the list of ordered items to calculate
	// the total and compile a list of item summaries.
	orderTotal := 0.0
	orderItemSummaries := make([]OrderItemSummary, 0)

	// Iterate through the list of ordered items to calculate
	for _, item := range customerOrder.OrderItems {
		product := findProduct(app.products, item.ProductID)
		if product == nil {
			log.Printf("product not found: %v", item.ProductID)
			continue
		}

		orderTotal += product.UnitPrice * float64(item.Quantity)
		orderItemSummaries = append(orderItemSummaries, OrderItemSummary{
			ProductID:   product.ProductID,
			ProductName: product.ProductName,
			Quantity:    item.Quantity,
			UnitCost:    product.UnitCost,
			UnitPrice:   product.UnitPrice,
			ImageUrl:    product.ImageUrl,
		})
	}

	// Initialize and return the order summary
	summary := OrderSummary{
		OrderID:    uuid.New().String(),
		OrderDate:  time.Now(),
		StoreID:    customerOrder.StoreID,
		FirstName:  customerOrder.FirstName,
		LastName:   customerOrder.LastName,
		LoyaltyID:  customerOrder.LoyaltyId,
		OrderItems: orderItemSummaries,
		OrderTotal: orderTotal,
	}

	return summary, nil
}

func (app *Config) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	// retrieve all the products
	if app.products == nil || len(app.products) == 0 {
		app.writeError(w, errors.New("no products found"), http.StatusNotFound)
	}

	log.Printf("Retrieving %d products", len(app.products))
	app.writeJSON(w, http.StatusOK, app.products)
}
