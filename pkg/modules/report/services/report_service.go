package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/pkg/modules/report/dto"
	"github.com/michaelchandrag/botfood-go/pkg/modules/report/entities"
	branch_channel_repository "github.com/michaelchandrag/botfood-go/pkg/modules/report/repositories/branch_channel"
	branch_channel_availability_report_repository "github.com/michaelchandrag/botfood-go/pkg/modules/report/repositories/branch_channel_availability_report"
	branch_channel_promotion_repository "github.com/michaelchandrag/botfood-go/pkg/modules/report/repositories/branch_channel_promotion"
	item_repository "github.com/michaelchandrag/botfood-go/pkg/modules/report/repositories/item"
	"github.com/michaelchandrag/botfood-go/utils"
	"github.com/xuri/excelize/v2"
)

type ReportService interface {
	ExportChannelReport(payload dto.ReportRequestPayload) (response dto.ChannelReportResponse)
	ExportBrandPromotion(payload dto.ReportRequestPayload) (response dto.PromotionReportResponse)
	ExportATPReport(payload dto.ReportRequestPayload) (response dto.ATPReportResponse)
}

type service struct {
	db database.MainDB
}

func RegisterReportService(db database.MainDB) ReportService {
	return &service{
		db: db,
	}
}

func (s *service) ExportChannelReport(payload dto.ReportRequestPayload) (response dto.ChannelReportResponse) {

	if payload.Brand.ID == 0 {
		response.Errors.AddHTTPError(400, errors.New("Brand is required"))
		return response
	}
	payloadBrandID := int(payload.Brand.ID)

	branchChannelRepository := branch_channel_repository.NewRepository(s.db)
	branchChannelFilter := branch_channel_repository.Filter{
		BrandID:   &payloadBrandID,
		BranchIDs: payload.BranchIDs,
	}
	branchChannels, err := branchChannelRepository.FindAll(branchChannelFilter)
	if err != nil {
		response.Errors.AddHTTPError(500, errors.New("Internal Server Error. Please contact our team for more information"))
		return response
	}

	itemRepository := item_repository.NewRepository(s.db)
	itemFilter := item_repository.Filter{
		BrandID:   &payloadBrandID,
		BranchIDs: payload.BranchIDs,
	}
	items, err := itemRepository.FindAll(itemFilter)
	if err != nil {
		response.Errors.AddHTTPError(500, errors.New("Internal Server Error. Please contact our team for more information"))
		return response
	}

	// excel
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	currentTime := time.Now()
	gofoodChannel := entities.BRANCH_CHANNEL_CHANNEL_GOFOOD
	grabfoodChannel := entities.BRANCH_CHANNEL_CHANNEL_GRABFOOD
	shopeefoodChannel := entities.BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD

	index, err := f.NewSheet(gofoodChannel)
	if err != nil {
		fmt.Println(err)
	}

	_, err = f.NewSheet(grabfoodChannel)
	if err != nil {
		fmt.Println(err)
	}

	_, err = f.NewSheet(shopeefoodChannel)
	if err != nil {
		fmt.Println(err)
	}

	rawGreenStyle := utils.ExcelizeStyle{
		Color: "#32CD32",
	}
	greenStyle, err := f.NewStyle(rawGreenStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawRedStyle := utils.ExcelizeStyle{
		Color: "#FF0000",
	}
	redStyle, err := f.NewStyle(rawRedStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawGrayStyle := utils.ExcelizeStyle{
		Color: "#DCDCDC",
	}
	grayStyle, err := f.NewStyle(rawGrayStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawOutletStyle := utils.ExcelizeStyle{
		Color:               "#000000",
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
		FontColor:           "#FFFFFF",
	}
	outletStyle, err := f.NewStyle(rawOutletStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawCellVerticalStyle := utils.ExcelizeStyle{
		VerticalAlignment:   "center",
		HorizontalAlignment: "center",
		Border:              "all",
	}

	cellVerticalStyle, err := f.NewStyle(rawCellVerticalStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawCellHorizontalStyle := utils.ExcelizeStyle{
		VerticalAlignment:   "center",
		HorizontalAlignment: "center",
		Border:              "all",
		TextRotation:        90,
	}

	cellHorizontalStyle, err := f.NewStyle(rawCellHorizontalStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawCellRedStyle := utils.ExcelizeStyle{
		Color:               "#FF0000",
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
	}
	cellRedStyle, err := f.NewStyle(rawCellRedStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawCellGreenStyle := utils.ExcelizeStyle{
		Color:               "#32CD32",
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
	}
	cellGreenStyle, err := f.NewStyle(rawCellGreenStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawCellGrayStyle := utils.ExcelizeStyle{
		Color:               "#DCDCDC",
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
	}
	cellGrayStyle, err := f.NewStyle(rawCellGrayStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawCellTotalStyle := utils.ExcelizeStyle{
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
		Bold:                true,
	}
	cellTotalStyle, err := f.NewStyle(rawCellTotalStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawHeaderHorizontalStyle := utils.ExcelizeStyle{
		Color:               "#FFFF00",
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
		Bold:                true,
		TextRotation:        90,
	}
	cellHeaderHorizontalStyle, err := f.NewStyle(rawHeaderHorizontalStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawAllPercentageStyle := utils.ExcelizeStyle{
		Color:               "#FFFF00",
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
		Bold:                true,
	}
	allPercentageStyle, err := f.NewStyle(rawAllPercentageStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	f.SetCellValue(gofoodChannel, "B3", "Channel")
	f.SetCellValue(gofoodChannel, "C3", gofoodChannel)
	f.SetCellValue(gofoodChannel, "B4", "Date")
	f.SetCellValue(gofoodChannel, "C4", currentTime.Format("02/Jan/2006"))
	f.SetCellValue(gofoodChannel, "B5", "Time")
	f.SetCellValue(gofoodChannel, "C5", currentTime.Format("15:04"))
	f.SetCellValue(gofoodChannel, "G4", "Buka / Aktif")
	f.SetCellValue(gofoodChannel, "G5", "Tutup / Tidak Aktif")
	f.SetCellValue(gofoodChannel, "L4", "Tidak tersedia")
	f.SetCellStyle(gofoodChannel, "G4", "G4", greenStyle)
	f.SetCellStyle(gofoodChannel, "G5", "G5", redStyle)
	f.SetCellStyle(gofoodChannel, "L4", "L4", grayStyle)
	f.SetCellValue(gofoodChannel, "B7", fmt.Sprintf("%s %s %s", gofoodChannel, payload.Brand.Name, currentTime.Format("02/Jan/2006")))
	f.SetCellValue(gofoodChannel, "B8", "Outlet")
	f.SetCellValue(gofoodChannel, "C8", "Rating")
	f.SetCellValue(gofoodChannel, "D8", "BUKA/TUTUP")
	f.SetColWidth(gofoodChannel, "B", "B", 40)
	f.SetColWidth(gofoodChannel, "C", "C", 4)
	f.SetColWidth(gofoodChannel, "D", "D", 4)
	f.SetRowHeight(gofoodChannel, 8, 150)
	f.SetCellStyle(gofoodChannel, "B8", "B8", outletStyle)
	f.SetCellStyle(gofoodChannel, "C8", "C8", cellHorizontalStyle)
	f.SetCellStyle(gofoodChannel, "D8", "D8", cellHorizontalStyle)

	f.SetCellValue(grabfoodChannel, "B3", "Channel")
	f.SetCellValue(grabfoodChannel, "C3", grabfoodChannel)
	f.SetCellValue(grabfoodChannel, "B4", "Date")
	f.SetCellValue(grabfoodChannel, "C4", currentTime.Format("02/Jan/2006"))
	f.SetCellValue(grabfoodChannel, "B5", "Time")
	f.SetCellValue(grabfoodChannel, "C5", currentTime.Format("15:04"))
	f.SetCellValue(grabfoodChannel, "G4", "Buka / Aktif")
	f.SetCellValue(grabfoodChannel, "G5", "Tutup / Tidak Aktif")
	f.SetCellValue(grabfoodChannel, "L4", "Tidak tersedia")
	f.SetCellStyle(grabfoodChannel, "G4", "G4", greenStyle)
	f.SetCellStyle(grabfoodChannel, "G5", "G5", redStyle)
	f.SetCellStyle(grabfoodChannel, "L4", "L4", grayStyle)
	f.SetCellValue(grabfoodChannel, "B7", fmt.Sprintf("%s %s %s", grabfoodChannel, payload.Brand.Name, currentTime.Format("02/Jan/2006")))
	f.SetCellValue(grabfoodChannel, "B8", "Outlet")
	f.SetCellValue(grabfoodChannel, "C8", "Rating")
	f.SetCellValue(grabfoodChannel, "D8", "BUKA/TUTUP")
	f.SetColWidth(grabfoodChannel, "B", "B", 40)
	f.SetColWidth(grabfoodChannel, "C", "C", 4)
	f.SetColWidth(grabfoodChannel, "D", "D", 4)
	f.SetRowHeight(grabfoodChannel, 8, 150)
	f.SetCellStyle(grabfoodChannel, "B8", "B8", outletStyle)
	f.SetCellStyle(grabfoodChannel, "C8", "C8", cellHorizontalStyle)
	f.SetCellStyle(grabfoodChannel, "D8", "D8", cellHorizontalStyle)

	f.SetCellValue(shopeefoodChannel, "B3", "Channel")
	f.SetCellValue(shopeefoodChannel, "C3", shopeefoodChannel)
	f.SetCellValue(shopeefoodChannel, "B4", "Date")
	f.SetCellValue(shopeefoodChannel, "C4", currentTime.Format("02/Jan/2006"))
	f.SetCellValue(shopeefoodChannel, "B5", "Time")
	f.SetCellValue(shopeefoodChannel, "C5", currentTime.Format("15:04"))
	f.SetCellValue(shopeefoodChannel, "G4", "Buka / Aktif")
	f.SetCellValue(shopeefoodChannel, "G5", "Tutup / Tidak Aktif")
	f.SetCellValue(shopeefoodChannel, "L4", "Tidak tersedia")
	f.SetCellStyle(shopeefoodChannel, "G4", "G4", greenStyle)
	f.SetCellStyle(shopeefoodChannel, "G5", "G5", redStyle)
	f.SetCellStyle(shopeefoodChannel, "L4", "L4", grayStyle)
	f.SetCellValue(shopeefoodChannel, "B7", fmt.Sprintf("%s %s %s", shopeefoodChannel, payload.Brand.Name, currentTime.Format("02/Jan/2006")))
	f.SetCellValue(shopeefoodChannel, "B8", "Outlet")
	f.SetCellValue(shopeefoodChannel, "C8", "Rating")
	f.SetCellValue(shopeefoodChannel, "D8", "BUKA/TUTUP")
	f.SetColWidth(shopeefoodChannel, "B", "B", 40)
	f.SetColWidth(shopeefoodChannel, "C", "C", 4)
	f.SetColWidth(shopeefoodChannel, "D", "D", 4)
	f.SetRowHeight(shopeefoodChannel, 8, 150)
	f.SetCellStyle(shopeefoodChannel, "B8", "B8", outletStyle)
	f.SetCellStyle(shopeefoodChannel, "C8", "C8", cellHorizontalStyle)
	f.SetCellStyle(shopeefoodChannel, "D8", "D8", cellHorizontalStyle)

	gofoodRow := 9
	grabfoodRow := 9
	shopeefoodRow := 9

	var channelReport entities.ChannelReportData
	var gofoodData entities.ChannelReport

	var gofood = make(map[string]map[string]bool)
	var grabfood = make(map[string]map[string]bool)
	var shopeefood = make(map[string]map[string]bool)

	for _, item := range items {
		inStock := false
		if item.PayloadInStock == 1 {
			inStock = true
		}

		var branchChannelID string
		branchChannelID = strconv.Itoa(item.BranchChannelID)

		if item.BranchChannelChannel == entities.BRANCH_CHANNEL_CHANNEL_GOFOOD {
			if _, ok := gofood[item.Name]; !ok {
				gofood[item.Name] = make(map[string]bool)
			}
			gofood[item.Name][branchChannelID] = inStock
		} else if item.BranchChannelChannel == entities.BRANCH_CHANNEL_CHANNEL_GRABFOOD {
			if _, ok := grabfood[item.Name]; !ok {
				grabfood[item.Name] = make(map[string]bool)
			}
			grabfood[item.Name][branchChannelID] = inStock
		} else if item.BranchChannelChannel == entities.BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD {
			if _, ok := shopeefood[item.Name]; !ok {
				shopeefood[item.Name] = make(map[string]bool)
			}
			shopeefood[item.Name][branchChannelID] = inStock
		}
	}
	channelReport.GoFoodReport = gofoodData

	columns := utils.GetExcelColumns()

	type itemDictionary struct {
		Index             int
		ItemName          string
		TotalItem         int
		TotalItemActive   int
		TotalItemInactive int
	}
	var gofoodItemIndex []itemDictionary
	var grabfoodItemIndex []itemDictionary
	var shopeefoodItemIndex []itemDictionary

	//gofood
	idxCol := 4
	idxItem := 0
	for itemName, _ := range gofood {
		f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_GOFOOD, fmt.Sprintf("%s%d", columns[idxCol], 8), itemName)
		f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_GOFOOD, fmt.Sprintf("%s%d", columns[idxCol], 8), fmt.Sprintf("%s%d", columns[idxCol], 8), cellHorizontalStyle)
		f.SetColWidth(gofoodChannel, columns[idxCol], columns[idxCol], 4)
		newDictionary := itemDictionary{
			Index:             idxItem,
			ItemName:          itemName,
			TotalItem:         0,
			TotalItemActive:   0,
			TotalItemInactive: 0,
		}
		gofoodItemIndex = append(gofoodItemIndex, newDictionary)
		idxItem++
		idxCol++
	}
	err = f.MergeCell(entities.BRANCH_CHANNEL_CHANNEL_GOFOOD, fmt.Sprintf("%s%d", columns[idxCol], 8), fmt.Sprintf("%s%d", columns[idxCol+1], 8))
	if err != nil {
		fmt.Println(err)
	}
	f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_GOFOOD, fmt.Sprintf("%s%d", columns[idxCol], 8), "Total Menu Aktif")
	f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_GOFOOD, fmt.Sprintf("%s%d", columns[idxCol], 8), fmt.Sprintf("%s%d", columns[idxCol+1], 8), cellHeaderHorizontalStyle)

	// grabfood
	idxCol = 4
	idxItem = 0
	for itemName, _ := range grabfood {
		f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_GRABFOOD, fmt.Sprintf("%s%d", columns[idxCol], 8), itemName)
		f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_GRABFOOD, fmt.Sprintf("%s%d", columns[idxCol], 8), fmt.Sprintf("%s%d", columns[idxCol], 8), cellHorizontalStyle)
		f.SetColWidth(grabfoodChannel, columns[idxCol], columns[idxCol], 4)
		newDictionary := itemDictionary{
			Index:             idxItem,
			ItemName:          itemName,
			TotalItem:         0,
			TotalItemActive:   0,
			TotalItemInactive: 0,
		}
		grabfoodItemIndex = append(grabfoodItemIndex, newDictionary)
		idxItem++
		idxCol++
	}
	err = f.MergeCell(entities.BRANCH_CHANNEL_CHANNEL_GRABFOOD, fmt.Sprintf("%s%d", columns[idxCol], 8), fmt.Sprintf("%s%d", columns[idxCol+1], 8))
	if err != nil {
		fmt.Println(err)
	}
	f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_GRABFOOD, fmt.Sprintf("%s%d", columns[idxCol], 8), "Total Menu Aktif")
	f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_GRABFOOD, fmt.Sprintf("%s%d", columns[idxCol], 8), fmt.Sprintf("%s%d", columns[idxCol+1], 8), cellHeaderHorizontalStyle)

	// shopeefood
	idxCol = 4
	idxItem = 0
	for itemName, _ := range shopeefood {
		f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD, fmt.Sprintf("%s%d", columns[idxCol], 8), itemName)
		f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD, fmt.Sprintf("%s%d", columns[idxCol], 8), fmt.Sprintf("%s%d", columns[idxCol], 8), cellHorizontalStyle)
		f.SetColWidth(shopeefoodChannel, columns[idxCol], columns[idxCol], 4)
		newDictionary := itemDictionary{
			Index:             idxItem,
			ItemName:          itemName,
			TotalItem:         0,
			TotalItemActive:   0,
			TotalItemInactive: 0,
		}
		shopeefoodItemIndex = append(shopeefoodItemIndex, newDictionary)
		idxItem++
		idxCol++
	}
	err = f.MergeCell(entities.BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD, fmt.Sprintf("%s%d", columns[idxCol], 8), fmt.Sprintf("%s%d", columns[idxCol+1], 8))
	if err != nil {
		fmt.Println(err)
	}
	f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD, fmt.Sprintf("%s%d", columns[idxCol], 8), "Total Menu Aktif")
	f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD, fmt.Sprintf("%s%d", columns[idxCol], 8), fmt.Sprintf("%s%d", columns[idxCol+1], 8), cellHeaderHorizontalStyle)

	totalGofoodItemActive := 0
	totalGofoodItem := 0
	totalGrabfoodItemActive := 0
	totalGrabfoodItem := 0
	totalShopeefoodItemActive := 0
	totalShopeefoodItem := 0
	for _, branchChannel := range branchChannels {
		rating := ""
		if branchChannel.Rating != nil {
			rating = fmt.Sprintf("%g", *branchChannel.Rating)
		}
		useRow := -1
		if branchChannel.Channel == gofoodChannel {
			useRow = gofoodRow
		} else if branchChannel.Channel == grabfoodChannel {
			useRow = grabfoodRow
		} else if branchChannel.Channel == shopeefoodChannel {
			useRow = shopeefoodRow
		}
		currentChannel := branchChannel.Channel
		f.SetCellValue(currentChannel, fmt.Sprintf("%s%d", "B", useRow), branchChannel.Name)
		f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", "B", useRow), fmt.Sprintf("%s%d", "B", useRow), cellVerticalStyle)
		f.SetCellValue(currentChannel, fmt.Sprintf("%s%d", "C", useRow), rating)
		f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", "C", useRow), fmt.Sprintf("%s%d", "C", useRow), cellVerticalStyle)
		if branchChannel.PayloadIsOpen == 1 {
			f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", "D", useRow), fmt.Sprintf("%s%d", "D", useRow), cellGreenStyle)
		} else if branchChannel.PayloadIsOpen == 0 {
			f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", "D", useRow), fmt.Sprintf("%s%d", "D", useRow), cellRedStyle)
		}

		var branchChannelIDText string
		branchChannelIDText = strconv.Itoa(branchChannel.ID)
		if branchChannel.Channel == gofoodChannel {
			idxColBc := 4
			totalItem := 0
			totalItemActive := 0
			totalAllItem := 0
			totalAllItemActive := 0
			for keyIndex, itemIndex := range gofoodItemIndex {
				if _, ok := gofood[itemIndex.ItemName][branchChannelIDText]; ok {
					if gofood[itemIndex.ItemName][branchChannelIDText] == true {
						f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%s%d", columns[idxColBc], useRow), cellGreenStyle)
						totalItemActive++
						totalAllItemActive++
						totalGofoodItemActive++
						gofoodItemIndex[keyIndex].TotalItemActive++
					} else {
						f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%s%d", columns[idxColBc], useRow), cellRedStyle)
						gofoodItemIndex[keyIndex].TotalItemInactive++
					}
					gofoodItemIndex[keyIndex].TotalItem++
					totalItem++
					totalAllItem++
					totalGofoodItem++
				} else {
					f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%s%d", columns[idxColBc], useRow), cellGrayStyle)
				}
				idxColBc++
			}
			textPercentageActive := 0
			if totalItem > 0 {
				textPercentageActive = totalItemActive * 100 / totalItem
			}
			f.SetCellValue(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%d", totalItemActive))
			f.SetCellValue(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc+1], useRow), fmt.Sprintf("%d", textPercentageActive)+"%")
			f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%s%d", columns[idxColBc], useRow), cellTotalStyle)
			f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%s%d", columns[idxColBc+1], useRow), cellTotalStyle)

			gofoodRow++
		} else if branchChannel.Channel == grabfoodChannel {
			idxColBc := 4
			totalItem := 0
			totalItemActive := 0
			totalAllItem := 0
			totalAllItemActive := 0
			for keyIndex, itemIndex := range grabfoodItemIndex {
				if _, ok := grabfood[itemIndex.ItemName][branchChannelIDText]; ok {
					if grabfood[itemIndex.ItemName][branchChannelIDText] == true {
						f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%s%d", columns[idxColBc], useRow), cellGreenStyle)
						totalItemActive++
						totalAllItemActive++
						totalGrabfoodItemActive++
						grabfoodItemIndex[keyIndex].TotalItemActive++
					} else {
						f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%s%d", columns[idxColBc], useRow), cellRedStyle)
						grabfoodItemIndex[keyIndex].TotalItemInactive++
					}
					grabfoodItemIndex[keyIndex].TotalItem++
					totalItem++
					totalAllItem++
					totalGrabfoodItem++
				} else {
					f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%s%d", columns[idxColBc], useRow), cellGrayStyle)
				}
				idxColBc++
			}
			textPercentageActive := 0
			if totalItem > 0 {
				textPercentageActive = totalItemActive * 100 / totalItem
			}
			f.SetCellValue(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%d", totalItemActive))
			f.SetCellValue(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc+1], useRow), fmt.Sprintf("%d", textPercentageActive)+"%")
			f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%s%d", columns[idxColBc], useRow), cellTotalStyle)
			f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%s%d", columns[idxColBc+1], useRow), cellTotalStyle)

			grabfoodRow++
		} else if branchChannel.Channel == shopeefoodChannel {
			idxColBc := 4
			totalItem := 0
			totalItemActive := 0
			totalAllItem := 0
			totalAllItemActive := 0
			for keyIndex, itemIndex := range shopeefoodItemIndex {
				if _, ok := shopeefood[itemIndex.ItemName][branchChannelIDText]; ok {
					if shopeefood[itemIndex.ItemName][branchChannelIDText] == true {
						f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%s%d", columns[idxColBc], useRow), cellGreenStyle)
						totalItemActive++
						totalAllItemActive++
						totalShopeefoodItemActive++
						shopeefoodItemIndex[keyIndex].TotalItemActive++
					} else {
						f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%s%d", columns[idxColBc], useRow), cellRedStyle)
						shopeefoodItemIndex[keyIndex].TotalItemInactive++
					}
					shopeefoodItemIndex[keyIndex].TotalItem++
					totalItem++
					totalAllItem++
					totalShopeefoodItem++
				} else {
					f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%s%d", columns[idxColBc], useRow), cellGrayStyle)
				}
				idxColBc++
			}
			textPercentageActive := 0
			if totalItem > 0 {
				textPercentageActive = totalItemActive * 100 / totalItem
			}
			f.SetCellValue(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%d", totalItemActive))
			f.SetCellValue(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc+1], useRow), fmt.Sprintf("%d", textPercentageActive)+"%")
			f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%s%d", columns[idxColBc], useRow), cellTotalStyle)
			f.SetCellStyle(currentChannel, fmt.Sprintf("%s%d", columns[idxColBc], useRow), fmt.Sprintf("%s%d", columns[idxColBc+1], useRow), cellTotalStyle)

			shopeefoodRow++
		}
	}

	// gofood
	textAllGofoodPercentageActive := 0
	if totalGofoodItem > 0 {
		textAllGofoodPercentageActive = totalGofoodItemActive * 100 / totalGofoodItem
	}
	f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_GOFOOD, fmt.Sprintf("%s%d", columns[len(gofoodItemIndex)+1+4], gofoodRow), fmt.Sprintf("%d", textAllGofoodPercentageActive)+"%")
	f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_GOFOOD, fmt.Sprintf("%s%d", columns[len(gofoodItemIndex)+1+4], gofoodRow), fmt.Sprintf("%s%d", columns[len(gofoodItemIndex)+1+4], gofoodRow), allPercentageStyle)

	// grabfood
	textAllGrabfoodPercentageActive := 0
	if totalGrabfoodItem > 0 {
		textAllGrabfoodPercentageActive = totalGrabfoodItemActive * 100 / totalGrabfoodItem
	}
	f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_GRABFOOD, fmt.Sprintf("%s%d", columns[len(grabfoodItemIndex)+1+4], grabfoodRow), fmt.Sprintf("%d", textAllGrabfoodPercentageActive)+"%")
	f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_GRABFOOD, fmt.Sprintf("%s%d", columns[len(grabfoodItemIndex)+1+4], grabfoodRow), fmt.Sprintf("%s%d", columns[len(grabfoodItemIndex)+1+4], grabfoodRow), allPercentageStyle)

	// shopeefood
	textAllShopeefoodPercentageActive := 0
	if totalShopeefoodItem > 0 {
		textAllShopeefoodPercentageActive = totalShopeefoodItemActive * 100 / totalShopeefoodItem
	}
	f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD, fmt.Sprintf("%s%d", columns[len(shopeefoodItemIndex)+1+4], shopeefoodRow), fmt.Sprintf("%d", textAllShopeefoodPercentageActive)+"%")
	f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD, fmt.Sprintf("%s%d", columns[len(shopeefoodItemIndex)+1+4], shopeefoodRow), fmt.Sprintf("%s%d", columns[len(shopeefoodItemIndex)+1+4], shopeefoodRow), allPercentageStyle)

	rawFooterStyle := utils.ExcelizeStyle{
		Color:               "#FFA500",
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
		Bold:                true,
	}
	footerStyle, err := f.NewStyle(rawFooterStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawFooterRedStyle := utils.ExcelizeStyle{
		Color:               "#FF0000",
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
		FontColor:           "#FFFFFF",
		Bold:                true,
	}
	footerRedStyle, err := f.NewStyle(rawFooterRedStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	// gofood
	err = f.MergeCell(entities.BRANCH_CHANNEL_CHANNEL_GOFOOD, fmt.Sprintf("%s%d", "B", gofoodRow), fmt.Sprintf("%s%d", "D", gofoodRow+1))
	if err != nil {
		fmt.Println(err)
	}
	f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_GOFOOD, fmt.Sprintf("%s%d", "B", gofoodRow), "PRESENTASI ALL OUTLET PERMENU (SOLD OUT)")
	f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_GOFOOD, fmt.Sprintf("%s%d", "B", gofoodRow), fmt.Sprintf("%s%d", "D", gofoodRow+1), footerStyle)

	idxColItem := 4
	for _, itemIndex := range gofoodItemIndex {
		f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_GOFOOD, fmt.Sprintf("%s%d", columns[idxColItem], gofoodRow), itemIndex.TotalItemInactive)
		f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_GOFOOD, fmt.Sprintf("%s%d", columns[idxColItem], gofoodRow), fmt.Sprintf("%s%d", columns[idxColItem], gofoodRow), footerStyle)
		textPercentageItemInactive := 0
		if itemIndex.TotalItemInactive > 0 {
			textPercentageItemInactive = itemIndex.TotalItemInactive * 100 / itemIndex.TotalItem
		}
		f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_GOFOOD, fmt.Sprintf("%s%d", columns[idxColItem], gofoodRow+1), fmt.Sprintf("%d", textPercentageItemInactive)+"%")
		f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_GOFOOD, fmt.Sprintf("%s%d", columns[idxColItem], gofoodRow+1), fmt.Sprintf("%s%d", columns[idxColItem], gofoodRow+1), footerRedStyle)
		idxColItem++
	}

	// grabfood
	err = f.MergeCell(entities.BRANCH_CHANNEL_CHANNEL_GRABFOOD, fmt.Sprintf("%s%d", "B", grabfoodRow), fmt.Sprintf("%s%d", "D", grabfoodRow+1))
	if err != nil {
		fmt.Println(err)
	}
	f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_GRABFOOD, fmt.Sprintf("%s%d", "B", grabfoodRow), "PRESENTASI ALL OUTLET PERMENU (SOLD OUT)")
	f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_GRABFOOD, fmt.Sprintf("%s%d", "B", grabfoodRow), fmt.Sprintf("%s%d", "D", grabfoodRow+1), footerStyle)
	idxColItem = 4
	for _, itemIndex := range grabfoodItemIndex {
		f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_GRABFOOD, fmt.Sprintf("%s%d", columns[idxColItem], grabfoodRow), itemIndex.TotalItemInactive)
		f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_GRABFOOD, fmt.Sprintf("%s%d", columns[idxColItem], grabfoodRow), fmt.Sprintf("%s%d", columns[idxColItem], grabfoodRow), footerStyle)
		textPercentageItemInactive := 0
		if itemIndex.TotalItemInactive > 0 {
			textPercentageItemInactive = itemIndex.TotalItemInactive * 100 / itemIndex.TotalItem
		}
		f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_GRABFOOD, fmt.Sprintf("%s%d", columns[idxColItem], grabfoodRow+1), fmt.Sprintf("%d", textPercentageItemInactive)+"%")
		f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_GRABFOOD, fmt.Sprintf("%s%d", columns[idxColItem], grabfoodRow+1), fmt.Sprintf("%s%d", columns[idxColItem], grabfoodRow+1), footerRedStyle)
		idxColItem++
	}

	// shopeefood
	err = f.MergeCell(entities.BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD, fmt.Sprintf("%s%d", "B", shopeefoodRow), fmt.Sprintf("%s%d", "D", shopeefoodRow+1))
	if err != nil {
		fmt.Println(err)
	}
	f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD, fmt.Sprintf("%s%d", "B", shopeefoodRow), "PRESENTASI ALL OUTLET PERMENU (SOLD OUT)")
	f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD, fmt.Sprintf("%s%d", "B", shopeefoodRow), fmt.Sprintf("%s%d", "D", shopeefoodRow+1), footerStyle)
	idxColItem = 4
	for _, itemIndex := range shopeefoodItemIndex {
		f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD, fmt.Sprintf("%s%d", columns[idxColItem], shopeefoodRow), itemIndex.TotalItemInactive)
		f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD, fmt.Sprintf("%s%d", columns[idxColItem], shopeefoodRow), fmt.Sprintf("%s%d", columns[idxColItem], shopeefoodRow), footerStyle)
		textPercentageItemInactive := 0
		if itemIndex.TotalItemInactive > 0 {
			textPercentageItemInactive = itemIndex.TotalItemInactive * 100 / itemIndex.TotalItem
		}
		f.SetCellValue(entities.BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD, fmt.Sprintf("%s%d", columns[idxColItem], shopeefoodRow+1), fmt.Sprintf("%d", textPercentageItemInactive)+"%")
		f.SetCellStyle(entities.BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD, fmt.Sprintf("%s%d", columns[idxColItem], shopeefoodRow+1), fmt.Sprintf("%s%d", columns[idxColItem], shopeefoodRow+1), footerRedStyle)
		idxColItem++
	}

	f.SetActiveSheet(index)

	err = f.DeleteSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
	}
	response.File.Excel = f

	response.Data.BranchChannels = branchChannels
	response.Data.Items = items
	response.Data.ChannelReportData = channelReport

	return response
}

