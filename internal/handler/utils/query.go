package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidFileUpload = errors.New("invalid file upload")
	ErrConvFileToBytes   = errors.New("failed to copy image into bytes")
	ErrMaxFilesCount     = errors.New("max files count exceeded")

	ErrMustbeProvided = errors.New("parameter must be provided")
	ErrMustBeInteger  = errors.New("parameter must be an integer")
)

func GetIntFromParams(c gin.Params, param string) (int64, error) {
	p := c.ByName(param)
	if p == "" {
		return 0, fmt.Errorf("%s %w", param, ErrMustbeProvided)
	}

	i, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s %w", param, ErrMustBeInteger)
	}

	return i, nil
}

func GetIntFromQuery(c *gin.Context, query string) (int, error) {
	q, ok := c.GetQuery(query)
	if !ok {
		return 0, fmt.Errorf("%s query %w", query, ErrMustbeProvided)
	}
	res, err := strconv.Atoi(q)
	if err != nil {
		return 0, fmt.Errorf("%s query %w", query, ErrMustBeInteger)
	}

	return res, nil
}

func GetFileByName(c *gin.Context, name string) (domain.File, error) {
	fileHeader, err := c.FormFile(name)
	if err != nil {
		return domain.File{}, ErrInvalidFileUpload
	}

	file, err := fileHeader.Open()
	if err != nil {
		return domain.File{}, ErrInvalidFileUpload
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return domain.File{}, ErrConvFileToBytes
	}

	return domain.File{
		Name:  fileHeader.Filename,
		Bytes: buf.Bytes(),
		Size:  fileHeader.Size,
		Type:  fileHeader.Header.Get("Content-Type"),
	}, nil
}

func GetFilesByName(c *gin.Context, name string, maxFilesCount int) ([]domain.File, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	files := form.File[name]
	if len(files) == 0 {
		return nil, fmt.Errorf("no files found with name %s", name)
	}
	if len(files) > maxFilesCount {
		return nil, fmt.Errorf("%w: max files count is %d", ErrMaxFilesCount, maxFilesCount)
	}

	var res []domain.File
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, ErrInvalidFileUpload
		}
		defer file.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			return nil, ErrConvFileToBytes
		}

		res = append(res, domain.File{
			Name:  fileHeader.Filename,
			Bytes: buf.Bytes(),
			Size:  fileHeader.Size,
			Type:  fileHeader.Header.Get("Content-Type"),
		})
	}

	return res, nil
}

func GetBoolFromQuery(c *gin.Context, query string) (bool, error) {
	q, ok := c.GetQuery(query)
	if !ok {
		return false, fmt.Errorf("%s query parameter must be provided", query)
	}

	res, err := strconv.ParseBool(q)
	if err != nil {
		return false, fmt.Errorf("%s query parameter must be a boolean", query)
	}

	return res, nil
}
