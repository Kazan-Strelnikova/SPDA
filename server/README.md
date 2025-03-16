# Event planner

## Пререквизиты

Перед установкой и запуском проекта необходимо, чтобы на устройстве был:

- docker
- docker-compose
- go (версия указанная [здесь](/server/go.mod))
- Taskfile (опционально, но сильно упростит процесс, [инструкция по установке](https://taskfile.dev/installation/))

## База данных

Для успешного запуска базы данных нужно создать `.env` файл, и положить туда следующие секреты:

- POSTGRES_USER
- POSTGRES_PASSWORD
- POSTGRES_DB
- POSTGRES_PORT

Например,

```env
POSTGRES_USER=planner-user
POSTGRES_PASSWORD=StrongRandomPassword123!
POSTGRES_DB=planner-db
POSTGRES_PORT=5433
```

В проекте уже имеются настроенные миграции для нашей схемы базы данных. Пример алгоритма по запоску дб:

1. Убедиться, что есть все необходимые переменные окружения
2. Запустить команду `docker-compose up -d` или `task docker-up`
3. Дождаться полной иницализации контейнера
4. Установить нужные модули: `go mod download`
5. Запустить команду `task migrate-up`. Если нет утилиты taskfile, можете подсмотреть соответствующую команду в [taskfile.yml](/server/Taskfile.yml)

## Запуск приложения 

В данной главе описан алгоритм для успешного запуска бэкенда. 

### Переменные окружения 

Вот пример полноценного .env необходимого для запуска приложения

```env
POSTGRES_USER=planner-user
POSTGRES_PASSWORD=StrongRandomPassword123!
POSTGRES_DB=planner_db
POSTGRES_PORT=5433
POSTGRES_HOST=localhost
SERVER_PORT=8080
RW_TIMEOUT=5s
IDLE_TIMEOUT=15m
```

### Старт сервиса

В данный момент приложение еще не докеризовано. Перед его запуском пройдите главу связанную с базой данных. После этого используйте команду `task run`,
сервер станет доступен на порту, указанном в переменной SERVER_PORT. Как только приложение полностью запустится, эндпоинт `/ping` начнёт возвращать код 200.