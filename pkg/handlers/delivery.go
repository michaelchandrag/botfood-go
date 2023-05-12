package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/michaelchandrag/botfood-go/pkg/protocols/error"
)

type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

type ErrorPayload struct {
	Code   int    `json:"code"`
	Detail string `json:"detail"`
}

type ErrorResponse struct {
	Errors  []ErrorPayload `json:"errors"`
	Success bool           `json:"success"`
}

func (h *Handler) deliverJSON(c *gin.Context, payload interface{}) {
	response := SuccessResponse{
		Data:    payload,
		Success: true,
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, response)
	return
}

func (h *Handler) deliverError(c *gin.Context, e error.ErrorCollection) {
	var collections ErrorResponse
	defaultHTTPCode := 400
	for key, val := range e.Errors {
		if key == 0 {
			defaultHTTPCode = val.HTTPCode
		}
		err := ErrorPayload{
			Code:   val.HTTPCode,
			Detail: val.Message,
		}
		collections.Errors = append(collections.Errors, err)
	}
	c.Header("Content-Type", "application/json")
	c.JSON(defaultHTTPCode, collections)
	return
}
