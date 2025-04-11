package common

import (
	"github.com/WlayRay/order-demo/common/tracing"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseResponse struct{}

type response struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"traceID"`
	Data    any    `json:"data"`
}

func (base BaseResponse) Response(c *gin.Context, err error, data any) {
	if err != nil {
		base.error(c, err)
	} else {
		base.success(c, data)
	}
}

func (base BaseResponse) success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, response{
		Success: true,
		Code:    0,
		Message: "",
		Data:    data,
		TraceID: tracing.TraceID(c.Request.Context()),
	})
}

func (base BaseResponse) error(c *gin.Context, err error) {
	c.JSON(http.StatusOK, response{
		Success: false,
		Code:    2,
		Message: err.Error(),
		Data:    nil,
		TraceID: tracing.TraceID(c.Request.Context()),
	})
}
