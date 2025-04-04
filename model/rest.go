package model

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpResponse struct {
	Meta HttpResponseMeta
	Data interface{}
}

type HttpResponseMeta struct {
	StatusCode int
	Errors     []string
}

func (h *HttpResponse) AddData(data interface{}) *HttpResponse {
	h.Data = data
	return h
}

func (h *HttpResponse) AddErrors(errors ...error) *HttpResponse {
	return h
}

func (h *HttpResponse) GinSuccess(c *gin.Context) {
	statusCode := http.StatusOK
	h.Meta.StatusCode = statusCode
	c.JSON(statusCode, h)
}

func (h *HttpResponse) GinErrorBadRequest(c *gin.Context) {
	statusCode := http.StatusBadRequest
	h.Meta.StatusCode = statusCode
	c.JSON(statusCode, h)
}

func (h *HttpResponse) GinErrorTimeout(c *gin.Context) {
	statusCode := http.StatusGatewayTimeout
	h.Meta.StatusCode = statusCode
	c.JSON(statusCode, h)
}

func (h *HttpResponse) GinErrorInternal(c *gin.Context) {
	statusCode := http.StatusInternalServerError
	h.Meta.StatusCode = statusCode
	c.JSON(statusCode, h)
}
