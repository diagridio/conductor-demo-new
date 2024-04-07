REGISTRY=rochabr

# ## up: starts all containers in the background without forcing build
# up:
# 	@echo "Starting Docker images..."
# 	docker-compose up -d
# 	@echo "Docker images started!"

# ## up_build: stops docker-compose (if running), builds all projects and starts docker compose
# up_build: build_order_service build_receipt_generation_service build_loyalty_service build_virtual_customer build_make_line_service
# 	@echo "Stopping docker images (if running...)"
# 	docker-compose down
# 	@echo "Building (when required) and starting docker images..."
# 	docker-compose up --build -d
# 	@echo "Docker images built and started!"

# ## down: stop docker compose
# down:
# 	@echo "Stopping docker compose..."
# 	docker-compose down
# 	@echo "Done!"

## build_loyalty_service: builds the loyalty service binary container and push it to the registry
build_loyalty_service:
	@echo "Building loyalty service binary..."
	docker buildx build --no-cache --platform=linux/amd64,linux/arm64 -t ${REGISTRY}/loyalty-service:latest --push ./../go/loyalty-service
	@echo "Done!"

## build_make_line_service: builds the make-line service binary container and push it to the registry
build_make_line_service:
	@echo "Building make-line service binary..."
	docker buildx build --no-cache --platform=linux/amd64,linux/arm64 -t ${REGISTRY}/make-line-service:latest --push ./../go/make-line-service	
	@echo "Done!"

## build_receipt_generation_service: builds the receipt generation service binary binary container and push it to the registry
build_receipt_generation_service:
	@echo "Building receipt generation service binary..."
	docker buildx build --no-cache --platform=linux/amd64,linux/arm64 -t ${REGISTRY}/receipt-generation-service:latest --push ./../go/receipt-generation-service	
	@echo "Done!"

## build_order_service: builds the order service binary binary container and push it to the registry
build_order_service:
	@echo "Building order service binary..."
	docker buildx build --no-cache --platform=linux/amd64,linux/arm64 -t ${REGISTRY}/order-service:latest --push ./../go/order-service
	@echo "Done!"

## build_order_processor_service: builds the order processor service container and push it to the registry
build_order_processor_service:
	@echo "Building order processor service binary..."
	docker buildx build --no-cache --platform=linux/amd64,linux/arm64 -t ${REGISTRY}/order-processor-service:latest --push ./../dotnet/OrderProcessorService
	@echo "Done!"

## build_virtual_customer_service: builds the virtual customer binary container and push it to the registry
build_virtual_customer:
	@echo "Building virtual customer binary..."
	docker buildx build --no-cache --platform=linux/amd64,linux/arm64 -t ${REGISTRY}/virtual-customer:latest --push ./../go/virtual-customer
	@echo "Done!"

# build and push all containers
build_push: build_order_processor_service build_order_service build_loyalty_service build_make_line_service build build_virtual_customer build_receipt_generation_service
	@echo "Building and pushing images to the registry..."
	@echo "Done!"
