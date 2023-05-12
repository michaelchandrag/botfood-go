package dto

import (
	entity "github.com/michaelchandrag/botfood-go/pkg/modules/report/entities"
)

type ReportRequestPayload struct {
	Brand entity.Brand
}
