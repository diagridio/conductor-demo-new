RELEASE_VERSION  = 0.7.5
IMAGE_NAME = worker-actor
SERVICE_NAME = worker-actor
REGISTRY=rochabr
APP_PORT=5500
DAPR_HTTP_PORT=5580

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

