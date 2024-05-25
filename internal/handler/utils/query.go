package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/gin-gonic/gin"
	"io"
	"strconv"
)

var (
	ErrInvalidFileUpload = errors.New("invalid file upload")
	ErrConvFileToBytes   = errors.New("failed to copy image into bytes")
)

func GetIntFromParams(c gin.Params, param string) (int64, error) {
	p := c.ByName(param)
	if p == "" {
		return 0, fmt.Errorf("%s parameter must be provided", param)
	}

	i, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s parameter must be an integer", param)
	}

	return i, nil
}

func GetIntFromQuery(c *gin.Context, query string) (int, error) {
	q, ok := c.GetQuery(query)
	if !ok {
		return 0, fmt.Errorf("%s query parameter must be provided", query)
	}
	res, err := strconv.Atoi(q)
	if err != nil {
		return 0, fmt.Errorf("%s query parameter must be an integer", query)
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
