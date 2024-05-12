package image

import (
	"errors"
	"github.com/google/uuid"
	"github.com/h2non/bimg"
	"strings"
)

var (
	ErrImageQuality = errors.New("image quality must be between 1 and 100")
	ErrImageFormat  = errors.New("unsupported image format")
	ErrImageIsEmpty = errors.New("image buffer is empty")
)

func CompressImage(buffer []byte, quality int) ([]byte, string, error) {
	if quality <= 0 || quality > 100 {
		return nil, "", ErrImageQuality
	}

	filename := strings.Replace(uuid.New().String(), "-", "", -1) + ".webp"

	converted, err := bimg.NewImage(buffer).Convert(bimg.WEBP)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unsupported image format"):
			return nil, "", ErrImageFormat
		case strings.Contains(err.Error(), "Image buffer is empty"):
			return nil, "", ErrImageIsEmpty
		}
		return nil, "", err
	}

	processed, err := bimg.NewImage(converted).Process(bimg.Options{Quality: quality})
	if err != nil {
		return nil, "", err
	}

	return processed, filename, nil
}
