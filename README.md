# Лабораторная работа №13 — Мультиагентные системы

**Студент:** Ковыльников Роман Александрович  
**Вариант:** 10 — Обработка заказов в ресторане  
**Сложность:** Средняя (задания 1–10)

## Стек

| Компонент | Технология |
|-----------|-----------|
| Агенты (4 типа) | Go 1.22 + nats.go |
| Оркестратор | Python 3.12 + nats-py + asyncio |
| REST API | FastAPI + uvicorn |
| Брокер сообщений | NATS 2.10 |
| Инфраструктура | Docker Compose |

## Быстрый старт (Docker)

```bash
mkdir -p logs
docker compose up --build
```

API доступен на http://localhost:8000  
NATS monitoring: http://localhost:8222

### Пример запроса

```bash
curl -X POST http://localhost:8000/orders \
  -H "Content-Type: application/json" \
  -d '{
    "table_id": 3,
    "customer_name": "Иван",
    "items": [
      {"name": "burger", "qty": 2, "price": 250},
      {"name": "soup",   "qty": 1, "price": 180}
    ]
  }'
```

## Тесты

### Go

```bash
cd agent
go test ./...
```

### Python

```bash
cd orchestrator
pip install -r requirements.txt
pytest
```

## Локальный запуск агентов (без Docker)

Нужен запущенный NATS: `docker run -p 4222:4222 nats:2.10-alpine`

```bash
cd agent
go run . --type order
go run . --type kitchen --log logs/kitchen-1.log &
go run . --type kitchen --log logs/kitchen-2.log &
go run . --type table
go run . --type delivery
```

## Структура репозитория

```
agent/          — Go-агенты (один бинарник, тип выбирается флагом --type)
orchestrator/   — Python-оркестратор + FastAPI REST API + pytest
docs/           — архитектурная документация и диаграммы
docker-compose.yml
```

## Архитектура

См. [docs/architecture.md](docs/architecture.md)

## API

| Метод | Путь | Описание |
|-------|------|----------|
| POST | /orders | Принять заказ |
| POST | /kitchen | Отправить заказ на кухню |
| POST | /tables | Обновить статус стола |
| POST | /delivery | Назначить доставку |
| GET | /stats | Счётчик обработанных задач |
