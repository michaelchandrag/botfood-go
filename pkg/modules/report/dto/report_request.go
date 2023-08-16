package dto

import (
	entity "github.com/michaelchandrag/botfood-go/pkg/modules/report/entities"
)

type ReportRequestPayload struct {
	Brand     entity.Brand
	Date      string `form:"date"`
	BranchIDs []int
}
