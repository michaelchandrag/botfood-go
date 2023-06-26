package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/dto"
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/entities"
	"github.com/michaelchandrag/botfood-go/pkg/protocols/error"
	"github.com/michaelchandrag/botfood-go/utils"
)

func (h *Handler) PostConsumeMessageQueueAction(c *gin.Context) {

	rawMessage, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(rawMessage))

	var customError error.ErrorCollection
	var payload dto.ConsumerRequestPayload
	if err := c.BindJSON(&payload); err != nil {
		customError.AddHTTPError(400, err)
		h.deliverError(c, customError)
		return
	}

	if payload.Type == "activity-outlet" {
		for key, data := range payload.DataOutlets {
			payload.DataOutlets[key].IsOpen = utils.Itob(data.PayloadIsOpen)
		}
	} else if payload.Type == "activity-item" {
		for key, data := range payload.DataItems.ItemNew {
			payload.DataItems.ItemNew[key].InStock = utils.Itob(data.PayloadInStock)
		}
		for key, data := range payload.DataItems.ItemChange {
			payload.DataItems.ItemChange[key].InStock = utils.Itob(data.PayloadInStock)
		}
		for key, data := range payload.DataItems.ItemDelete {
			payload.DataItems.ItemDelete[key].InStock = utils.Itob(data.PayloadInStock)
		}
	}
	payload.RawMessage = string(rawMessage)
	action := h.openApiService.ConsumeActivityMessage(payload)
	if action.Errors.HasErrors() {
		h.deliverError(c, action.Errors)
		return
	}

	h.deliverJSON(c, action.MessageQueue)
	return
}

func (h *Handler) GetOpenApiBranchChannelListAction(c *gin.Context) {

	var customError error.ErrorCollection
	authBrand, existsAuthBrand := c.Get("open_api_brand")
	if !existsAuthBrand {
		customError.AddHTTPError(401, errors.New("Unauthorized. Missing partner"))
		h.deliverError(c, customError)
		return
	}

	brand := authBrand.(entities.Brand)

	var branchChannelPayload dto.OpenApiBranchChannelRequestPayload
	c.Bind(&branchChannelPayload)
	authBrandID := int(brand.ID)
	branchChannelPayload.BrandID = &authBrandID

	if len(branchChannelPayload.PayloadIsOpen) > 0 {
		isOpen, _ := strconv.Atoi(branchChannelPayload.PayloadIsOpen)
		branchChannelPayload.IsOpen = &isOpen
	}

	page, err := strconv.Atoi(branchChannelPayload.PayloadPage)
	if err == nil {
		branchChannelPayload.Page = &page
	}

	data, err := strconv.Atoi(branchChannelPayload.PayloadData)
	if err == nil {
		branchChannelPayload.Data = &data
	}

	serviceResult := h.openApiService.GetBranchChannels(branchChannelPayload)
	if serviceResult.Errors.HasErrors() {
		h.deliverError(c, serviceResult.Errors)
		return
	}

	h.deliverJSON(c, serviceResult.Data)
	return
}

func (h *Handler) GetOpenApiBranchChannelDetailAction(c *gin.Context) {

	var customError error.ErrorCollection
	authBrand, existsAuthBrand := c.Get("open_api_brand")
	if !existsAuthBrand {
		customError.AddHTTPError(401, errors.New("Unauthorized. Missing partner"))
		h.deliverError(c, customError)
		return
	}

	brand := authBrand.(entities.Brand)

	var branchChannelPayload dto.OpenApiBranchChannelRequestPayload
	authBrandID := int(brand.ID)
	branchChannelID, _ := strconv.Atoi(c.Param("branch_channel_id"))
	branchChannelPayload.BrandID = &authBrandID
	branchChannelPayload.ID = &branchChannelID

	serviceResult := h.openApiService.GetBranchChannelDetail(branchChannelPayload)
	if serviceResult.Errors.HasErrors() {
		h.deliverError(c, serviceResult.Errors)
		return
	}

	h.deliverJSON(c, serviceResult.Data)
	return
}

