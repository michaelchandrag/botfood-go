package services

import (
	"errors"

	"github.com/michaelchandrag/botfood-go/pkg/modules/health/dto"
)

type HealthService interface {
	GetHealth(payload dto.HealthRequestParams) (response dto.HealthResponse)
	GetError() (response dto.HealthResponse)
}

type service struct {
}

func RegisterHealthService() HealthService {
	return &service{}
}

func (service *service) GetHealth(payload dto.HealthRequestParams) (response dto.HealthResponse) {
	if payload.Message != "" {
		response.Message = payload.Message
		return response
	}
	response.Message = "ok"
	return response
}

func (service *service) GetError() (response dto.HealthResponse) {
	e := errors.New("First Error")
	response.Errors.AddHTTPError(400, e)
	e = errors.New("Second Error")
	response.Errors.AddHTTPError(400, e)
	return response
}
