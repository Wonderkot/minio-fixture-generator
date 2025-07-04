package generator

import (
	"fmt"
	"math/rand"
	"time"
)

type FileType string

const (
	TextFile  FileType = "text"
	ImageFile FileType = "image"
)

type FileData struct {
	Name    string
	Content []byte
	Type    FileType
}

func GenerateFile(fileType FileType, index int) (*FileData, error) {
	switch fileType {
	case TextFile:
		return generateTextFile(index)
	case ImageFile:
		return generateImageFile(index)
	default:
		return nil, fmt.Errorf("неподдерживаемый тип файла: %s", fileType)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
