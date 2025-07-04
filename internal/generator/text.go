package generator

import (
	"fmt"

	"github.com/google/uuid"
)

func generateTextFile(index int) (*FileData, error) {
	id := uuid.New().String()
	filename := fmt.Sprintf("file_%03d.txt", index)

	return &FileData{
		Name:    filename,
		Content: []byte(id),
		Type:    TextFile,
	}, nil
}
