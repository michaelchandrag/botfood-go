package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/dto"
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/entities"
	brand_repository "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/repositories/brand"
	message_queue_repository "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/repositories/message_queue"
	webhook_log_repository "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/repositories/webhook_log"
	"github.com/michaelchandrag/botfood-go/utils"
)

func (s *service) ConsumeActivityMessage(payload dto.ConsumerRequestPayload) (response dto.ConsumerResponsePayload) {

	messageQueueRepository := message_queue_repository.NewRepository(s.db)
	queueFilter := message_queue_repository.Filter{
		BrandID:   &payload.BrandID,
		Type:      payload.Type,
		MessageID: payload.MessageID,
	}
	existsQueue, err := messageQueueRepository.FindOne(queueFilter)
	if err != nil {
		response.Errors.AddHTTPError(500, err)
		return response
	}
	if existsQueue.ID == 0 {
		newQueue := entities.MessageQueue{
			MessageID: payload.MessageID,
			Type:      payload.Type,
			BrandID:   payload.BrandID,
			Body:      payload.RawMessage,
		}
		actionQueue, err := messageQueueRepository.Create(newQueue)
		if err != nil {
			response.Errors.AddHTTPError(500, err)
			return response
		}
		response.MessageQueue = actionQueue

		brandRepository := brand_repository.NewRepository(s.db)
		brandFilter := brand_repository.Filter{
			ID: &payload.BrandID,
		}
		existsBrand, err := brandRepository.FindOne(brandFilter)
		if existsBrand.ID == 0 {
			response.Errors.AddHTTPError(400, errors.New("Brand is required"))
			return response
		}
		// send webhooks
		webhookBody := dto.WebhookRequestPayload{
			Type: payload.Type,
		}

		if payload.Type == "activity-outlet" {
			var webhookOutlets []dto.WebhookOutletRequestPayload
			for _, val := range payload.DataOutlets {
				webhookOutlets = append(webhookOutlets, val.ToWebhookOutletRequestPayload())
			}
			webhookBody.DataOutlets = &webhookOutlets
		} else if payload.Type == "activity-item" {
			var webhookDataItem dto.WebhookDataItem
			var webhookNewItems []dto.WebhookItemRequestPayload
			var webhookChangeItems []dto.WebhookItemRequestPayload
			var webhookDeletedItems []dto.WebhookItemRequestPayload
			for _, val := range payload.DataItems.ItemNew {
				webhookNewItems = append(webhookNewItems, val.ToWebhookItemRequestPayload())
			}
			for _, val := range payload.DataItems.ItemChange {
				webhookChangeItems = append(webhookNewItems, val.ToWebhookItemRequestPayload())
			}
			for _, val := range payload.DataItems.ItemDelete {
				webhookDeletedItems = append(webhookNewItems, val.ToWebhookItemRequestPayload())
			}
			webhookDataItem.ItemNew = &webhookNewItems
			webhookDataItem.ItemChange = &webhookChangeItems
			webhookDataItem.ItemDelete = &webhookDeletedItems
			webhookBody.DataItems = &webhookDataItem
		}

		go s.sendWebhook(existsBrand.WebhookURL, webhookBody, existsBrand)

	} else {
		response.Errors.AddHTTPError(400, errors.New("Message Queue already exists"))
		return response
	}

	return response
}

func (s *service) sendWebhook(url string, payload dto.WebhookRequestPayload, brand entities.Brand) {
	requestBody, _ := json.Marshal(payload)
	logRepository := webhook_log_repository.NewRepository(s.db)
	logObj := entities.WebhookLog{
		BrandID:     int(brand.ID),
		RequestURL:  url,
		RequestBody: string(requestBody),
	}

	trials := [3]int{0, 1, 2}
	for _, secs := range trials {
		time.Sleep(time.Duration(secs) * time.Second)
		now := time.Now()
		timestamp := now.Unix()
		recipe := fmt.Sprintf("%s:%d", *brand.ApiKey, timestamp)
		partnerToken := utils.GenerateHMAC256(recipe, *brand.SecretKey)
		var theResult interface{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("X-Timestamp", fmt.Sprintf("%d", timestamp)).
			SetHeader("X-Partner-Token", partnerToken).
			SetHeader("Content-Type", "application/json").
			SetBody(payload).
			SetResult(&theResult). // or SetResult(AuthSuccess{}).
			Post(url)
		fmt.Println(resp.StatusCode())
		fmt.Println(resp)
		if err != nil {
			fmt.Println(err)
		}
		logObj.ResponseBody = string(resp.Body())
		newLog, err := logRepository.Create(logObj)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(newLog.ID)
		if resp.StatusCode() == 200 {
			break
		}
	}

}
