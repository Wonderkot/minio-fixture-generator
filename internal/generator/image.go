package generator

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"
)

func generateImageFile(index int) (*FileData, error) {
	filename := fmt.Sprintf("file_%03d.jpg", index)

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// случайный цвет
	c := color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}

	// заливаем всё изображение
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, c)
		}
	}

	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации jpg: %w", err)
	}

	return &FileData{
		Name:    filename,
		Content: buf.Bytes(),
		Type:    ImageFile,
	}, nil
}
