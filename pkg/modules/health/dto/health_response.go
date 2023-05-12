package dto

import "github.com/michaelchandrag/botfood-go/pkg/protocols/error"

type HealthResponse struct {
	Message string                `json:"message"`
	Errors  error.ErrorCollection `json:"errors"`
}
