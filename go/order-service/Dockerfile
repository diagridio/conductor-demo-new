FROM alpine:latest

RUN mkdir /app

COPY orderServiceApp /app

# Copying product files to container
COPY /cmd/api/product-definitions /app

EXPOSE 80

CMD [ "/app/orderServiceApp" ]
