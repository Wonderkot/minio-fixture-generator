package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"minio-fixture-generator/internal/config"
	"minio-fixture-generator/internal/generator"
	"minio-fixture-generator/internal/miniohandler"
)

func main() {
	fmt.Println("os.Args:", os.Args)

	configPath := flag.String("config", "config1.json", "Путь до JSON-конфига")
	flag.Parse()
	fmt.Println("Используемый путь к конфигу:", *configPath)

	logger := log.New(os.Stdout, "", log.LstdFlags)

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	logger.Printf("Загружен конфиг: %+v", *cfg)

	ctx := context.Background()
	client, err := miniohandler.New(ctx, logger)
	if err != nil {
		logger.Fatalf("Ошибка инициализации MinIO клиента: %v", err)
	}

	// создаем бакеты
	for _, bucket := range cfg.Buckets {
		err := client.CreateBucketIfNotExists(ctx, bucket)
		if err != nil {
			logger.Fatalf("Ошибка создания бакета %s: %v", bucket, err)
		}
	}

	logger.Printf("Начинаем генерацию %d файлов...", cfg.FileCount)

	for i := 0; i < cfg.FileCount; i++ {
		bucket := cfg.Buckets[i%len(cfg.Buckets)]
		fileType := cfg.FileTypes[i%len(cfg.FileTypes)]

		file, err := generator.GenerateFile(generator.FileType(fileType), i+1)
		if err != nil {
			logger.Printf("Ошибка генерации файла: %v", err)
			continue
		}

		tags := generator.GenerateTags(cfg.Tags, cfg.SkipTagsProbability)

		err = client.UploadFile(ctx, bucket, file.Name, file.Content, tags)
		if err != nil {
			logger.Printf("Ошибка загрузки файла %s: %v", file.Name, err)
		} else {
			logger.Printf("✅ Загружен файл: %s → %s, теги: %+v", file.Name, bucket, tags)
		}
	}

	logger.Println("Генерация и загрузка завершены.")
}
