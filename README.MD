# MinIO Fixture Generator

Генератор тестовых файлов (текста и изображений) с тэгами для MinIO, реализован на Go.  
Позволяет быстро наполнить MinIO данными для тестирования или отладки интеграций.

---

## 🧰 Возможности

- Генерация файлов:
  - text → текстовый файл с UUID в содержимом
  - image → однотонное JPEG изображение (100×100 пикселей) с случайным цветом
- Поддержка тегов:
  - uuid
  - random_hash
  - date или datetime
  - random_date
- Опциональная очистка бакетов перед загрузкой
- Многопоточная загрузка файлов
- Опциональная отправка метаданных файлов в Kafka
- Гибкая конфигурация через `config.json`
- Управляемый прогресс репортинг через переменные окружения
- Запуск в Docker

---

## ⚙️ Установка

Поддерживается Docker или Go >= 1.23.5.

### Локальная сборка

```bash
go build -o minio-fixture-generator ./cmd/generator
```

### Docker

Собрать Docker-образ:

```bash
docker compose build
```

---

## 🚀 Запуск

Через Docker Compose:

```bash
docker compose up
```

---

## 📁 Структура проекта

```
.
├── cmd/generator
│    └── main.go
├── internal/
│    ├── config
│    ├── generator
│    ├── miniohandler
│    ├── kafka
│    └── service
├── config/config.json
├── Dockerfile
├── docker-compose.yml
├── .env
└── README.md
```

---

## ⚙️ Конфигурация

Конфигурация хранится в JSON-файле. Пример **финального конфига**:

```json
{
  "file_count": 10,
  "file_types": ["text", "image"],
  "buckets": ["test-bucket", "images"],
  "tags": {
    "entity_id": "uuid",
    "created_at": "random_date",
    "source": "random_hash"
  },
  "skip_tags_probability": 0.2,
  "clean_buckets": true,
  "num_workers": 5,
  "kafka": {
    "enabled": true,
    "brokers": ["kafka:9092"],
    "topic": "file-metadata"
  }
}
```

---

## 📝 Описание параметров конфига

| Поле                   | Описание                                                            |
|------------------------|---------------------------------------------------------------------|
| file_count             | Кол-во генерируемых файлов                                          |
| file_types             | Массив типов файлов (`text`, `image`)                               |
| buckets                | Массив бакетов MinIO для загрузки файлов                            |
| tags                   | Словарь тегов (ключ → способ генерации или дефолтное значение)      |
| skip_tags_probability  | Вероятность **не устанавливать тег** для объекта (0–1)              |
| clean_buckets          | Если `true`, очищает бакеты перед загрузкой                         |
| num_workers            | Кол-во воркеров для многопоточной генерации                         |
| kafka                  | Опциональный блок для публикации метаданных в Kafka                 |

---

## ⚙️ Переменные окружения

Генератор использует следующие переменные окружения, которые можно задать в `.env` или в `docker-compose.yml`:

| Переменная               | Описание                                                           | Дефолт |
|--------------------------|--------------------------------------------------------------------|--------|
| `MINIO_ROOT_USER`        | Логин для доступа к MinIO                                          | -      |
| `MINIO_ROOT_PASSWORD`    | Пароль для доступа к MinIO                                         | -      |
| `MINIO_PORT`             | Порт MinIO API                                                     | -      |
| `PROGRESS_INTERVAL_SEC`  | Интервал обновления прогресса в логах (в секундах)                 | 5      |

### Пример `.env` файла

```
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin
MINIO_PORT=9000
MINIO_CONSOLE_PORT=9001
PROGRESS_INTERVAL_SEC=5
```

---

## ✅ Kafka

Kafka полностью опциональна. Если в конфиге нет блока `kafka`, сервис просто не будет пытаться отправлять сообщения.

Пример отправляемых метаданных:

```json
{
  "bucket": "test-bucket",
  "object_name": "file_001_a1b2c3d4-5678-9abc-def0-123456789abc.txt",
  "tags": {
    "entity_id": "9b3b644e-...",
    "created_at": "2023-07-08T..."
  }
}
```

---

## 🧪 Итоговый лог

После завершения генерации сервис выведет:

```
[PROGRESS] Загружено файлов: 2200 / 5000 (44.0%), ETA: 38s, Kafka: 2198, Данных: 12.3 MB
[INFO] Генерация завершена за 1m2.431s.
[INFO] Всего загружено файлов: 5000
[INFO] Всего отправлено сообщений в Kafka: 5000
[INFO] Общий размер загруженных данных: 27.84 MB
```

Если Kafka выключена → количество сообщений будет `0`.

---

## 🧹 Очистка

Остановить и удалить все контейнеры:

```bash
docker compose down -v
```

Удалить образы:

```bash
docker image prune -f
```

---

## 💡 Возможности расширения

- Поддержка других типов файлов (PDF, ZIP)
- Генерация сложных структур тегов
- Гибкие стратегии распределения файлов по бакетам
- Интеграция с другими системами (например, RabbitMQ)

---

## 📜 Лицензия

This project is Open Source software released under the [Apache License 2.0](http://www.apache.org/licenses/LICENSE-2.0).

---
