package dto

import "github.com/michaelchandrag/botfood-go/pkg/protocols/error"

type ReportResponse struct {
	Message string                `json:"message"`
	Errors  error.ErrorCollection `json:"errors"`
}
