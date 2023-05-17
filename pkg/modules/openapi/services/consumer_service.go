package services

import (
	"errors"

	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/dto"
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/entities"
	message_queue_repository "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/repositories/message_queue"
)

type ConsumerService interface {
	ConsumeActivityMessage(payload dto.ConsumerRequestPayload) (response dto.ConsumerResponsePayload)
}

type service struct {
	db database.MainDB
}

func RegisterConsumerService(db database.MainDB) ConsumerService {
	return &service{
		db: db,
	}
}

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
	} else {
		response.Errors.AddHTTPError(400, errors.New("Message Queue already exists"))
		return response
	}

	return response
}
