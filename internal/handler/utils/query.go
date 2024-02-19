package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
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
