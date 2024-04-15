from fastapi import FastAPI, status, HTTPException
from fastapi.responses import JSONResponse
from dapr.ext.fastapi import DaprApp

import uvicorn

import os
import json
import logging
import requests

import models

from dapr.clients import DaprClient

app_port = os.getenv('APP_PORT', '5300')

app = FastAPI()
dapr_app = DaprApp(app)


logging.basicConfig(level = logging.INFO)
BINDING_OPERATION = 'create' 
BINDING_NAME = 'oms.binding.receipt'

@app.get("/healthz")
def healthz():
    port = int(os.getenv('DAPR_HTTP_PORT'))
    if port is None:
        port = 5380
    resp = requests.get(f'http://localhost:{port}/v1.0/healthz')

    if resp.status_code != 200:
        return JSONResponse({'status': 'Healthy'}, status_code=status.HTTP_200_OK)
    
    return JSONResponse({'status': 'Not found'}, status_code=status.HTTP_404_NOT_FOUND)

@app.get("/ready")
def ready():
    return JSONResponse({'status': 'Dapr is ready to go!'}, status_code=status.HTTP_200_OK)

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
            logging.info(f'Error saving receipt for order: {order.orderId}')

            raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail="Error saving receipt")


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=int(app_port))
