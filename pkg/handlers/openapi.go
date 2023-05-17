package handlers

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/dto"
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
		for key, data := range payload.DataOutlet {
			payload.DataOutlet[key].IsOpen = utils.Itob(data.PayloadIsOpen)
		}
	} else if payload.Type == "activity-item" {
		for key, data := range payload.DataItem.ItemNew {
			payload.DataItem.ItemNew[key].InStock = utils.Itob(data.PayloadInStock)
		}
		for key, data := range payload.DataItem.ItemChange {
			payload.DataItem.ItemChange[key].InStock = utils.Itob(data.PayloadInStock)
		}
		for key, data := range payload.DataItem.ItemDelete {
			payload.DataItem.ItemDelete[key].InStock = utils.Itob(data.PayloadInStock)
		}
	}
	payload.RawMessage = string(rawMessage)
	action := h.consumerService.ConsumeActivityMessage(payload)
	if action.Errors.HasErrors() {
		h.deliverError(c, action.Errors)
		return
	}

	h.deliverJSON(c, action.MessageQueue)
	return
}
