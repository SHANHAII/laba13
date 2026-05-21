import asyncio
import json
import logging
import uuid
from typing import Optional

import nats

AGENT_SUBJECTS = {
    "order":    "restaurant.order.new",
    "kitchen":  "restaurant.kitchen.new",
    "table":    "restaurant.table.update",
    "delivery": "restaurant.delivery.new",
}

MAX_RETRIES = 3
DEFAULT_TIMEOUT = 30


def _make_logger(name: str, log_file: Optional[str] = None) -> logging.Logger:
    logger = logging.getLogger(name)
    logger.setLevel(logging.DEBUG)
    fmt = logging.Formatter("%(asctime)s [%(name)s] %(levelname)s %(message)s")

    ch = logging.StreamHandler()
    ch.setFormatter(fmt)
    logger.addHandler(ch)

    if log_file:
        fh = logging.FileHandler(log_file)
        fh.setFormatter(fmt)
        logger.addHandler(fh)

    return logger


class AgentOrchestrator:
    def __init__(
        self,
        nats_url: str = "nats://localhost:4222",
        log_file: Optional[str] = None,
    ) -> None:
        self.nats_url = nats_url
        self.nc: Optional[nats.NATS] = None
        self._pending: dict[str, asyncio.Future] = {}
        self._processed = 0
        self.logger = _make_logger("orchestrator", log_file)

    async def connect(self) -> None:
        self.nc = await nats.connect(self.nats_url)
        await self.nc.subscribe("restaurant.tasks.completed", cb=self._on_result)
        self.logger.info("connected to NATS at %s", self.nats_url)

    async def close(self) -> None:
        if self.nc:
            await self.nc.close()

    async def _on_result(self, msg) -> None:
        try:
            result = json.loads(msg.data.decode())
            task_id = result.get("task_id")
            if task_id in self._pending:
                self._pending[task_id].set_result(result)
                del self._pending[task_id]
        except Exception as exc:
            self.logger.error("failed to handle result message: %s", exc)

    async def send_task(
        self,
        task_type: str,
        payload: dict,
        timeout: int = DEFAULT_TIMEOUT,
        retries: int = MAX_RETRIES,
    ) -> dict:
        subject = AGENT_SUBJECTS.get(task_type)
        if subject is None:
            raise ValueError(f"unknown task type: {task_type!r}")

        last_error: Optional[Exception] = None

        for attempt in range(1, retries + 1):
            task_id = str(uuid.uuid4())
            task = {"id": task_id, "type": task_type, "payload": json.dumps(payload)}

            loop = asyncio.get_running_loop()
            future: asyncio.Future = loop.create_future()
            self._pending[task_id] = future

            await self.nc.publish(subject, json.dumps(task).encode())
            self.logger.info(
                "task %s sent (type=%s, attempt=%d/%d)", task_id, task_type, attempt, retries
            )

            try:
                result = await asyncio.wait_for(future, timeout)
                if result.get("success"):
                    self._processed += 1
                    self.logger.info(
                        "task %s succeeded, total processed: %d", task_id, self._processed
                    )
                    return result
                last_error = RuntimeError(result.get("error", "agent returned failure"))
                self.logger.error("task %s failed: %s", task_id, last_error)
            except asyncio.TimeoutError:
                self._pending.pop(task_id, None)
                last_error = TimeoutError(f"task {task_id} timed out after {timeout}s")
                self.logger.error("task %s timed out (attempt %d/%d)", task_id, attempt, retries)

        raise last_error or RuntimeError("all retries exhausted")

    @property
    def processed_count(self) -> int:
        return self._processed
