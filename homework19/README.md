# Отчёт по практической работе №19
## Структурированное логирование: zap + request-id

---

## 1. Запуск проекта

### Команды запуска

**Терминал 1 — Auth service:**
```bash
cd services/auth
AUTH_PORT=8081 go run ./cmd/auth
```

**Терминал 2 — Tasks service:**
```bash
cd services/tasks
TASKS_PORT=8082 AUTH_BASE_URL=http://localhost:8081 go run ./cmd/tasks
```

---

## 2. Выбранный логгер — zap

Использован `go.uber.org/zap`.

zap выбран по двум причинам.
Во-первых, он пишет JSON без дополнительной настройки — каждое лог-событие сразу является структурированной записью, удобной для парсинга в системах сбора логов. Во-вторых, zap не использует рефлексию и работает быстрее logrus, что важно при высоких нагрузках.

Логгер создаётся один раз в `shared/logger/logger.go` и прокидывается через dependency injection в роутеры и обработчики.

---

## 3. Стандарт полей логов

Каждое лог-событие содержит следующие поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `ts` | float | Unix timestamp (добавляется zap автоматически) |
| `level` | string | Уровень: `info`, `warn`, `error` |
| `msg` | string | Описание события |
| `request_id` | string | ID запроса из `X-Request-ID` или сгенерированный UUID |
| `method` | string | HTTP метод: `GET`, `POST`, и т.д. |
| `path` | string | Путь запроса: `/v1/tasks` |
| `status` | int | HTTP статус ответа |
| `duration_ms` | int | Длительность обработки в миллисекундах |
| `has_auth` | bool | Наличие заголовка `Authorization` (без значения) |


---

## 4. Примеры лог-событий

### Успешный запрос — список задач

```bash
curl -i http://localhost:8082/v1/tasks \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: pz19-001"
```

Лог Tasks service:
```json
{
    "level":"info",
    "ts":1772986609.0902636,
    "caller":"middleware/accesslog.go:34",
    "msg":"request completed",
    "request_id":"pz19-001",
    "method":"GET",
    "path":"/v1/tasks",
    "status":200,
    "duration_ms":15,
    "has_auth":true
} 
```

Лог Auth service (вызов verify):
```json
{
    "level":"info",
    "ts":1772986609.0892274,
    "caller":"middleware/accesslog.go:34",
    "msg":"request completed",
    "request_id":"pz19-001",
    "method":"GET",
    "path":"/v1/auth/verify",
    "status":200,
    "duration_ms":0,
    "has_auth":true
}
```

---

### Запрос с ошибкой — невалидный токен (401)

```bash
curl -i http://localhost:8082/v1/tasks \
  -H "Authorization: Bearer wrong-token" \
  -H "X-Request-ID: pz19-003"
```

Лог Tasks service:
```json
{
    "level":"info",
    "ts":1772986790.1813898,
    "caller":"middleware/accesslog.go:34",
    "msg":"request completed",
    "request_id":"pz19-003",
    "method":"GET",
    "path":"/v1/tasks",
    "status":401,
    "duration_ms":8,
    "has_auth":true
} 
```

Лог Auth service:
```json
{
    "level":"info",
    "ts":1772986790.1808279,
    "caller":"middleware/accesslog.go:34",
    "msg":"request completed",
    "request_id":"pz19-003",
    "method":"GET",
    "path":"/v1/auth/verify",
    "status":401,
    "duration_ms":0,
    "has_auth":true
}
```

---

### Межсервисный вызов — создание задачи с прокидыванием request-id

```bash
curl -i -X POST http://localhost:8082/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: pz19-002" \
  -d '{"title":"Logs","description":"Implement zap/logrus","due_date":"2026-01-12"}'
```

Лог Auth service — verify вызван Tasks с тем же `pz19-002`:
```json
{
    "level":"info",
    "ts":1772986828.4843075,
    "caller":"middleware/accesslog.go:34",
    "msg":"request completed",
    "request_id":"pz19-002",
    "method":"GET",
    "path":"/v1/auth/verify",
    "status":200,
    "duration_ms":0,
    "has_auth":true
}
```

Лог Tasks service — итоговый ответ клиенту:
```json
{
    "level":"info",
    "ts":1772986828.4848332,
    "caller":"middleware/accesslog.go:34",
    "msg":"request completed",
    "request_id":"pz19-002",
    "method":"POST",
    "path":"/v1/tasks",
    "status":201,
    "duration_ms":8,
    "has_auth":true
}
```

Оба лог-события содержат `"request_id": "pz19-002"` — по этому полю можно найти всю цепочку запроса сразу в двух сервисах.
