# Практическое занятие №3

## Тема: Реализация простого HTTP-сервера на стандартной библиотеке net/http. Обработка запросов GET/POST

**Студент:** Наумов А.Е.
**Группа:** ЭФМО-01-25

### Запуск
1. Клонирование репозитория
```bash
git clone git@github.com:d1vide/golang.git
```
2. Переход в директорию проекта 
```bash
cd homework3/pz3-http/
```
3. Запуск приложения
   3.1. Запуск через `go run`
   ```bash
   go run ./cmd/server/
   ```
# Описание проекта:
Простой HTTP API сервер на Go, предоставляющий базовые эндпоинты.

## Примеры ответов на запросы
1. Создание таски

![alt text](screenshots/image-1.png)

2. Просмотр списка тасок

![alt text](screenshots/image-2.png)

3. Просмотр конкретной таски

![alt text](screenshots/image-3.png)

4. Изменение поля done таски

![alt text](screenshots/image-4.png)

5. Проверка изменения поля

![alt text](screenshots/image-5.png)

6. Удаление таски

![alt text](screenshots/image-6.png)

7. Проверка удаления тасок

![alt text](screenshots/image-7.png)

8. Проверка удаления таски

![alt text](screenshots/image.png)

9. Проверка возможности done: false для таски

![alt text](screenshots/image-8.png)

(копипасты для curl запросов в `requests.md`)

## CORS
С помощью middleware реализовано добавление заголовка для проверка CORS

Тестирование:
Написан index.html делающий запрос на /health
В коде middleware захардкожено значение хоста для страницы:
`w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")`

В такой реализации результат будет корректно возвращаться только для адреса `127.0.0.1` но не для `localhost`

`127.0.0.1`

![alt text](screenshots/image3.png)

`localhost`

![alt text](screenshots/image1.png)

![alt text](screenshots/image4.png)

## Graceful shutdown
Обеспечивает плавное завершение работы сервера без обрыва текущих соединений

Обрабатываемые сигналы:
- SIGINT - прерывание от пользователя
- SIGTERM - запрос на завершение от системы

Контекст задает таймаут 30 секунд на graceful shutdown:
- Сервер перестает принимать новые запросы
- Ждет завершения активных обработчиков (30 сек)
- Если не успевает - принудительно закрывает оставшиеся соединения


![alt text](screenshots/image2.png)
