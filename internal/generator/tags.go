package generator

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GenerateTags(baseTags map[string]string, skipProbability float64) map[string]string {
	tags := make(map[string]string)

	for key, rule := range baseTags {
		if rand.Float64() < skipProbability {
			continue // пропустить тэг с вероятностью skipProbability
		}

		lowerRule := strings.ToLower(rule)
		var value string

		switch {
		case lowerRule == "uuid":
			value = uuid.New().String()
		case lowerRule == "random_hash":
			value = randomHash()
		case lowerRule == "random_date":
			value = randomDate().Format(time.RFC3339)
		case strings.Contains(lowerRule, "date"):
			value = time.Now().Format(time.RFC3339)
		default:
			value = rule
		}

		tags[key] = value
	}

	return tags
}

func randomHash() string {
	randomBytes := []byte(uuid.New().String())
	hash := sha1.Sum(randomBytes)
	return hex.EncodeToString(hash[:])
}

func randomDate() time.Time {
	start := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now()
	// Случайный момент времени между start и end
	sec := rand.Int63n(end.Unix() - start.Unix())
	return time.Unix(start.Unix()+sec, 0)
}
