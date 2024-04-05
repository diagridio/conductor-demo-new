LOYALTY_SERVICE_BINARY=loyaltyServiceApp
MAKE_LINE_SERVICE_BINARY=makeLineServiceApp
ORDER_SERVICE_BINARY=orderServiceApp
RECEIPT_GENERATION_SERVICE_BINARY=receiptGenerationServiceApp
VIRTUAL_CUSTOMER_BINARY=virtualCustomerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_order_service build_receipt_generation_service build_loyalty_service build_virtual_customer build_make_line_service
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_loyalty_service: builds the loyalty service binary as a linux executable
build_loyalty_service:
	@echo "Building loyalty service binary..."
	cd ../go/loyalty-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOYALTY_SERVICE_BINARY} ./cmd/api
	@echo "Done!"

## build_make_line_service: builds the make-line service binary as a linux executable
build_make_line_service:
	@echo "Building make-line service binary..."
	cd ../go/make-line-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAKE_LINE_SERVICE_BINARY} ./cmd/api
	@echo "Done!"

## build_order_service: builds the order service binary as a linux executable
build_order_service:
	@echo "Building order service binary..."
	cd ../go/order-service && env GOOS=linux CGO_ENABLED=0 go build -o ${ORDER_SERVICE_BINARY} ./cmd/api
	@echo "Done!"

## build_receipt_generation_service: builds the receipt generation service binary as a linux executable
build_receipt_generation_service:
	@echo "Building receipt generation service binary..."
	cd ../go/receipt-generation-service && env GOOS=linux CGO_ENABLED=0 go build -o ${RECEIPT_GENERATION_SERVICE_BINARY} ./cmd/api
	@echo "Done!"

## build_virtual_customer_service: builds the virtual customer binary as a linux executable
build_virtual_customer:
	@echo "Building virtual customer binary..."
	cd ../go/virtual-customer && env GOOS=linux CGO_ENABLED=0 go build -o ${VIRTUAL_CUSTOMER_BINARY} ./cmd
	@echo "Done!"
