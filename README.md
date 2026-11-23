# TODO Scheduler

Веб-сервер — планировщик задач с REST API, SQLite и поддержкой повторяющихся задач.

## Выполненные задания со звёздочкой

- Порт из `TODO_PORT`
- Путь к БД из `TODO_DBFILE`
- Поиск задач (`/api/tasks?search=...`)
- Правила повторения `w` и `m`
- Аутентификация по паролю (`TODO_PASSWORD`) с JWT
- Docker-образ

## Локальный запуск

```bash
go run main.go