func (h *Handler) GetOpenApiItemListAction(c *gin.Context) {

	var customError error.ErrorCollection
	authBrand, existsAuthBrand := c.Get("open_api_brand")
	if !existsAuthBrand {
		customError.AddHTTPError(401, errors.New("Unauthorized. Missing partner"))
		h.deliverError(c, customError)
		return
	}

	brand := authBrand.(entities.Brand)

	var itemPayload dto.OpenApiItemRequestPayload
	c.Bind(&itemPayload)
	authBrandID := int(brand.ID)
	itemPayload.BrandID = &authBrandID

	if len(itemPayload.PayloadBranchChannelID) > 0 {
		branchChannelID, _ := strconv.Atoi(itemPayload.PayloadBranchChannelID)
		itemPayload.BranchChannelID = &branchChannelID
	}

	if len(itemPayload.PayloadInStock) > 0 {
		inStock, _ := strconv.Atoi(itemPayload.PayloadInStock)
		itemPayload.InStock = &inStock
	}

	page, err := strconv.Atoi(itemPayload.PayloadPage)
	if err == nil {
		itemPayload.Page = &page
	}

	data, err := strconv.Atoi(itemPayload.PayloadData)
	if err == nil {
		itemPayload.Data = &data
	}

	serviceResult := h.openApiService.GetItems(itemPayload)
	if serviceResult.Errors.HasErrors() {
		h.deliverError(c, serviceResult.Errors)
		return
	}

	h.deliverJSON(c, serviceResult.Data)
	return
}

func (h *Handler) GetOpenApiReviewListAction(c *gin.Context) {

	var customError error.ErrorCollection
	authBrand, existsAuthBrand := c.Get("open_api_brand")
	if !existsAuthBrand {
		customError.AddHTTPError(401, errors.New("Unauthorized. Missing partner"))
		h.deliverError(c, customError)
		return
	}

	brand := authBrand.(entities.Brand)

	var reviewPayload dto.OpenApiReviewsRequestPayload
	c.Bind(&reviewPayload)
	authBrandID := int(brand.ID)
	reviewPayload.BrandID = &authBrandID
	if reviewPayload.QueryWithComment == "0" {
		reviewPayload.WithComment = true
	}

	if reviewPayload.QueryWithImages == "0" {
		reviewPayload.WithImages = true
	}

	if reviewPayload.QueryWithMerchantReply == "0" {
		reviewPayload.WithMerchantReply = true
	}

	if reviewPayload.PayloadBranchChannelID != "" {
		bcId, _ := strconv.Atoi(reviewPayload.PayloadBranchChannelID)
		reviewPayload.BranchChannelID = &bcId
	}

	if reviewPayload.PayloadRating != "" {
		rating, _ := strconv.Atoi(reviewPayload.PayloadRating)
		reviewPayload.Rating = &rating
	}

	serviceResult := h.openApiService.GetReviews(reviewPayload)
	if serviceResult.Errors.HasErrors() {
		h.deliverError(c, serviceResult.Errors)
		return
	}

	h.deliverJSON(c, serviceResult.Data)
	return
}

