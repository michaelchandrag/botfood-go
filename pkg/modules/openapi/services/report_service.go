package services

import (
	"encoding/json"
	"fmt"
	"time"

	dto "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/dto"
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/entities"
	"github.com/michaelchandrag/botfood-go/utils"

	item_availability_report_repository "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/repositories/item_availability_report"
)

func (s *service) GetItemAvailabilityReports(payload dto.OpenApiReportItemAvailabilityReportsRequestPayload) (response dto.OpenApiReportItemAvailabilityReportListResponse) {
	itemReportRepository := item_availability_report_repository.NewRepository(s.db)

	now := time.Now()
	date := now.AddDate(0, 0, -1).Format("2006-01-02")
	if payload.Date != "" {
		payloadDate, err := time.Parse("2006-01-02", payload.Date)
		if err != nil {
			fmt.Println(err)
		} else {
			date = payloadDate.Format("2006-01-02")
		}
	}
	reportFilter := item_availability_report_repository.Filter{
		BrandID:         payload.BrandID,
		BranchChannelID: payload.BranchChannelID,
		Date:            date,
		Data:            payload.Data,
		Page:            payload.Page,
	}
	itemReports, err := itemReportRepository.FindPaginated(reportFilter)
	if err != nil {
		response.Errors.AddHTTPError(500, err)
		return response
	}

	for key, itemReport := range itemReports.Data {
		if itemReport.PayloadRemarks != nil {
			remarks := *itemReport.PayloadRemarks
			if utils.IsJSON(*itemReport.PayloadRemarks) {
				err = json.Unmarshal([]byte(remarks), &itemReports.Data[key].Remarks)
				if err != nil {
					fmt.Println(err)
					response.Errors.AddHTTPError(500, err)
					return response
				}
			}
		}
	}
	if itemReports.Data != nil {
		response.Data.Reviews = itemReports.Data
	} else {
		response.Data.Reviews = []entities.ItemAvailabilityReport{}
	}
	response.Data.CurrentPage = itemReports.CurrentPage
	response.Data.LimitData = itemReports.LimitData
	response.Data.TotalPage = itemReports.TotalPage
	response.Data.TotalData = itemReports.TotalData
	return response
}
