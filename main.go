package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"minio-fixture-generator/internal/config"
	"minio-fixture-generator/internal/kafka"
	"minio-fixture-generator/internal/miniohandler"
	"minio-fixture-generator/internal/service"
)

func main() {
	configPath := flag.String("config", "config.json", "Путь до JSON-конфига")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	logger.Printf("Используем конфиг: %+v", *cfg)

	ctx := context.Background()

	// подключение к MinIO
	client, err := miniohandler.New(ctx, logger)
	if err != nil {
		logger.Fatalf("Ошибка инициализации MinIO клиента: %v", err)
	}

	start := time.Now()

	var publisher *kafka.Publisher
	if cfg.Kafka != nil && cfg.Kafka.Enabled {
		publisher, err = kafka.NewPublisher(cfg.Kafka.Brokers, cfg.Kafka.Topic, logger)
		if err != nil {
			logger.Fatalf("Ошибка инициализации Kafka publisher: %v", err)
		}
		defer publisher.Close()
	}

	// сервис генерации
	genService := service.NewFileGenerator(cfg, client, publisher, logger)

	stats, err := genService.Run(ctx)
	if err != nil {
		logger.Fatalf("Ошибка генерации файлов: %v", err)
	}

	duration := time.Since(start)
	logger.Println("✅ Работа сервиса завершена успешно.")
	logger.Printf("Всего загружено файлов: %d", stats.TotalFiles)
	logger.Printf("Всего отправлено сообщений в Kafka: %d", stats.TotalKafka)
	logger.Printf("Общий размер загруженных данных: %s", humanReadableBytes(stats.TotalBytes))
	logger.Printf("Продолжительность генерации: %s", duration)
}

func humanReadableBytes(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
