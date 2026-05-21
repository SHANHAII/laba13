import asyncio
import json
from unittest.mock import AsyncMock

import pytest

from orchestrator import AgentOrchestrator


@pytest.fixture
def orch():
    o = AgentOrchestrator()
    o.nc = AsyncMock()
    return o


# ---------- helpers ----------

def _resolve_on_publish(orch: AgentOrchestrator, success: bool, error: str = ""):
    """Возвращает side_effect для nc.publish, который сразу резолвит pending future."""
    async def fake_publish(subject, data):
        payload = json.loads(data)
        task_id = payload["id"]
        future = orch._pending.get(task_id)
        if future and not future.done():
            result = {"task_id": task_id, "success": success, "output": "", "error": error}
            future.set_result(result)

    return fake_publish


# ---------- tests ----------

@pytest.mark.asyncio
async def test_send_task_success(orch):
    orch.nc.publish = AsyncMock(side_effect=_resolve_on_publish(orch, success=True))

    result = await orch.send_task("order", {"table_id": 1, "items": []})

    assert result["success"] is True
    assert orch.processed_count == 1


@pytest.mark.asyncio
async def test_send_task_timeout(orch):
    orch.nc.publish = AsyncMock()  # future никогда не резолвится

    with pytest.raises(TimeoutError):
        await orch.send_task("order", {}, timeout=0.05, retries=1)


@pytest.mark.asyncio
async def test_send_task_retry_on_agent_failure(orch):
    call_count = 0

    async def fake_publish(subject, data):
        nonlocal call_count
        call_count += 1
        payload = json.loads(data)
        task_id = payload["id"]
        future = orch._pending.get(task_id)
        if future and not future.done():
            ok = call_count >= 3  # первые 2 — ошибка, 3-й — успех
            future.set_result({"task_id": task_id, "success": ok, "output": "", "error": "fail"})

    orch.nc.publish = AsyncMock(side_effect=fake_publish)

    result = await orch.send_task("order", {}, retries=3)

    assert result["success"] is True
    assert call_count == 3


@pytest.mark.asyncio
async def test_send_task_all_retries_exhausted(orch):
    orch.nc.publish = AsyncMock(side_effect=_resolve_on_publish(orch, success=False, error="boom"))

    with pytest.raises(RuntimeError, match="boom"):
        await orch.send_task("order", {}, retries=3)


@pytest.mark.asyncio
async def test_unknown_task_type(orch):
    with pytest.raises(ValueError, match="unknown task type"):
        await orch.send_task("unknown_type", {})


@pytest.mark.asyncio
async def test_processed_count_accumulates(orch):
    orch.nc.publish = AsyncMock(side_effect=_resolve_on_publish(orch, success=True))

    await orch.send_task("order", {})
    await orch.send_task("kitchen", {})
    await orch.send_task("delivery", {})

    assert orch.processed_count == 3


@pytest.mark.asyncio
async def test_on_result_ignores_unknown_task_id(orch):
    class FakeMsg:
        data = json.dumps({"task_id": "nonexistent-id", "success": True}).encode()

    # не должно бросить исключение
    await orch._on_result(FakeMsg())
