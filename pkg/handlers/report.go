package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	meDTO "github.com/michaelchandrag/botfood-go/pkg/modules/me/dto"
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

	auth := authBrand.(middlewareEntity.Brand)
	payload := meDTO.MeAuthRequestPayload{
		AuthBrand: auth,
	}
	serviceAuth := h.meService.FormatAuthFromMiddleware(payload)
	if serviceAuth.Errors.HasErrors() {
		h.deliverError(c, serviceAuth.Errors)
		return
	}

	brand := entities.Brand{
		ID:   serviceAuth.Auth.Brand.ID,
		Name: serviceAuth.Auth.Brand.Name,
		Slug: serviceAuth.Auth.Brand.Slug,
	}

	reportPayload := dto.ReportRequestPayload{
		Brand:     brand,
		BranchIDs: serviceAuth.Auth.BranchIDs,
	}

	reportAction := h.reportService.ExportChannelReport(reportPayload)
	if reportAction.Errors.HasErrors() {
		h.deliverError(c, reportAction.Errors)
		return
	}
	h.deliverExcel(c, reportAction.File.Excel)
	return
}
