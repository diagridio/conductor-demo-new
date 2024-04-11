package main

import (
	"time"
)

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
	MarkedForDeletion  bool               `json:"markedForDeletion"`
}
