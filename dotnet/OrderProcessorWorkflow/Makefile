RELEASE_VERSION  = 0.7.3
IMAGE_NAME = order-processor-workflow
SERVICE_NAME = order-processor-workflow
REGISTRY=rochabr
APP_PORT=5700
DAPR_HTTP_PORT=5780
DAPR_GRPC_PORT=5701

.PHONY: clean
clean: ## Cleans the project
	dotnet clean


.PHONY: run
run: clean ## Runs uncompiled code in Dapr
	dapr run \
		--app-id $(SERVICE_NAME) \
		--app-port ${APP_PORT} \
		--app-protocol http \
		--dapr-http-port ${DAPR_HTTP_PORT} \
		--dapr-grpc-port ${DAPR_GRPC_PORT} \
		--resources-path ../../components/local \
		dotnet run

.PHONY: build_image
build_image: clean ## Builds local release binary
	docker buildx build --no-cache --platform=linux/amd64,linux/arm64 -t ${REGISTRY}/${IMAGE_NAME}:${RELEASE_VERSION} --push .

