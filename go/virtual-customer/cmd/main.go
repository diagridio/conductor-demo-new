package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"sync"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

const (
	OrderServiceDaprId      = "order-service"
	StoreId                 = "Seattle"
	MaxItemQuantity         = 1
	MaxUniqueItemsPerOrder  = 10
	MinSecondsToPlaceOrder  = 1
	MaxSecondsToPlaceOrder  = 3
	MinSecondsBetweenOrders = 1
	MaxSecondsBetweenOrders = 3
	NumOrders               = 20
)

var storeId = StoreId
var maxItemQuantity = MaxItemQuantity
var maxUniqueItemsPerOrder = MaxUniqueItemsPerOrder
var minSecondsToPlaceOrder = MinSecondsToPlaceOrder
var maxSecondsToPlaceOrder = MaxSecondsToPlaceOrder
var minSecondsBetweenOrders = MinSecondsBetweenOrders
var maxSecondsBetweenOrders = MaxSecondsBetweenOrders
var numOrders = NumOrders

type Config struct{}

var client dapr.Client

func main() {

	ctx := context.Background()
	populateConstants()

	_client, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("error creating dapr client: %v", err)
	}
	defer _client.Close()

	client = _client

	// Get products syncronously
	waitGroup := sync.WaitGroup{}

	products, err := getProducts(ctx)
	if err != nil {
		log.Fatalf("Error getting products: %v", err)
	}

	waitGroup.Wait()

	ordersCreated := 0
	for {

		sleepTime := rand.Intn((maxSecondsBetweenOrders - minSecondsBetweenOrders) + minSecondsToPlaceOrder)
		time.Sleep(time.Duration(sleepTime) * time.Second)

		// Create a new customer order
		order, err := createOrder(products)
		if err != nil {
			log.Printf("Error creating order: %v", err)
		}

		sendOrder(ctx, order)
		if err != nil {
			log.Fatalf("Error sending order for customer %s with error: %v", order.LoyaltyId, err)
		} else {
			ordersCreated++
		}

		if ordersCreated >= numOrders && numOrders != -1 {
			break
		}
	}

}

func populateConstants() {
	// get store id
	if sid := os.Getenv("STORE_ID"); sid != "" {
		storeId = sid
	}

	// get max item quantity
	if miq := os.Getenv("MAX_ITEM_QUANTITY"); miq != "" {
		maxItemQuantity, _ = strconv.Atoi(miq)
	}

	// get max unique items per order
	if mui := os.Getenv("MAX_UNIQUE_ITEMS_PER_ORDER"); mui != "" {
		maxUniqueItemsPerOrder, _ = strconv.Atoi(mui)
	}

	// get min seconds to place order
	if msto := os.Getenv("MIN_SEC_TO_PLACE_ORDER"); msto != "" {
		minSecondsToPlaceOrder, _ = strconv.Atoi(msto)
	}

	// get max seconds to place order
	if msto := os.Getenv("MAX_SEC_TO_PLACE_ORDER"); msto != "" {
		maxSecondsToPlaceOrder, _ = strconv.Atoi(msto)
	}

	// get min seconds between orders
	if msbo := os.Getenv("MIN_SEC_BETWEEN_ORDERS"); msbo != "" {
		minSecondsBetweenOrders, _ = strconv.Atoi(msbo)
	}

	// get max seconds between orders
	if msbo := os.Getenv("MAX_SEC_BETWEEN_ORDERS"); msbo != "" {
		maxSecondsBetweenOrders, _ = strconv.Atoi(msbo)
	}

	// get number of orders
	if no := os.Getenv("NUM_ORDERS"); no != "" {
		numOrders, _ = strconv.Atoi(no)
	}

}

// creates a new customer order
func createOrder(products []Product) (CustomerOrder, error) {
	// get random customer from customers array
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	customer := customers[rng.Intn(len(customers)-1)]
	//customer := customers[75]

	log.Printf("Customer %v %v  with loyalty id %d is placing an order.", customer.FirstName, customer.LastName, customer.ID)

	// Get total number of menu items (t) and ramdomly choose a number of them to order (n)
	numProducts := len(products)
	numOrderItems := rng.Intn(min(numProducts, maxUniqueItemsPerOrder))

	var orderItems []OrderItem
	//fill orderitems array with random products that cannot repeat
	for i := 0; i < numOrderItems; i++ {

		// randomly choose non-duplicated product
		product := products[rng.Intn(numProducts-1)]
		log.Printf("Product %v %v  with product id %d is being added to the order.", product.ProductName, product.Description, product.ProductID)
		for ok := true; ok; ok = slices.ContainsFunc(orderItems,
			func(orderItem OrderItem) bool {
				//log.Printf("Checking if product Id %d is equal to productId %d in order", product.ProductID, orderItem.ProductID)
				return orderItem.ProductID == product.ProductID // compare current order item product id with Product Id
			}) {
			product = products[rng.Intn(numProducts-1)]
		}

		orderItems = append(orderItems, OrderItem{
			ProductID: product.ProductID,
			Quantity:  rng.Intn(maxItemQuantity) + 1,
		})

	}

	// create customer order
	order := CustomerOrder{
		StoreID:    storeId,
		FirstName:  customer.FirstName,
		LastName:   customer.LastName,
		LoyaltyId:  fmt.Sprint((customer.ID)),
		OrderItems: orderItems,
	}

	return order, nil

}

func getProducts(ctx context.Context) ([]Product, error) {

	response, err := client.InvokeMethod(ctx, OrderServiceDaprId, "product", "get")
	if err != nil {
		log.Printf("error calling service: %v", err)
		return nil, err
	}

	var products []Product
	err = json.Unmarshal(response, &products)
	if err != nil {
		log.Printf("error unmarshalling response: %v", err)
		return nil, err
	}

	log.Printf("Loaded %d products\n", len(products))

	return products, nil
}

func sendOrder(ctx context.Context, order CustomerOrder) (string, error) {
	// Create dapr request content from order
	data, err := json.Marshal(order)
	if err != nil {
		log.Printf("error marshalling order: %v", err)
		return "", err
	}

	log.Printf("Sending order: %v", string(data))

	content := &dapr.DataContent{
		ContentType: "application/json",
		Data:        data,
	}

	// Call order service
	response, err := client.InvokeMethodWithContent(ctx, OrderServiceDaprId, "order", "post", content)
	if err != nil {
		sleepTime := rand.Intn((maxSecondsToPlaceOrder - minSecondsToPlaceOrder) + minSecondsToPlaceOrder)
		time.Sleep(time.Duration(sleepTime) * time.Second)
		log.Printf("Error calling service: %v", err)
		return "", err
	}

	log.Printf("Order %s completed!", string(response))

	return string(response), nil
}
