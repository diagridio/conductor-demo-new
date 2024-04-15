import pydantic
import pydantic
from typing import List
from datetime import datetime

class OrderItemSummary(pydantic.BaseModel):
    productId: int
    productName: str
    quantity: int
    unitCost: float
    unitPrice: float
    imageUrl: str

class OrderSummary(pydantic.BaseModel):
    orderId: str
    orderDate: datetime
    orderCompletedDate: datetime = None
    storeId: str
    firstName: str
    lastName: str
    loyaltyId: str
    orderItems: List[OrderItemSummary]
    orderTotal: float