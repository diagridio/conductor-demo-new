from fastapi import FastAPI, status, HTTPException
from fastapi.responses import JSONResponse

import uvicorn

import os
import json
import logging


import models

from dapr.clients import DaprClient

app = FastAPI()

logging.basicConfig(level = logging.INFO)
BINDING_OPERATION = 'create' 
BINDING_NAME = 'oms.binding.receipt'

@app.get("/healthz")
def healthz():
    return JSONResponse(status_code=status.HTTP_200_OK)

@app.get("/ready")
def ready():
    return JSONResponse(status_code=status.HTTP_200_OK)

@app.post("/receipt")
def handleReceipt(orderSummary: models.OrderSummary):
    logging.info(f'Received order: {orderSummary.orderId}')
    
    resp = saveReceipt(orderSummary)
    
    return resp

def saveReceipt(order: models.OrderSummary):
    with DaprClient() as d:

        logging.info(f'Saving receipt for order: {order.orderId}')

        # Create a typed message with content type and body
        binding_key = {
            'key': order.orderId
        }
        binding_data = {
            'orderId': order.orderId,
            'storeId': order.storeId,
            'firstName': order.firstName,
            'lastName': order.lastName,
            'loyaltyId': order.loyaltyId,
        }

        # Induce error when using Redis. No key "key" in the request
        #req_data = {'receiptName': order.orderId}

        # Invoke binding
        try:
            # Insert order using Dapr output binding via HTTP Post
            resp = d.invoke_binding(BINDING_NAME, BINDING_OPERATION, json.dumps(binding_data), binding_key)
            logging.info(f'Saved receipt for order: {order.orderId}')
            
            return order.orderId
        except Exception as e:
            print(e, flush=True)
            logging.info(f'Saved receipt for order: {order.orderId}')

            raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail="Error saving receipt")


if __name__ == "__main__":
    port = int(os.getenv('APP_PORT'))
    if port is None:
        port = 5300
    uvicorn.run(app, port=port)
