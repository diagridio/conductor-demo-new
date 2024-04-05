package main

import "time"

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

type Customer struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// initialize an array of 200 customers with random names
var customers []Customer = []Customer{
	{ID: 1, FirstName: "John", LastName: "Doe"},
	{ID: 2, FirstName: "Jane", LastName: "Doe"},
	{ID: 3, FirstName: "Alice", LastName: "Smith"},
	{ID: 4, FirstName: "Bob", LastName: "Smith"},
	{ID: 5, FirstName: "Charlie", LastName: "Brown"},
	{ID: 6, FirstName: "Lucy", LastName: "Brown"},
	{ID: 7, FirstName: "Eve", LastName: "Johnson"},
	{ID: 8, FirstName: "David", LastName: "Johnson"},
	{ID: 9, FirstName: "Frank", LastName: "Williams"},
	{ID: 10, FirstName: "Grace", LastName: "Williams"},
	{ID: 11, FirstName: "Heather", LastName: "Jones"},
	{ID: 12, FirstName: "Ivan", LastName: "Jones"},
	{ID: 13, FirstName: "Mallory", LastName: "Davis"},
	{ID: 14, FirstName: "Nancy", LastName: "Davis"},
	{ID: 15, FirstName: "Oliver", LastName: "Martinez"},
	{ID: 16, FirstName: "Pamela", LastName: "Martinez"},
	{ID: 17, FirstName: "Quinn", LastName: "Garcia"},
	{ID: 18, FirstName: "Randy", LastName: "Garcia"},
	{ID: 19, FirstName: "Samantha", LastName: "Rodriguez"},
	{ID: 20, FirstName: "Trevor", LastName: "Rodriguez"},
	{ID: 21, FirstName: "Ursula", LastName: "Lopez"},
	{ID: 22, FirstName: "Victor", LastName: "Lopez"},
	{ID: 23, FirstName: "Wendy", LastName: "Perez"},
	{ID: 24, FirstName: "Xavier", LastName: "Perez"},
	{ID: 25, FirstName: "Yvonne", LastName: "Torres"},
	{ID: 26, FirstName: "Zachary", LastName: "Torres"},
	{ID: 27, FirstName: "Amy", LastName: "Hernandez"},
	{ID: 28, FirstName: "Brian", LastName: "Hernandez"},
	{ID: 29, FirstName: "Cindy", LastName: "Gonzalez"},
	{ID: 30, FirstName: "Dennis", LastName: "Gonzalez"},
	{ID: 31, FirstName: "Emily", LastName: "Wilson"},
	{ID: 32, FirstName: "Fred", LastName: "Wilson"},
	{ID: 33, FirstName: "Gina", LastName: "Anderson"},
	{ID: 34, FirstName: "Henry", LastName: "Anderson"},
	{ID: 35, FirstName: "Isabel", LastName: "Thomas"},
	{ID: 36, FirstName: "Jack", LastName: "Thomas"},
	{ID: 37, FirstName: "Kathy", LastName: "Jackson"},
	{ID: 38, FirstName: "Larry", LastName: "Jackson"},
	{ID: 39, FirstName: "Megan", LastName: "White"},
	{ID: 40, FirstName: "Nathan", LastName: "White"},
	{ID: 41, FirstName: "Olivia", LastName: "Harris"},
	{ID: 42, FirstName: "Peter", LastName: "Harris"},
	{ID: 43, FirstName: "Quincy", LastName: "Young"},
	{ID: 44, FirstName: "Rachel", LastName: "Young"},
	{ID: 45, FirstName: "Sally", LastName: "King"},
	{ID: 46, FirstName: "Tom", LastName: "King"},
	{ID: 47, FirstName: "Alex", LastName: "Taylor"},
	{ID: 48, FirstName: "Benjamin", LastName: "Taylor"},
	{ID: 49, FirstName: "Catherine", LastName: "Brown"},
	{ID: 50, FirstName: "Daniel", LastName: "Brown"},
	{ID: 51, FirstName: "Emma", LastName: "Johnson"},
	{ID: 52, FirstName: "Franklin", LastName: "Johnson"},
	{ID: 53, FirstName: "Grace", LastName: "Davis"},
	{ID: 54, FirstName: "Henry", LastName: "Davis"},
	{ID: 55, FirstName: "Isabella", LastName: "Martinez"},
	{ID: 56, FirstName: "Jacob", LastName: "Martinez"},
	{ID: 57, FirstName: "Katherine", LastName: "Garcia"},
	{ID: 58, FirstName: "Liam", LastName: "Garcia"},
	{ID: 59, FirstName: "Mia", LastName: "Rodriguez"},
	{ID: 60, FirstName: "Noah", LastName: "Rodriguez"},
	{ID: 61, FirstName: "Olivia", LastName: "Lopez"},
	{ID: 62, FirstName: "Patrick", LastName: "Lopez"},
	{ID: 63, FirstName: "Quinn", LastName: "Perez"},
	{ID: 64, FirstName: "Ryan", LastName: "Perez"},
	{ID: 65, FirstName: "Sophia", LastName: "Torres"},
	{ID: 66, FirstName: "Thomas", LastName: "Torres"},
	{ID: 67, FirstName: "Ursula", LastName: "Hernandez"},
	{ID: 68, FirstName: "Vincent", LastName: "Hernandez"},
	{ID: 69, FirstName: "Wendy", LastName: "Gonzalez"},
	{ID: 70, FirstName: "Xavier", LastName: "Gonzalez"},
	{ID: 71, FirstName: "Yvonne", LastName: "Wilson"},
	{ID: 72, FirstName: "Zachary", LastName: "Wilson"},
	{ID: 73, FirstName: "Amy", LastName: "Anderson"},
	{ID: 74, FirstName: "Brian", LastName: "Anderson"},
	{ID: 75, FirstName: "Cindy", LastName: "Thomas"},
	{ID: 76, FirstName: "David", LastName: "Thomas"},
	{ID: 77, FirstName: "Emily", LastName: "Jackson"},
	{ID: 78, FirstName: "Fred", LastName: "Jackson"},
	{ID: 79, FirstName: "Gina", LastName: "White"},
	{ID: 80, FirstName: "Henry", LastName: "White"},
	{ID: 81, FirstName: "Isabel", LastName: "Harris"},
	{ID: 82, FirstName: "Jack", LastName: "Harris"},
	{ID: 83, FirstName: "Kathy", LastName: "Young"},
	{ID: 84, FirstName: "Larry", LastName: "Young"},
	{ID: 85, FirstName: "Megan", LastName: "King"},
	{ID: 86, FirstName: "Nathan", LastName: "King"},
	{ID: 87, FirstName: "Oliver", LastName: "Taylor"},
	{ID: 88, FirstName: "Penelope", LastName: "Taylor"},
	{ID: 89, FirstName: "Quincy", LastName: "Brown"},
	{ID: 90, FirstName: "Rachel", LastName: "Brown"},
	{ID: 91, FirstName: "Sally", LastName: "Johnson"},
	{ID: 92, FirstName: "Thomas", LastName: "Johnson"},
	{ID: 93, FirstName: "Ursula", LastName: "Davis"},
	{ID: 94, FirstName: "Victor", LastName: "Davis"},
	{ID: 95, FirstName: "Wendy", LastName: "Martinez"},
	{ID: 96, FirstName: "Xavier", LastName: "Martinez"},
	{ID: 97, FirstName: "Yvonne", LastName: "Garcia"},
	{ID: 98, FirstName: "Zachary", LastName: "Garcia"},
	{ID: 99, FirstName: "Amy", LastName: "Rodriguez"},
	{ID: 100, FirstName: "Brian", LastName: "Rodriguez"},
	{ID: 101, FirstName: "Cindy", LastName: "Lopez"},
	{ID: 102, FirstName: "David", LastName: "Lopez"},
	{ID: 103, FirstName: "Emily", LastName: "Perez"},
	{ID: 104, FirstName: "Fred", LastName: "Perez"},
	{ID: 105, FirstName: "Gina", LastName: "Torres"},
	{ID: 106, FirstName: "Henry", LastName: "Torres"},
	{ID: 107, FirstName: "Isabel", LastName: "Hernandez"},
	{ID: 108, FirstName: "Jack", LastName: "Hernandez"},
	{ID: 109, FirstName: "Kathy", LastName: "Gonzalez"},
	{ID: 110, FirstName: "Larry", LastName: "Gonzalez"},
	{ID: 111, FirstName: "Megan", LastName: "Wilson"},
	{ID: 112, FirstName: "Nathan", LastName: "Wilson"},
	{ID: 113, FirstName: "Olivia", LastName: "Anderson"},
	{ID: 114, FirstName: "Peter", LastName: "Anderson"},
	{ID: 115, FirstName: "Quincy", LastName: "Thomas"},
	{ID: 116, FirstName: "Rachel", LastName: "Thomas"},
	{ID: 117, FirstName: "Sally", LastName: "Jackson"},
	{ID: 118, FirstName: "Tom", LastName: "Jackson"},
	{ID: 119, FirstName: "Ursula", LastName: "White"},
	{ID: 120, FirstName: "Victor", LastName: "White"},
	{ID: 121, FirstName: "Wendy", LastName: "Harris"},
	{ID: 122, FirstName: "Xavier", LastName: "Harris"},
	{ID: 123, FirstName: "Yvonne", LastName: "Young"},
	{ID: 124, FirstName: "Zachary", LastName: "Young"},
	{ID: 125, FirstName: "Amy", LastName: "King"},
	{ID: 126, FirstName: "Brian", LastName: "King"},
	{ID: 127, FirstName: "Cindy", LastName: "Taylor"},
	{ID: 128, FirstName: "David", LastName: "Taylor"},
	{ID: 129, FirstName: "Emily", LastName: "Brown"},
	{ID: 130, FirstName: "Fred", LastName: "Brown"},
	{ID: 131, FirstName: "Gina", LastName: "Johnson"},
	{ID: 132, FirstName: "Henry", LastName: "Johnson"},
	{ID: 133, FirstName: "Isabel", LastName: "Davis"},
	{ID: 134, FirstName: "Jack", LastName: "Davis"},
	{ID: 135, FirstName: "Kathy", LastName: "Martinez"},
	{ID: 136, FirstName: "Larry", LastName: "Martinez"},
	{ID: 137, FirstName: "Megan", LastName: "Garcia"},
	{ID: 138, FirstName: "Nathan", LastName: "Garcia"},
	{ID: 139, FirstName: "Olivia", LastName: "Rodriguez"},
	{ID: 140, FirstName: "Peter", LastName: "Rodriguez"},
	{ID: 141, FirstName: "Quincy", LastName: "Lopez"},
	{ID: 142, FirstName: "Rachel", LastName: "Lopez"},
	{ID: 143, FirstName: "Sally", LastName: "Perez"},
	{ID: 144, FirstName: "Tom", LastName: "Perez"},
	{ID: 145, FirstName: "Ursula", LastName: "Torres"},
	{ID: 146, FirstName: "Victor", LastName: "Torres"},
	{ID: 147, FirstName: "Wendy", LastName: "Hernandez"},
	{ID: 148, FirstName: "Xavier", LastName: "Hernandez"},
	{ID: 149, FirstName: "Yvonne", LastName: "Gonzalez"},
	{ID: 150, FirstName: "Zachary", LastName: "Gonzalez"},
	{ID: 151, FirstName: "Amy", LastName: "Wilson"},
	{ID: 152, FirstName: "Brian", LastName: "Wilson"},
	{ID: 153, FirstName: "Cindy", LastName: "Anderson"},
	{ID: 154, FirstName: "David", LastName: "Anderson"},
	{ID: 155, FirstName: "Emily", LastName: "Thomas"},
	{ID: 156, FirstName: "Fred", LastName: "Thomas"},
	{ID: 157, FirstName: "Gina", LastName: "Jackson"},
	{ID: 158, FirstName: "Henry", LastName: "Jackson"},
	{ID: 159, FirstName: "Isabel", LastName: "White"},
	{ID: 160, FirstName: "Jack", LastName: "White"},
	{ID: 161, FirstName: "Kathy", LastName: "Harris"},
	{ID: 162, FirstName: "Larry", LastName: "Harris"},
	{ID: 163, FirstName: "Megan", LastName: "Young"},
	{ID: 164, FirstName: "Nathan", LastName: "Young"},
	{ID: 165, FirstName: "Oliver", LastName: "King"},
	{ID: 166, FirstName: "Penelope", LastName: "King"},
	{ID: 167, FirstName: "Quincy", LastName: "Taylor"},
	{ID: 168, FirstName: "Rachel", LastName: "Taylor"},
	{ID: 169, FirstName: "Sally", LastName: "Brown"},
	{ID: 170, FirstName: "Thomas", LastName: "Brown"},
	{ID: 171, FirstName: "Ursula", LastName: "Johnson"},
	{ID: 172, FirstName: "Victor", LastName: "Johnson"},
	{ID: 173, FirstName: "Wendy", LastName: "Davis"},
	{ID: 174, FirstName: "Xavier", LastName: "Davis"},
	{ID: 175, FirstName: "Yvonne", LastName: "Martinez"},
	{ID: 176, FirstName: "Zachary", LastName: "Martinez"},
	{ID: 177, FirstName: "Amy", LastName: "Garcia"},
	{ID: 178, FirstName: "Brian", LastName: "Garcia"},
	{ID: 179, FirstName: "Cindy", LastName: "Rodriguez"},
	{ID: 180, FirstName: "David", LastName: "Rodriguez"},
	{ID: 181, FirstName: "Emily", LastName: "Lopez"},
	{ID: 182, FirstName: "Fred", LastName: "Lopez"},
	{ID: 183, FirstName: "Gina", LastName: "Perez"},
	{ID: 184, FirstName: "Henry", LastName: "Perez"},
	{ID: 185, FirstName: "Isabel", LastName: "Torres"},
	{ID: 186, FirstName: "Jack", LastName: "Torres"},
	{ID: 187, FirstName: "Kathy", LastName: "Hernandez"},
	{ID: 188, FirstName: "Larry", LastName: "Hernandez"},
	{ID: 189, FirstName: "Megan", LastName: "Gonzalez"},
	{ID: 190, FirstName: "Nathan", LastName: "Gonzalez"},
	{ID: 191, FirstName: "Olivia", LastName: "Wilson"},
	{ID: 192, FirstName: "Peter", LastName: "Wilson"},
	{ID: 193, FirstName: "Quincy", LastName: "Anderson"},
	{ID: 194, FirstName: "Rachel", LastName: "Anderson"},
	{ID: 195, FirstName: "Sally", LastName: "Thomas"},
	{ID: 196, FirstName: "Tom", LastName: "Thomas"},
	{ID: 197, FirstName: "Ursula", LastName: "Jackson"},
	{ID: 198, FirstName: "Victor", LastName: "Jackson"},
	{ID: 199, FirstName: "Wendy", LastName: "White"},
	{ID: 200, FirstName: "Xavier", LastName: "White"},
}
