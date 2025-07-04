# ---------- Этап сборки ----------
FROM golang:1.23.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/minio-fixture-generator ./

# ---------- Финальный образ ----------
FROM alpine:3.20

WORKDIR /app

# Установим CA-сертификаты
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/minio-fixture-generator /app/minio-fixture-generator

ENTRYPOINT ["/app/minio-fixture-generator"]
