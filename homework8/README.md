# Практическое занятие №8

## Тема: Работа с MongoDB: подключение, создание коллекции, CRUD-операции

**Студент:** Наумов А.Е.
**Группа:** ЭФМО-01-25


## Подготовка окружения

Контейнер с MongoDB
```bash
docker compose up -d
```

## Примеры запросов

- POST `/api/v1/notes`

```bash
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Accept: application/json" \
  -d '{"title":"MyTitle1","content":"MyContent1"}'
```

В MongoDB коллекции notes:
![alt text](screenshots/image.png)

- GET `/api/v1/notes`
```bash
curl "http://localhost:8080/api/v1/notes"
```
![alt text](screenshots/image-1.png)


- GET `/api/v1/notes/{id}`

![alt text](screenshots/image-2.png)

- PATCH `/api/v1/notes/{id}`

![alt text](screenshots/image-3.png)

![alt text](screenshots/image-4.png)

- DELETE `/api/v1/notes/{id}`

![alt text](screenshots/image-5.png)

Запись успешно удалена


### Дополнительные задания

- GET `api/v1/notes?q=MyContent1`

![alt text](screenshots/image-6.png)

- GET `/api/v1/notes/stats`

![alt text](screenshots/image-7.png)

- TTL
Объекты создаются с TTL (24 часа по умолчанию)

![alt text](screenshots/image-8.png)