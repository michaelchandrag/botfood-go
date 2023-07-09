package services

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/pkg/modules/report/dto"
	"github.com/michaelchandrag/botfood-go/pkg/modules/report/entities"
	branch_channel_repository "github.com/michaelchandrag/botfood-go/pkg/modules/report/repositories/branch_channel"
	item_repository "github.com/michaelchandrag/botfood-go/pkg/modules/report/repositories/item"
	"github.com/michaelchandrag/botfood-go/utils"
	"github.com/xuri/excelize/v2"
)

type ReportService interface {
	ExportChannelReport(payload dto.ReportRequestPayload) (response dto.ChannelReportResponse)
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
