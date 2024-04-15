RELEASE_VERSION  = 0.7.3
SERVICE_NAME = receipt-generation-service
IMAGE_NAME = order-system/receipt-generation
REGISTRY=rochabr
APP_PORT=5300
DAPR_HTTP_PORT=5380
DAPR_GRPC_PORT=5301

.PHONY: tidy
tidy: ## Updates the go modules and vendors all dependencies 
	go mod tidy



.PHONY: run
run: tidy ## Runs uncompiled code in Dapr
	dapr run \
		--app-id $(SERVICE_NAME) \
		--app-port ${APP_PORT} \
		--app-protocol http \
		--log-level debug \
		--dapr-http-port ${DAPR_HTTP_PORT} \
		--dapr-grpc-port ${DAPR_GRPC_PORT} \
		--resources-path ../../components/local \
		go run ./cmd/api

.PHONY: build_image
build_image: tidy ## Builds local release binary
	docker buildx build --no-cache --platform=linux/amd64,linux/arm64 -t ${REGISTRY}/${IMAGE_NAME}:${RELEASE_VERSION} --push .


