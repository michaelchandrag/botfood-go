package dto

import (
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/entities"
	"github.com/michaelchandrag/botfood-go/pkg/protocols/error"
)

type ConsumerResponsePayload struct {
	MessageQueue entities.MessageQueue
	Errors       error.ErrorCollection `json:"errors"`
}
