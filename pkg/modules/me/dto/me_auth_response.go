package dto

import (
	"github.com/michaelchandrag/botfood-go/pkg/modules/me/entities"
	"github.com/michaelchandrag/botfood-go/pkg/protocols/error"
)

type MeAuthResponse struct {
	Auth   entities.Auth         `json:"auth"`
	Errors error.ErrorCollection `json:"errors"`
}
