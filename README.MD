# MinIO Fixture Generator

Генератор тестовых файлов (текста и изображений) с тэгами для MinIO, реализован на Go.  
Запускается в Docker и подключается к локальному MinIO-серверу.

---

## 🧰 Возможности

- Генерация текстовых и JPEG-файлов
- Установка кастомных тегов (включая UUID, хэши, текущие и случайные даты)
- Гибкая конфигурация через `config.json`
- Поддержка Docker Compose
- Переменные окружения в `.env`

---

## 📁 Структура проекта
```

├── cmd/generator # точка входа main.go
├── internal/ # внутренние пакеты: генератор, клиент, конфиг
├── config/config.json # конфигурация генерации
├── Dockerfile
├── docker-compose.yml
├── .env
├── Makefile
└── README.md
```


---

## ⚙️ Использование

### 1. Установи Docker и Docker Compose

Поддерживается Docker Desktop или обычный Docker Engine + Compose plugin.

### 2. Проверь `.env`

```env
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin
MINIO_PORT=9000
MINIO_CONSOLE_PORT=9001
```

### 3. Проверь config/config.json
Пример: 
```json
{
  "file_count": 10,
  "file_types": ["text", "image"],
  "buckets": ["test-bucket", "images"],
  "tags": {
    "entity_id": "uuid",
    "created_at": "random_date"
  },
  "skip_tags_probability": 0.2
}

```

### 4. Запусти
```bash
make build
make up
```

### 5. Останови
```bash
make down
```

🧪 Проверка с логами
```bash
make logs
```


🧹 Полная очистка
```bash
make clean
```

---

## 📜 Лицензия

This project is Open Source software released under the [Apache License 2.0](http://www.apache.org/licenses/LICENSE-2.0).

---