func (s *service) ExportBrandPromotion(payload dto.ReportRequestPayload) (response dto.PromotionReportResponse) {
	if payload.Brand.ID == 0 {
		response.Errors.AddHTTPError(400, errors.New("Brand is required"))
		return response
	}

	payloadBrandID := int(payload.Brand.ID)
	promotionRepository := branch_channel_promotion_repository.NewRepository(s.db)
	promotionFilter := branch_channel_promotion_repository.Filter{
		BrandID: &payloadBrandID,
	}
	promotions, err := promotionRepository.FindAll(promotionFilter)
	if err != nil {
		fmt.Println("ERROR PROMO VOUCHER")
		response.Errors.AddHTTPError(500, errors.New("Internal Server Error. Please contact our team for more information"))
		return response
	}

	var promotionGroup = make(map[string][]entities.BranchChannelPromotion)
	for _, promotion := range promotions {
		if _, ok := promotionGroup[promotion.BranchChannelChannel]; !ok {
			var emptyPromotions []entities.BranchChannelPromotion
			promotionGroup[promotion.BranchChannelChannel] = emptyPromotions
		}
		promotionGroup[promotion.BranchChannelChannel] = append(promotionGroup[promotion.BranchChannelChannel], promotion)
	}

	itemRepository := item_repository.NewRepository(s.db)
	hasSellingPrice := true
	itemDiscountFilter := item_repository.Filter{
		BrandID:         &payloadBrandID,
		HasSellingPrice: &hasSellingPrice,
	}
	itemDiscounts, err := itemRepository.FindAll(itemDiscountFilter)
	if err != nil {
		fmt.Println("ERROR DISKON CORET")
		response.Errors.AddHTTPError(500, errors.New("Internal Server Error. Please contact our team for more information"))
		return response
	}
	isBundle := 1
	itemBundleFilter := item_repository.Filter{
		BrandID:  &payloadBrandID,
		IsBundle: &isBundle,
	}
	itemBundles, err := itemRepository.FindAll(itemBundleFilter)
	if err != nil {
		fmt.Println("ERROR BUNDLE")
		response.Errors.AddHTTPError(500, errors.New("Internal Server Error. Please contact our team for more information"))
		return response
	}

	// excel
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	voucherSheet := "PROMO VOUCHER"
	discountSheet := "PROMO ITEM DISKON CORET"
	bundleSheet := "PROMO ITEM BUNDLE"
	index, err := f.NewSheet(voucherSheet)
	if err != nil {
		fmt.Println(err)
	}

	rawChannelHeaderVoucherStyle := utils.ExcelizeStyle{
		Color:               "#FFA500",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
		Bold:                true,
	}
	channelHeaderVoucherStyle, err := f.NewStyle(rawChannelHeaderVoucherStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawHeaderVoucherStyle := utils.ExcelizeStyle{
		Color:               "#FFFF00",
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
		Bold:                true,
	}
	headerVoucherStyle, err := f.NewStyle(rawHeaderVoucherStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawCellVoucherStyle := utils.ExcelizeStyle{
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
	}
	cellVoucherStyle, err := f.NewStyle(rawCellVoucherStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	currentTime := time.Now()
	f.SetCellValue(voucherSheet, "A1", "Tanggal Laporan")
	f.SetCellValue(voucherSheet, "B1", currentTime.Format("02/Jan/2006"))

	f.SetColWidth(voucherSheet, "A", "A", 40)
	f.SetColWidth(voucherSheet, "B", "B", 40)
	f.SetColWidth(voucherSheet, "C", "C", 50)
	f.SetColWidth(voucherSheet, "D", "D", 50)
	f.SetColWidth(voucherSheet, "E", "E", 20)
	f.SetColWidth(voucherSheet, "F", "F", 25)
	f.SetColWidth(voucherSheet, "G", "G", 25)
	f.SetColWidth(voucherSheet, "H", "H", 25)
	idxRow := 3
	for channel, promotions := range promotionGroup {
		f.SetCellValue(voucherSheet, fmt.Sprintf("A%d", idxRow), channel)
		f.SetCellStyle(voucherSheet, fmt.Sprintf("A%d", idxRow), fmt.Sprintf("A%d", idxRow), channelHeaderVoucherStyle)

		idxRow++
		f.SetCellValue(voucherSheet, fmt.Sprintf("A%d", idxRow), "Outlet")
		f.SetCellValue(voucherSheet, fmt.Sprintf("B%d", idxRow), "Judul Voucher")
		f.SetCellValue(voucherSheet, fmt.Sprintf("C%d", idxRow), "Deskripsi")
		f.SetCellValue(voucherSheet, fmt.Sprintf("D%d", idxRow), "Catatan")
		f.SetCellValue(voucherSheet, fmt.Sprintf("E%d", idxRow), "Tipe Diskon")
		f.SetCellValue(voucherSheet, fmt.Sprintf("F%d", idxRow), "Nilai Diskon")
		f.SetCellValue(voucherSheet, fmt.Sprintf("G%d", idxRow), "Minimal Pembelian")
		f.SetCellValue(voucherSheet, fmt.Sprintf("H%d", idxRow), "Maksimal Nilai Diskon")

		f.SetCellStyle(voucherSheet, fmt.Sprintf("A%d", idxRow), fmt.Sprintf("H%d", idxRow), headerVoucherStyle)
		idxRow++

		for _, promotion := range promotions {
			pr := promotion.ToReport()
			f.SetCellValue(voucherSheet, fmt.Sprintf("A%d", idxRow), pr.BranchChannelName)
			f.SetCellValue(voucherSheet, fmt.Sprintf("B%d", idxRow), pr.Title)
			f.SetCellValue(voucherSheet, fmt.Sprintf("C%d", idxRow), pr.Description)
			f.SetCellValue(voucherSheet, fmt.Sprintf("D%d", idxRow), pr.TagsInText)
			f.SetCellValue(voucherSheet, fmt.Sprintf("E%d", idxRow), pr.DiscountType)
			f.SetCellValue(voucherSheet, fmt.Sprintf("F%d", idxRow), pr.DiscountValue)
			f.SetCellValue(voucherSheet, fmt.Sprintf("G%d", idxRow), pr.MinSpend)
			f.SetCellValue(voucherSheet, fmt.Sprintf("H%d", idxRow), pr.MaxDiscountAmount)
			f.SetCellStyle(voucherSheet, fmt.Sprintf("A%d", idxRow), fmt.Sprintf("H%d", idxRow), cellVoucherStyle)
			idxRow++
		}

		idxRow++
	}

	_, err = f.NewSheet(discountSheet)

	f.SetCellValue(discountSheet, "A1", "Tanggal Laporan")
	f.SetCellValue(discountSheet, "B1", currentTime.Format("02/Jan/2006"))

	f.SetColWidth(discountSheet, "A", "A", 30)
	f.SetColWidth(discountSheet, "B", "B", 30)
	f.SetColWidth(discountSheet, "C", "C", 40)
	f.SetColWidth(discountSheet, "D", "D", 40)
	f.SetColWidth(discountSheet, "E", "E", 25)
	f.SetColWidth(discountSheet, "F", "F", 25)

	idxRow = 3
	f.SetCellValue(discountSheet, fmt.Sprintf("A%d", idxRow), "Outlet")
	f.SetCellValue(discountSheet, fmt.Sprintf("B%d", idxRow), "Channel")
	f.SetCellValue(discountSheet, fmt.Sprintf("C%d", idxRow), "Nama Item")
	f.SetCellValue(discountSheet, fmt.Sprintf("D%d", idxRow), "Deskripsi")
	f.SetCellValue(discountSheet, fmt.Sprintf("E%d", idxRow), "Harga Normal")
	f.SetCellValue(discountSheet, fmt.Sprintf("F%d", idxRow), "Harga Diskon")
	f.SetCellStyle(discountSheet, fmt.Sprintf("A%d", idxRow), fmt.Sprintf("F%d", idxRow), headerVoucherStyle)
	idxRow++
	for _, item := range itemDiscounts {
		ir := item.ToReport()
		f.SetCellValue(discountSheet, fmt.Sprintf("A%d", idxRow), ir.BranchChannelName)
		f.SetCellValue(discountSheet, fmt.Sprintf("B%d", idxRow), ir.BranchChannelChannel)
		f.SetCellValue(discountSheet, fmt.Sprintf("C%d", idxRow), ir.Name)
		f.SetCellValue(discountSheet, fmt.Sprintf("D%d", idxRow), ir.Description)
		f.SetCellValue(discountSheet, fmt.Sprintf("E%d", idxRow), ir.PriceInText)
		f.SetCellValue(discountSheet, fmt.Sprintf("F%d", idxRow), ir.SellingPriceInText)

		f.SetCellStyle(discountSheet, fmt.Sprintf("A%d", idxRow), fmt.Sprintf("F%d", idxRow), cellVoucherStyle)
		idxRow++
	}

	_, err = f.NewSheet(bundleSheet)

	f.SetCellValue(bundleSheet, "A1", "Tanggal Laporan")
	f.SetCellValue(bundleSheet, "B1", currentTime.Format("02/Jan/2006"))

	f.SetColWidth(bundleSheet, "A", "A", 30)
	f.SetColWidth(bundleSheet, "B", "B", 30)
	f.SetColWidth(bundleSheet, "C", "C", 40)
	f.SetColWidth(bundleSheet, "D", "D", 40)
	f.SetColWidth(bundleSheet, "E", "E", 25)
	f.SetColWidth(bundleSheet, "F", "F", 25)

	idxRow = 3
	f.SetCellValue(bundleSheet, fmt.Sprintf("A%d", idxRow), "Outlet")
	f.SetCellValue(bundleSheet, fmt.Sprintf("B%d", idxRow), "Channel")
	f.SetCellValue(bundleSheet, fmt.Sprintf("C%d", idxRow), "Nama Item")
	f.SetCellValue(bundleSheet, fmt.Sprintf("D%d", idxRow), "Deskripsi")
	f.SetCellValue(bundleSheet, fmt.Sprintf("E%d", idxRow), "Variant")
	f.SetCellValue(bundleSheet, fmt.Sprintf("F%d", idxRow), "Harga Normal")
	f.SetCellValue(bundleSheet, fmt.Sprintf("G%d", idxRow), "Harga Diskon")
	f.SetCellStyle(bundleSheet, fmt.Sprintf("A%d", idxRow), fmt.Sprintf("G%d", idxRow), headerVoucherStyle)
	idxRow++

	type variantDict struct {
		IndexRow     int
		VariantNames []string
	}

	variantMap := make(map[string]variantDict)
	var itemIds []int

	for _, item := range itemBundles {
		ir := item.ToReport()

		if item.BranchChannelChannel == entities.BRANCH_CHANNEL_CHANNEL_GOFOOD {
			idString := strconv.Itoa(item.ID)
			variantMap[idString] = variantDict{
				IndexRow: idxRow,
			}
			itemIds = append(itemIds, item.ID)
		}

		f.SetCellValue(bundleSheet, fmt.Sprintf("A%d", idxRow), ir.BranchChannelName)
		f.SetCellValue(bundleSheet, fmt.Sprintf("B%d", idxRow), ir.BranchChannelChannel)
		f.SetCellValue(bundleSheet, fmt.Sprintf("C%d", idxRow), ir.Name)
		f.SetCellValue(bundleSheet, fmt.Sprintf("D%d", idxRow), ir.Description)
		f.SetCellValue(bundleSheet, fmt.Sprintf("F%d", idxRow), ir.PriceInText)
		f.SetCellValue(bundleSheet, fmt.Sprintf("G%d", idxRow), ir.SellingPriceInText)

		f.SetCellStyle(bundleSheet, fmt.Sprintf("A%d", idxRow), fmt.Sprintf("G%d", idxRow), cellVoucherStyle)
		idxRow++
	}

	if len(itemBundles) > 0 {
		variants, err := itemRepository.FindVariantByItemIDs(itemIds)
		if err != nil {
			response.Errors.AddHTTPError(500, errors.New("Internal Server Error. Please contact our team for more information"))
			return response
		}
		for _, variant := range variants {
			itemId := strconv.Itoa(variant.ItemID)
			if entry, ok := variantMap[itemId]; ok {
				entry.VariantNames = append(entry.VariantNames, variant.Name)
				variantMap[itemId] = entry
			}
		}

		for _, val := range variantMap {
			f.SetCellValue(bundleSheet, fmt.Sprintf("E%d", val.IndexRow), strings.Join(val.VariantNames, ","))
		}
	}

	f.SetActiveSheet(index)

	err = f.DeleteSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
	}
	response.File.Excel = f

	response.Data.Promotions = promotions
	response.Data.ItemDiscounts = itemDiscounts
	response.Data.ItemBundles = itemBundles
	return response
}

func (s *service) ExportATPReport(payload dto.ReportRequestPayload) (response dto.ATPReportResponse) {

	if payload.Brand.ID == 0 {
		response.Errors.AddHTTPError(400, errors.New("Brand is required"))
		return response
	}
	payloadBrandID := int(payload.Brand.ID)

	bcReportRepository := branch_channel_availability_report_repository.NewRepository(s.db)
	bcReportFilter := branch_channel_availability_report_repository.Filter{
		BrandID: &payloadBrandID,
		Date:    payload.Date,
	}
	bcReports, err := bcReportRepository.FindAll(bcReportFilter)
	if err != nil {
		response.Errors.AddHTTPError(500, errors.New("Internal Server Error. Please contact our team for more information"))
		return response
	}
	fmt.Println(bcReports)

	// excel
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	currentTime, err := time.Parse("2006-01-02", payload.Date)
	if err != nil {
		response.Errors.AddHTTPError(500, errors.New("Internal Server Error. Please contact our team for more information"))
		return response
	}
	gofoodChannel := entities.BRANCH_CHANNEL_CHANNEL_GOFOOD
	grabfoodChannel := entities.BRANCH_CHANNEL_CHANNEL_GRABFOOD
	shopeefoodChannel := entities.BRANCH_CHANNEL_CHANNEL_SHOPEEFOOD

	index, err := f.NewSheet(gofoodChannel)
	if err != nil {
		fmt.Println(err)
	}

	_, err = f.NewSheet(grabfoodChannel)
	if err != nil {
		fmt.Println(err)
	}

	_, err = f.NewSheet(shopeefoodChannel)
	if err != nil {
		fmt.Println(err)
	}

	rawGreenStyle := utils.ExcelizeStyle{
		Color: "#32CD32",
	}
	greenStyle, err := f.NewStyle(rawGreenStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawRedStyle := utils.ExcelizeStyle{
		Color: "#FF0000",
	}
	redStyle, err := f.NewStyle(rawRedStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawGrayStyle := utils.ExcelizeStyle{
		Color: "#DCDCDC",
	}
	grayStyle, err := f.NewStyle(rawGrayStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawYellowStyle := utils.ExcelizeStyle{
		Color: "#FFFF00",
	}
	yellowStyle, err := f.NewStyle(rawYellowStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawHeaderStyle := utils.ExcelizeStyle{
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
	}
	headerStyle, err := f.NewStyle(rawHeaderStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	/* rawOutletStyle := utils.ExcelizeStyle{
		Color:               "#000000",
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
		FontColor:           "#FFFFFF",
	}
	outletStyle, err := f.NewStyle(rawOutletStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	} */

	/* rawCellVerticalStyle := utils.ExcelizeStyle{
		VerticalAlignment:   "center",
		HorizontalAlignment: "center",
		Border:              "all",
	}

	cellVerticalStyle, err := f.NewStyle(rawCellVerticalStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawCellHorizontalStyle := utils.ExcelizeStyle{
		VerticalAlignment:   "center",
		HorizontalAlignment: "center",
		Border:              "all",
		TextRotation:        90,
	}

	cellHorizontalStyle, err := f.NewStyle(rawCellHorizontalStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawCellRedStyle := utils.ExcelizeStyle{
		Color:               "#FF0000",
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
	}
	cellRedStyle, err := f.NewStyle(rawCellRedStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawCellGreenStyle := utils.ExcelizeStyle{
		Color:               "#32CD32",
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
	}
	cellGreenStyle, err := f.NewStyle(rawCellGreenStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawCellGrayStyle := utils.ExcelizeStyle{
		Color:               "#DCDCDC",
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
	}
	cellGrayStyle, err := f.NewStyle(rawCellGrayStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawCellTotalStyle := utils.ExcelizeStyle{
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
		Bold:                true,
	}
	cellTotalStyle, err := f.NewStyle(rawCellTotalStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawHeaderHorizontalStyle := utils.ExcelizeStyle{
		Color:               "#FFFF00",
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
		Bold:                true,
		TextRotation:        90,
	}
	cellHeaderHorizontalStyle, err := f.NewStyle(rawHeaderHorizontalStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	}

	rawAllPercentageStyle := utils.ExcelizeStyle{
		Color:               "#FFFF00",
		Border:              "all",
		HorizontalAlignment: "center",
		VerticalAlignment:   "center",
		Bold:                true,
	}
	allPercentageStyle, err := f.NewStyle(rawAllPercentageStyle.GenerateStyle())
	if err != nil {
		fmt.Println(err)
	} */

	f.SetCellValue(gofoodChannel, "A1", "Channel")
	f.SetCellValue(gofoodChannel, "B1", gofoodChannel)
	f.SetCellValue(gofoodChannel, "A2", "Date")
	f.SetCellValue(gofoodChannel, "B2", currentTime.Format("02/Jan/2006"))
	f.SetCellValue(gofoodChannel, "A3", "NB: Jika salah satu item availability menunjukkan 10% - maka item tersebut available (aktif) 10% dari total 660 menit toko tersebut buka, yakni selama 66 menit.")
	f.SetCellValue(gofoodChannel, "E1", "100%")
	f.SetCellValue(gofoodChannel, "F1", "99%-60%")
	f.SetCellValue(gofoodChannel, "G1", "59%-0%")
	f.SetCellValue(gofoodChannel, "H1", "N/A")
	f.SetCellStyle(gofoodChannel, "E1", "E1", greenStyle)
	f.SetCellStyle(gofoodChannel, "F1", "F1", redStyle)
	f.SetCellStyle(gofoodChannel, "G1", "G1", yellowStyle)
	f.SetCellStyle(gofoodChannel, "H1", "H1", grayStyle)

	f.SetCellValue(gofoodChannel, "A5", "Outlet")
	f.SetCellValue(gofoodChannel, "A6", "Jam Operasional")
	f.SetCellValue(gofoodChannel, "A7", "Timeline Botfood")
	f.SetCellValue(gofoodChannel, "A8", "Durasi Outlet Buka (menit)")
	f.SetCellValue(gofoodChannel, "A9", "Performa Item")
	f.SetCellStyle(gofoodChannel, "A5", "A9", headerStyle)

	f.SetColWidth(gofoodChannel, "A", "A", 25)

	excelColumns := utils.GetExcelColumns()
	gofoodIndex := 1
	bcDictionary := make(map[string]map[string]int)
	bcDictionary[gofoodChannel] = make(map[string]int)
	for _, report := range bcReports {
		if _, k := bcDictionary[report.BranchChannelChannel]; k {
			if _, ok := bcDictionary[report.BranchChannelChannel][report.BranchChannelName]; !ok {
				bcDictionary[report.BranchChannelChannel][report.BranchChannelName] = gofoodIndex
				report.ToText()
				f.SetCellValue(report.BranchChannelChannel, fmt.Sprintf("%s%d", excelColumns[gofoodIndex], 5), report.BranchChannelName)
				f.SetCellStyle(report.BranchChannelChannel, fmt.Sprintf("%s%d", excelColumns[gofoodIndex], 5), fmt.Sprintf("%s%d", excelColumns[gofoodIndex], 5), headerStyle)
				f.SetCellValue(report.BranchChannelChannel, fmt.Sprintf("%s%d", excelColumns[gofoodIndex], 6), report.OperationalHours)
				f.SetCellStyle(report.BranchChannelChannel, fmt.Sprintf("%s%d", excelColumns[gofoodIndex], 6), fmt.Sprintf("%s%d", excelColumns[gofoodIndex], 6), headerStyle)
				f.SetCellValue(report.BranchChannelChannel, fmt.Sprintf("%s%d", excelColumns[gofoodIndex], 7), report.Timeline)
				f.SetCellStyle(report.BranchChannelChannel, fmt.Sprintf("%s%d", excelColumns[gofoodIndex], 7), fmt.Sprintf("%s%d", excelColumns[gofoodIndex], 7), headerStyle)
				f.SetCellValue(report.BranchChannelChannel, fmt.Sprintf("%s%d", excelColumns[gofoodIndex], 8), report.ActiveTime)
				f.SetCellStyle(report.BranchChannelChannel, fmt.Sprintf("%s%d", excelColumns[gofoodIndex], 8), fmt.Sprintf("%s%d", excelColumns[gofoodIndex], 8), headerStyle)
				f.SetCellValue(report.BranchChannelChannel, fmt.Sprintf("%s%d", excelColumns[gofoodIndex], 9), report.ItemAvailabilityPercentageText)
				f.SetCellStyle(report.BranchChannelChannel, fmt.Sprintf("%s%d", excelColumns[gofoodIndex], 9), fmt.Sprintf("%s%d", excelColumns[gofoodIndex], 9), headerStyle)
				f.SetColWidth(gofoodChannel, excelColumns[gofoodIndex], excelColumns[gofoodIndex], 40)
				gofoodIndex++
			}
		}
	}

	f.SetActiveSheet(index)

	err = f.DeleteSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
	}
	response.File.Excel = f

	return response
}
