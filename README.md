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

## Сборка и запуск DOCKER

docker build -t todo-scheduler:v1 .

## Запуск без пароля
docker run -d --name todo-scheduler -p 7540:7540 todo-scheduler:v1

## Запуск с паролем
docker run -d   -p 7540:7540   -e TODO_PASSWORD=12345   -v /e/learn/final_project_step1/.web/scheduler.db:/app/scheduler.db   --name todo   todo-scheduler:v1
