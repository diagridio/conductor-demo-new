RELEASE_VERSION  = 0.7.4
SERVICE_NAME = virtual-worker
IMAGE_NAME = virtual-worker
REGISTRY=rochabr
APP_PORT=5900
DAPR_HTTP_PORT=5980


.PHONY: tidy
tidy: ## Updates the go modules and vendors all dependencies 
	go mod tidy

.PHONY: run
run: tidy ## Runs uncompiled code in Dapr
	dapr run \
		--app-id $(SERVICE_NAME) \
		--app-protocol http \
		--app-port ${APP_PORT} \
		--dapr-http-port ${DAPR_HTTP_PORT} \
		--resources-path ../../components/local \
		go run ./cmd/api

.PHONY: build_image
build_image: tidy ## Builds local release binary
	docker buildx build --no-cache --platform=linux/amd64,linux/arm64 -t ${REGISTRY}/${IMAGE_NAME}:${RELEASE_VERSION} --push .

