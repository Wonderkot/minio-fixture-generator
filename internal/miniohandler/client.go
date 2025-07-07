package miniohandler

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cenkalti/backoff/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	minioTags "github.com/minio/minio-go/v7/pkg/tags"
)

type Client struct {
	Minio  *minio.Client
	Logger *log.Logger
}

func New(ctx context.Context, logger *log.Logger) (*Client, error) {
	endpoint := os.Getenv("MINIO_HOST")
	accessKey := os.Getenv("MINIO_USER")
	secretKey := os.Getenv("MINIO_PASSWORD")

	if endpoint == "" || accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("переменные окружения MINIO_HOST, MINIO_USER и MINIO_PASSWORD обязательны")
	}

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: strings.HasPrefix(endpoint, "https://"),
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к MinIO: %w", err)
	}

	return &Client{
		Minio:  minioClient,
		Logger: logger,
	}, nil
}

func (c *Client) CreateBucketIfNotExists(ctx context.Context, bucket string) error {
	exists, err := c.Minio.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("ошибка проверки бакета %s: %w", bucket, err)
	}
	if !exists {
		err = c.Minio.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("ошибка создания бакета %s: %w", bucket, err)
		}
		c.Logger.Printf("Создан бакет: %s", bucket)
	}
	return nil
}

func (c *Client) UploadFile(ctx context.Context, bucket, objectName string, content []byte, tags map[string]string) error {
	opts := minio.PutObjectOptions{
		ContentType: guessContentType(objectName),
	}

	reader := bytes.NewReader(content)
	_, err := c.Minio.PutObject(ctx, bucket, objectName, reader, int64(len(content)), opts)
	if err != nil {
		return fmt.Errorf("ошибка загрузки %s в %s: %w", objectName, bucket, err)
	}

	if len(tags) > 0 {
		err := c.SetTagsWithRetry(ctx, bucket, objectName, tags)
		if err != nil {
			return fmt.Errorf("ошибка установки тегов на %s: %w", objectName, err)
		}
	}

	return nil
}

func (c *Client) SetTagsWithRetry(ctx context.Context, bucket, objectName string, tagMap map[string]string) error {
	operation := func() error {
		tagObj, err := minioTags.NewTags(tagMap, true)
		if err != nil {
			return err
		}
		return c.Minio.PutObjectTagging(ctx, bucket, objectName, tagObj, minio.PutObjectTaggingOptions{})
	}

	return backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3))
}

func guessContentType(name string) string {
	name = strings.ToLower(name)
	if strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg") {
		return "image/jpeg"
	}
	return "text/plain"
}

func (c *Client) ClearBucket(ctx context.Context, bucket string) error {
	objectCh := make(chan minio.ObjectInfo, 100)

	// запуск горутины для наполнения канала
	go func() {
		defer close(objectCh)
		for obj := range c.Minio.ListObjects(ctx, bucket, minio.ListObjectsOptions{
			Recursive: true,
		}) {
			if obj.Err != nil {
				c.Logger.Printf("Ошибка чтения объекта: %v", obj.Err)
				continue
			}
			objectCh <- obj
		}
	}()

	// удаление объектов пакетами
	for err := range c.Minio.RemoveObjects(ctx, bucket, objectCh, minio.RemoveObjectsOptions{}) {
		if err.Err != nil {
			return fmt.Errorf("ошибка удаления объекта %s: %v", err.ObjectName, err.Err)
		}
	}

	c.Logger.Printf("Бакет %s очищен.", bucket)
	return nil
}
