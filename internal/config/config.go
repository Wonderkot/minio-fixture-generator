package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	FileCount           int               `json:"file_count"`
	FileTypes           []string          `json:"file_types"`
	Buckets             []string          `json:"buckets"`
	Tags                map[string]string `json:"tags"`
	SkipTagsProbability float64           `json:"skip_tags_probability"`
	CleanBuckets        bool              `json:"clean_buckets"`
	Kafka               *KafkaConfig      `json:"kafka,omitempty"`
	NumWorkers          int               `json:"num_workers`
}

type KafkaConfig struct {
	Enabled bool     `json:"enabled"`
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
}

func Load(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать конфиг: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("не удалось распарсить конфиг: %w", err)
	}

	cfg.applyDefaults()
	return &cfg, nil
}

func (c *Config) applyDefaults() {
	if c.FileCount <= 0 {
		c.FileCount = 10
	}
	if len(c.FileTypes) == 0 {
		c.FileTypes = []string{"text"}
	}
	if len(c.Buckets) == 0 {
		c.Buckets = []string{"default-bucket"}
	}
	if c.Tags == nil {
		c.Tags = make(map[string]string)
	}
	if c.SkipTagsProbability < 0 || c.SkipTagsProbability > 1 {
		c.SkipTagsProbability = 0
	}
}

func (c *Config) WorkerCount() int {
	if c.NumWorkers <= 0 {
		return 5 // дефолт
	}
	return c.NumWorkers
}
