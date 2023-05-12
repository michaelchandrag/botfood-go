package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	middlewareEntity "github.com/michaelchandrag/botfood-go/pkg/modules/middleware/entities"
	"github.com/michaelchandrag/botfood-go/pkg/modules/report/dto"
	"github.com/michaelchandrag/botfood-go/pkg/modules/report/entities"
	"github.com/michaelchandrag/botfood-go/pkg/protocols/error"
)

func (h *Handler) GetChannelReportAction(c *gin.Context) {

	var customError error.ErrorCollection
	authBrand, existsAuthBrand := c.Get("auth_brand")
	if !existsAuthBrand {
		customError.AddHTTPError(401, errors.New("Invalid auth brand"))
		h.deliverError(c, customError)
		return
	}
	convertBrand := authBrand.(middlewareEntity.Brand)
	brand := entities.Brand{
		ID:   convertBrand.ID,
		Name: convertBrand.Name,
		Slug: convertBrand.Slug,
	}

	payload := dto.ReportRequestPayload{
		Brand: brand,
	}

	reportAction := h.reportService.ExportChannelReport(payload)
	if reportAction.Errors.HasErrors() {
		h.deliverError(c, reportAction.Errors)
		return
	}
	h.deliverJSON(c, reportAction.Message)
	return
}
