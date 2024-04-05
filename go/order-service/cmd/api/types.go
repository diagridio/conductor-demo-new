package main

import (
	"log"
	"os"
	"time"
)

const FileName = "DrugStoreProducts-categorized.json"

type CustomerOrder struct {
	StoreID    string      `json:"storeId"`
	FirstName  string      `json:"firstName"`
	LastName   string      `json:"lastName"`
	LoyaltyId  string      `json:"loyaltyId"`
	OrderItems []OrderItem `json:"orderItems"`
}

type OrderItem struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type OrderItemSummary struct {
	ProductID   int     `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	UnitCost    float64 `json:"unitCost"`
	UnitPrice   float64 `json:"unitPrice"`
	ImageUrl    string  `json:"imageUrl"`
}

type OrderSummary struct {
	OrderID            string             `json:"orderId"`
	OrderDate          time.Time          `json:"orderDate"`
	OrderCompletedDate time.Time          `json:"orderCompletedDate,omitempty"`
	StoreID            string             `json:"storeId"`
	FirstName          string             `json:"firstName"`
	LastName           string             `json:"lastName"`
	LoyaltyID          string             `json:"loyaltyId"`
	OrderItems         []OrderItemSummary `json:"orderItems"`
	OrderTotal         float64            `json:"orderTotal"`
}

type Product struct {
	ProductID   int     `json:"productId"`
	ProductName string  `json:"productName"`
	Description string  `json:"description"`
	UnitCost    float64 `json:"unitCost"`
	UnitPrice   float64 `json:"unitPrice"`
	ImageUrl    string  `json:"imageUrl"`
	CategoryID  string  `json:"categoryId"`
}

// Loads products from JSON file
func (app *Config) GetAppProducts() ([]Product, error) {

	fileName := ""
	if fn := os.Getenv("PRODUCT_DEFINITION_FILENAME"); fn != "" {
		fileName = fn
	} else {
		fileName = "cmd/api/product-definitions/" + FileName
	}

	log.Printf("Filename %s\n", fileName)

	//fileName = "app/" + fileName

	var products []Product
	err := app.readJSONFile(fileName, &products)
	if err != nil || len(products) == 0 {
		log.Fatalf("Error loading products from file %s. Error: %s", fileName, err)
		return nil, err
	}

	log.Printf("Loaded %d products from %s\n", len(products), fileName)

	return products, nil
}