func (h *Handler) GetOpenApiItemAvailabilityReportListAction(c *gin.Context) {

	var customError error.ErrorCollection
	authBrand, existsAuthBrand := c.Get("open_api_brand")
	if !existsAuthBrand {
		customError.AddHTTPError(401, errors.New("Unauthorized. Missing partner"))
		h.deliverError(c, customError)
		return
	}

	brand := authBrand.(entities.Brand)

	var reportPayload dto.OpenApiReportItemAvailabilityReportsRequestPayload
	c.Bind(&reportPayload)
	authBrandID := int(brand.ID)
	reportPayload.BrandID = &authBrandID

	if reportPayload.PayloadBranchChannelID != "" {
		bcId, _ := strconv.Atoi(reportPayload.PayloadBranchChannelID)
		reportPayload.BranchChannelID = &bcId
	}

	if reportPayload.PayloadDate != "" {
		date, err := time.Parse("2006-01-02", reportPayload.PayloadDate)
		if err != nil {
			fmt.Println(err)
		} else {
			reportPayload.Date = date.Format("2006-01-02")
		}
	}

	serviceResult := h.openApiService.GetItemAvailabilityReports(reportPayload)
	if serviceResult.Errors.HasErrors() {
		h.deliverError(c, serviceResult.Errors)
		return
	}

	h.deliverJSON(c, serviceResult.Data)
	return
}

func (h *Handler) GetOpenApiBranchChannelAvailabilityReportListAction(c *gin.Context) {

	var customError error.ErrorCollection
	authBrand, existsAuthBrand := c.Get("open_api_brand")
	if !existsAuthBrand {
		customError.AddHTTPError(401, errors.New("Unauthorized. Missing partner"))
		h.deliverError(c, customError)
		return
	}

	brand := authBrand.(entities.Brand)

	var reportPayload dto.OpenApiReportBranchChannelAvailabilityReportsRequestPayload
	c.Bind(&reportPayload)
	authBrandID := int(brand.ID)
	reportPayload.BrandID = &authBrandID

	if reportPayload.PayloadBranchChannelID != "" {
		bcId, _ := strconv.Atoi(reportPayload.PayloadBranchChannelID)
		reportPayload.BranchChannelID = &bcId
	}

	if reportPayload.PayloadDate != "" {
		date, err := time.Parse("2006-01-02", reportPayload.PayloadDate)
		if err != nil {
			fmt.Println(err)
		} else {
			reportPayload.Date = date.Format("2006-01-02")
		}
	}

	serviceResult := h.openApiService.GetBranchChannelAvailabilityReports(reportPayload)
	if serviceResult.Errors.HasErrors() {
		h.deliverError(c, serviceResult.Errors)
		return
	}

	h.deliverJSON(c, serviceResult.Data)
	return
}

func (h *Handler) GetOpenApiVariantListAction(c *gin.Context) {

	var customError error.ErrorCollection
	authBrand, existsAuthBrand := c.Get("open_api_brand")
	if !existsAuthBrand {
		customError.AddHTTPError(401, errors.New("Unauthorized. Missing partner"))
		h.deliverError(c, customError)
		return
	}

	brand := authBrand.(entities.Brand)

	var variantPayload dto.OpenApiVariantRequestPayload
	c.Bind(&variantPayload)
	authBrandID := int(brand.ID)
	variantPayload.BrandID = &authBrandID

	if len(variantPayload.PayloadBranchChannelID) > 0 {
		branchChannelID, _ := strconv.Atoi(variantPayload.PayloadBranchChannelID)
		variantPayload.BranchChannelID = &branchChannelID
	}

	if len(variantPayload.PayloadVariantCategoryID) > 0 {
		variantCategoryID, _ := strconv.Atoi(variantPayload.PayloadVariantCategoryID)
		variantPayload.VariantCategoryID = &variantCategoryID
	}

	if len(variantPayload.PayloadInStock) > 0 {
		inStock, _ := strconv.Atoi(variantPayload.PayloadInStock)
		variantPayload.InStock = &inStock
	}

	page, err := strconv.Atoi(variantPayload.PayloadPage)
	if err == nil {
		variantPayload.Page = &page
	}

	data, err := strconv.Atoi(variantPayload.PayloadData)
	if err == nil {
		variantPayload.Data = &data
	}

	serviceResult := h.openApiService.GetVariants(variantPayload)
	if serviceResult.Errors.HasErrors() {
		h.deliverError(c, serviceResult.Errors)
		return
	}

	h.deliverJSON(c, serviceResult.Data)
	return
}
