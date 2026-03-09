# PZ17 — Микросервисная система: Auth + Tasks

## Переменные окружения и порты

| Сервис | Переменная      | Значение по умолчанию       |
|--------|-----------------|----------------------------|
| Auth   | `AUTH_PORT`     | `8081`                     |
| Tasks  | `TASKS_PORT`    | `8082`                     |
| Tasks  | `AUTH_BASE_URL` | `http://localhost:8081`    |

---

## Auth service (`localhost:8081`)

### POST /v1/auth/login

Упрощённая авторизация. Принимает `username`/`password`, возвращает токен.

**Request**
```json
{
  "username": "student",
  "password": "student"
}
```

**Response 200**
```json
{
  "access_token": "demo-token",
  "token_type": "Bearer"
}
```

**Коды ответов**
| Код | Описание                        |
|-----|---------------------------------|
| 200 | Успешно, токен выдан            |
| 400 | Неверный формат запроса         |
| 401 | Неверные учётные данные         |

**curl**
```bash
curl -s -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: req-001" \
  -d '{"username":"student","password":"student"}'
```

---

### GET /v1/auth/verify

Проверяет валидность токена из заголовка `Authorization`.

**Headers**
```
Authorization: Bearer demo-token
X-Request-ID: req-002   (опционально)
```

**Response 200**
```json
{
  "valid": true,
  "subject": "student"
}
```

**Response 401**
```json
{
  "valid": false,
  "error": "unauthorized"
}
```

**Коды ответов**
| Код | Описание                        |
|-----|---------------------------------|
| 200 | Токен валиден                   |
| 401 | Токен невалиден или отсутствует |

**curl**
```bash
curl -i http://localhost:8081/v1/auth/verify \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: req-002"
```

---

## Tasks service (`localhost:8082`)

Все эндпоинты требуют заголовок `Authorization: Bearer <token>`.  
Токен проверяется через Auth service перед каждой операцией.

### POST /v1/tasks — создать задачу

**Request**
```json
{
  "title": "Read lecture",
  "description": "Prepare notes",
  "due_date": "2026-01-10"
}
```

**Response 201**
```json
{
  "id": "t_001",
  "title": "Read lecture",
  "description": "Prepare notes",
  "due_date": "2026-01-10",
  "done": false
}
```

**curl**
```bash
curl -i -X POST http://localhost:8082/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: req-003" \
  -d '{"title":"Do PZ17","description":"split services","due_date":"2026-01-10"}'
```

---

### GET /v1/tasks — список задач

**Response 200**
```json
[
  {"id":"t_001","title":"Read lecture","done":false},
  {"id":"t_002","title":"Do practice","done":true}
]
```

**curl**
```bash
curl -i http://localhost:8082/v1/tasks \
  -H "Authorization: Bearer demo-token"
```

---

### GET /v1/tasks/{id} — получить задачу

**Response 200**
```json
{
  "id": "t_001",
  "title": "Read lecture",
  "description": "Prepare notes",
  "due_date": "2026-01-10",
  "done": false
}
```

**curl**
```bash
curl -i http://localhost:8082/v1/tasks/t_001 \
  -H "Authorization: Bearer demo-token"
```

---

### PATCH /v1/tasks/{id} — обновить задачу

Все поля опциональны (частичное обновление).

**Request**
```json
{
  "title": "Read lecture (updated)",
  "done": true
}
```

**Response 200** — обновлённая задача.

**curl**
```bash
curl -i -X PATCH http://localhost:8082/v1/tasks/t_001 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-token" \
  -d '{"done":true}'
```

---

### DELETE /v1/tasks/{id} — удалить задачу

**Response 204** — тело отсутствует.

**curl**
```bash
curl -i -X DELETE http://localhost:8082/v1/tasks/t_001 \
  -H "Authorization: Bearer demo-token"
```

---

### Коды ответов Tasks

| Код | Описание                                          |
|-----|---------------------------------------------------|
| 200 | Успешно                                           |
| 201 | Задача создана                                    |
| 204 | Задача удалена                                    |
| 400 | Неверные данные запроса                           |
| 401 | Токен невалиден (Auth вернул 401)                 |
| 404 | Задача не найдена                                 |
| 503 | Auth service недоступен (fail-closed)             |

---

### Проверка без токена (должно вернуть 401)

```bash
curl -i http://localhost:8082/v1/tasks \
  -H "X-Request-ID: req-004"
```