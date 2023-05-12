package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/michaelchandrag/botfood-go/pkg/modules/health/dto"
	"github.com/michaelchandrag/botfood-go/pkg/protocols/error"
)

func (h *Handler) GetHealthAction(c *gin.Context) {

	var customError error.ErrorCollection

	var payload dto.HealthRequestParams
	if err := c.BindQuery(&payload); err != nil {
		customError.AddHTTPError(500, errors.New(err.Error()))
		h.deliverError(c, customError)
		return
	}

	healthAction := h.healthService.GetHealth(payload)
	h.deliverJSON(c, healthAction.Message)
	return
}

func (h *Handler) GetErrorAction(c *gin.Context) {

	errorAction := h.healthService.GetError()
	if errorAction.Errors.HasErrors() {
		h.deliverError(c, errorAction.Errors)
		return
	}
	h.deliverJSON(c, errorAction)
	return
}
