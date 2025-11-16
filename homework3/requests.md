## SUCCESS
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy milk"}'
```

```bash
curl -X GET http://localhost:8080/tasks
```

```bash
curl -X GET http://localhost:8080/tasks/1
```

```bash
curl -X PATCH http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"done": true}'
```

```bash
curl -X DELETE http://localhost:8080/tasks/1
```

## ERROR
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "LOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
  OOOOOOOOOOOO OOOOOOOO OOOOOOOOOO OOOOOO OOOOOOOOOOOOOOOOOOOOOOOOONG
  TEEEEEEEEEEEEEEEEEEEXTTTTTTTTTTTTTTTTTTT
  TTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTT
  TTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTT"}'
```

```bash
curl -X PATCH http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"done": false}'
```
