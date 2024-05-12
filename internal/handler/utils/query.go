package utils

import (
	"bytes"
	"errors"
	"fmt"
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

func GetFileByName(c *gin.Context, name string) ([]byte, int64, error) {
	fileHeader, err := c.FormFile(name)
	if err != nil {
		return nil, 0, ErrInvalidFileUpload
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, 0, ErrInvalidFileUpload
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, 0, ErrConvFileToBytes
	}

	return buf.Bytes(), fileHeader.Size, nil
}
