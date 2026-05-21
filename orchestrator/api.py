import os
from contextlib import asynccontextmanager
from typing import List

import uvicorn
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

from orchestrator import AgentOrchestrator

_nats_url = os.getenv("NATS_URL", "nats://localhost:4222")
_log_file = os.getenv("LOG_FILE", None)
orch = AgentOrchestrator(nats_url=_nats_url, log_file=_log_file)


@asynccontextmanager
async def lifespan(app: FastAPI):
    await orch.connect()
    yield
    await orch.close()


app = FastAPI(title="Restaurant MAS API", lifespan=lifespan)


# ---------- request models ----------

class OrderItemModel(BaseModel):
    name: str
    qty: int
    price: float


class OrderRequest(BaseModel):
    table_id: int
    customer_name: str
    items: List[OrderItemModel]


class KitchenRequest(BaseModel):
    order_id: str
    items: List[OrderItemModel]


class TableRequest(BaseModel):
    order_id: str
    table_id: int
    status: str  # occupied | free | reserved


class DeliveryRequest(BaseModel):
    order_id: str
    table_id: int


# ---------- endpoints ----------

@app.post("/orders")
async def create_order(req: OrderRequest):
    try:
        result = await orch.send_task("order", req.model_dump())
        return {"status": "ok", "result": result}
    except Exception as exc:
        raise HTTPException(status_code=500, detail=str(exc))


@app.post("/kitchen")
async def send_to_kitchen(req: KitchenRequest):
    try:
        result = await orch.send_task("kitchen", req.model_dump())
        return {"status": "ok", "result": result}
    except Exception as exc:
        raise HTTPException(status_code=500, detail=str(exc))


@app.post("/tables")
async def update_table(req: TableRequest):
    try:
        result = await orch.send_task("table", req.model_dump())
        return {"status": "ok", "result": result}
    except Exception as exc:
        raise HTTPException(status_code=500, detail=str(exc))


@app.post("/delivery")
async def request_delivery(req: DeliveryRequest):
    try:
        result = await orch.send_task("delivery", req.model_dump())
        return {"status": "ok", "result": result}
    except Exception as exc:
        raise HTTPException(status_code=500, detail=str(exc))


@app.get("/stats")
async def stats():
    return {"processed": orch.processed_count}


if __name__ == "__main__":
    uvicorn.run("api:app", host="0.0.0.0", port=8000, reload=False)
