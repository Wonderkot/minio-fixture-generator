package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"minio-fixture-generator/internal/config"
	"minio-fixture-generator/internal/generator"
	"minio-fixture-generator/internal/kafka"
	"minio-fixture-generator/internal/miniohandler"

	"github.com/google/uuid"
)

type FileGenerator struct {
	cfg       *config.Config
	client    *miniohandler.Client
	publisher *kafka.Publisher
	logger    *log.Logger
}

type GenerationStats struct {
	TotalFiles int
	TotalKafka int
	TotalBytes int64
	Duration   time.Duration
}

func NewFileGenerator(cfg *config.Config, client *miniohandler.Client, publisher *kafka.Publisher, logger *log.Logger) *FileGenerator {
	return &FileGenerator{
		cfg:       cfg,
		client:    client,
		publisher: publisher,
		logger:    logger,
	}
}

func (s *FileGenerator) Run(ctx context.Context) (*GenerationStats, error) {
	intervalSec := 5 // дефолт

	if envVal, ok := os.LookupEnv("PROGRESS_INTERVAL_SEC"); ok {
		val, err := strconv.Atoi(envVal)
		if err == nil && val > 0 {
			intervalSec = val
		}
	}

	start := time.Now()

	// Очистка бакетов если требуется
	if s.cfg.CleanBuckets {
		for _, bucket := range s.cfg.Buckets {
			err := s.client.ClearBucket(ctx, bucket)
			if err != nil {
				return nil, fmt.Errorf("ошибка очистки бакета %s: %w", bucket, err)
			}
		}
	}

	// Создание бакетов (если вдруг кто-то удалил)
	for _, bucket := range s.cfg.Buckets {
		err := s.client.CreateBucketIfNotExists(ctx, bucket)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания бакета %s: %w", bucket, err)
		}
	}

	s.logger.Printf("[INFO] Старт генерации %d файлов с %d воркерами...", s.cfg.FileCount, s.cfg.WorkerCount())

	jobs := make(chan int, s.cfg.FileCount)
	var wg sync.WaitGroup

	var totalKafka int
	var totalFiles int
	var totalBytes int64
	var mu sync.Mutex

	// Прогресс-репортёр
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				mu.Lock()
				files := totalFiles
				kafkaMsgs := totalKafka
				bytes := totalBytes
				mu.Unlock()

				progress := float64(files) / float64(s.cfg.FileCount) * 100

				eta := "н/д"
				elapsed := time.Since(start).Seconds()
				speed := float64(files) / elapsed
				if speed > 0 && files > 0 {
					remaining := float64(s.cfg.FileCount - files)
					etaSecs := remaining / speed
					eta = time.Duration(etaSecs * float64(time.Second)).Truncate(time.Second).String()
				}

				s.logger.Printf(
					"[PROGRESS] Загружено файлов: %d / %d (%.1f%%), ETA: %s, Kafka: %d, Данных: %.2f MB",
					files, s.cfg.FileCount, progress, eta, kafkaMsgs, float64(bytes)/(1024*1024),
				)
			}
		}
	}()

	for w := 0; w < s.cfg.WorkerCount(); w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := range jobs {
				bucket := s.cfg.Buckets[i%len(s.cfg.Buckets)]
				fileType := s.cfg.FileTypes[i%len(s.cfg.FileTypes)]

				file, err := generator.GenerateFile(generator.FileType(fileType), i+1)
				if err != nil {
					s.logger.Printf("[Worker %d] Ошибка генерации файла: %v", workerID, err)
					continue
				}

				ext := s.detectExtension(fileType)
				fileName := s.GenerateFileName("file", i+1, ext)
				tags := generator.GenerateTags(s.cfg.Tags, s.cfg.SkipTagsProbability)

				err = s.client.UploadFile(ctx, bucket, fileName, file.Content, tags)
				if err != nil {
					s.logger.Printf("[Worker %d] Ошибка загрузки %s: %v", workerID, fileName, err)
					continue
				}

				mu.Lock()
				totalFiles++
				totalBytes += int64(len(file.Content))
				mu.Unlock()

				if s.cfg.Kafka != nil && s.cfg.Kafka.Enabled && s.publisher != nil {
					meta := kafka.FileMetadata{
						Bucket:     bucket,
						ObjectName: fileName,
						Tags:       tags,
					}
					err := s.publisher.Send(ctx, meta)
					if err != nil {
						s.logger.Printf("[Worker %d] Ошибка отправки в Kafka: %v", workerID, err)
					} else {
						mu.Lock()
						totalKafka++
						mu.Unlock()
					}
				}
			}
		}(w + 1)
	}

	for i := 0; i < s.cfg.FileCount; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()

	// Останавливаем прогресс-репортёр
	close(done)

	duration := time.Since(start)

	stats := &GenerationStats{
		TotalFiles: totalFiles,
		TotalKafka: totalKafka,
		TotalBytes: totalBytes,
		Duration:   duration,
	}
	return stats, nil
}

func (s *FileGenerator) GenerateFileName(baseName string, index int, extension string) string {
	return fmt.Sprintf("%s_%03d_%s.%s",
		baseName,
		index,
		uuid.New().String(),
		extension,
	)
}

func (s *FileGenerator) detectExtension(fileType string) string {
	fileType = strings.ToLower(fileType)
	switch fileType {
	case "image":
		return "jpg"
	case "text":
		return "txt"
	default:
		return "bin"
	}
}
