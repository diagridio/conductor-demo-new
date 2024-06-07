RELEASE_VERSION  = 0.7.4
IMAGE_NAME = basic-actor-demos
SERVICE_NAME = basic-actor-demos
REGISTRY=rochabr
APP_PORT=7100
DAPR_HTTP_PORT=7180
DAPR_GRPC_PORT=7101

.PHONY: clean
clean: ## Cleans the project
	dotnet clean


.PHONY: run
run: clean ## Runs uncompiled code in Dapr
	dapr run \
		--app-id $(SERVICE_NAME) \
		--app-port ${APP_PORT} \
		--dapr-http-port ${DAPR_HTTP_PORT} \
		--resources-path ../../components/local \
		dotnet run

.PHONY: build_image
build_image: clean ## Builds local release binary
	docker buildx build --no-cache --platform=linux/amd64,linux/arm64 -t ${REGISTRY}/${IMAGE_NAME}:${RELEASE_VERSION} --push .

