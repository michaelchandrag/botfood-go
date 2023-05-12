package services

import (
	"errors"

	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/pkg/modules/report/dto"
)

type ReportService interface {
	ExportChannelReport(payload dto.ReportRequestPayload) (response dto.ReportResponse)
}

type service struct {
	db database.MainDB
}

func RegisterReportService(db database.MainDB) ReportService {
	return &service{
		db: db,
	}
}

func (s *service) ExportChannelReport(payload dto.ReportRequestPayload) (response dto.ReportResponse) {
	e := errors.New("First Error")
	response.Errors.AddHTTPError(400, e)
	return response
}
