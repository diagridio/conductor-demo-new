RELEASE_VERSION  = 0.6.1
SERVICE_NAME = receipt-generation-service
IMAGE_NAME = order-system/receipt-generation
REGISTRY=rochabr
APP_PORT=5300
DAPR_HTTP_PORT=5380
DAPR_GRPC_PORT=5301

.PHONY: setup
setup: # Manages dependencies using requirements.txt
	pip3 install -r requirements.txt

.PHONY: run
run: setup ## Runs uncompiled code in Dapr
	dapr run \
		--app-id $(SERVICE_NAME) \
		--app-port ${APP_PORT} \
		--app-protocol http \
		--log-level debug \
		--dapr-http-port ${DAPR_HTTP_PORT} \
		--dapr-grpc-port ${DAPR_GRPC_PORT} \
		--resources-path ../../components/local \
		python3 ./api/main.py

.PHONY: build_image
build_image: tidy ## Builds local release binary
	docker buildx build --no-cache --platform=linux/amd64,linux/arm64 -t ${REGISTRY}/${IMAGE_NAME}:${RELEASE_VERSION} --push .


