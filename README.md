# Demo Service

## ⚙️ Используемый стек

* **Go** (net/http, gorm)
* **PostgreSQL**
* **Kafka**
* **LRU Cache** (реализация на Go)
* **Swagger** (документация доступна по адресу `http://localhost:8081/swagger`)
* **Docker / Docker Compose**

## 🚀 Запуск проекта

### 1. Клонировать репозиторий

```bash
git clone https://github.com/ddovwll/wb-l0
```

### 2. Запустить сервисы через Docker Compose

```bash
docker-compose up --build
```

### 3. Веб-интерфейс для поиска заказов:

```
http://localhost:8081
```

### 4. Swagger-документация

```
http://localhost:8081/swagger
```

## 🧪 Тесты

Юнит-тесты запускаются на этапе сборки контейнера.
Запустить тесты вручную можно так:

```bash
go test ./... -v
```

## 🔑 Переменные окружения

Файл `.env` содержит все необходимые настройки:

```dotenv
# PostgreSQL
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=demo
POSTGRES_HOST=db
POSTGRES_PORT=5432

# Kafka
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC=orders
DLQ_TOPIC=orders_dlq
KAFKA_GROUP_ID=order-service-group
CONSUMER_WORKER_COUNT=12

# HTTP
HTTP_PORT=8081

# Cache
CACHE_CAPACITY=100
CACHE_PRELOAD=50
```

## 📂 Структура проекта

```
demoService/
 ├── docs/               # Документация swagger
 ├── src/ 
 │   ├── application/    # Бизнес-логика
 │   ├── domain/         # Модели и интерфейсы
 │   ├── infrastructure/ # Работа с БД, Kafka, кешем
 │   ├── tests/          # Юнит-тесты
 │   └── web/            # Контроллеры, роутинг, UI-шаблоны
 └── main.go
```