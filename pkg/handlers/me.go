package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/michaelchandrag/botfood-go/pkg/modules/me/dto"
	middlewareEntity "github.com/michaelchandrag/botfood-go/pkg/modules/middleware/entities"
	"github.com/michaelchandrag/botfood-go/pkg/protocols/error"
)

func (h *Handler) GetMeAction(c *gin.Context) {

	var customError error.ErrorCollection
	authBrand, existsAuthBrand := c.Get("auth_brand")
	if !existsAuthBrand {
		customError.AddHTTPError(401, errors.New("Invalid auth brand"))
		h.deliverError(c, customError)
		return
	}
	auth := authBrand.(middlewareEntity.Brand)
	payload := dto.MeAuthRequestPayload{
		AuthBrand: auth,
	}
	serviceAuth := h.meService.FormatAuthFromMiddleware(payload)
	if serviceAuth.Errors.HasErrors() {
		h.deliverError(c, serviceAuth.Errors)
		return
	}
	h.deliverJSON(c, serviceAuth.Auth)
	return
}

func (h *Handler) GetMeReviewsAction(c *gin.Context) {

	var customError error.ErrorCollection
	authBrand, existsAuthBrand := c.Get("auth_brand")
	if !existsAuthBrand {
		customError.AddHTTPError(401, errors.New("Invalid auth brand"))
		h.deliverError(c, customError)
		return
	}

	auth := authBrand.(middlewareEntity.Brand)
	payload := dto.MeAuthRequestPayload{
		AuthBrand: auth,
	}
	serviceAuth := h.meService.FormatAuthFromMiddleware(payload)
	if serviceAuth.Errors.HasErrors() {
		h.deliverError(c, serviceAuth.Errors)
		return
	}

	var reviewPayload dto.MeReviewsRequestPayload
	c.Bind(&reviewPayload)
	authBrandID := int(serviceAuth.Auth.Brand.ID)
	reviewPayload.BrandID = &authBrandID
	reviewPayload.BranchIDs = serviceAuth.Auth.BranchIDs
	if reviewPayload.QueryWithComment == "0" {
		reviewPayload.WithComment = true
	}

	if reviewPayload.QueryWithImages == "0" {
		reviewPayload.WithImages = true
	}

	if reviewPayload.QueryWithMerchantReply == "0" {
		reviewPayload.WithMerchantReply = true
	}

	serviceResult := h.meService.GetReviews(reviewPayload)
	if serviceResult.Errors.HasErrors() {
		h.deliverError(c, serviceAuth.Errors)
		return
	}

	h.deliverJSON(c, serviceResult.Data)
	return
}
