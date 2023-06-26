package repositories

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/entities"
	bqb "github.com/nullism/bqb"
)

type Repository interface {
	FindPaginated(filter Filter) (result PaginatedData, err error)
}

type repository struct {
	db database.MainDB
}

type Filter struct {
	BrandID   *int
	SortKey   string
	SortValue string

	BranchChannelID *int
	Date            string

	FromCreatedAt  string
	UntilCreatedAt string

	Page *int
	Data *int
}

type PaginatedData struct {
	Data        []entities.BranchChannelAvailabilityReport `json:"data"`
	CurrentPage int                                        `json:"current_page"`
	LimitData   int                                        `json:"limit_data"`
	TotalPage   int                                        `json:"total_page"`
	TotalData   int                                        `json:"total_data"`
}

func NewRepository(db database.MainDB) Repository {
	return &repository{
		db: db,
	}
}

func getQueryBuilder() string {
	query := fmt.Sprintf(`
		SELECT
			branch_channel_availability_reports.id,
			branch_channel_availability_reports.branch_channel_id,
			branch_channels.name as branch_channel_name,
			branch_channels.channel as branch_channel_channel,
			branch_channel_availability_reports.date,
			branch_channel_availability_reports.active_time,
			branch_channel_availability_reports.timeline,
			branch_channel_availability_reports.average_items_active_time_total,
			branch_channel_availability_reports.operational_hours,
			branch_channel_availability_reports.inactive_time,
			branch_channel_availability_reports.operational_hours_duration,
			branch_channel_availability_reports.open_state,
			branch_channel_availability_reports.close_state,
			branch_channel_availability_reports.has_break,
			branch_channel_availability_reports.items_availability_percentage,
			branch_channel_availability_reports.inactive_timeline
		FROM
			branch_channel_availability_reports
		JOIN branch_channels ON branch_channel_availability_reports.branch_channel_id = branch_channels.id
	`)
	return query
}

func getTotalQueryBuilder() string {
	query := fmt.Sprintf(`
		SELECT
			COUNT(*)
		FROM
			branch_channel_availability_reports
		JOIN branch_channels ON branch_channel_availability_reports.branch_channel_id = branch_channels.id
	`)
	return query
}

func generateFilter(filter Filter) string {

	where := bqb.Optional("WHERE")

	where.And("branch_channel_availability_reports.deleted_at IS NULL")

	if filter.BrandID != nil {
		where.And("branch_channels.brand_id = ?", *filter.BrandID)
	}

	if filter.BranchChannelID != nil {
		where.And("branch_channel_availability_reports.branch_channel_id = ?", *filter.BranchChannelID)
	}

	if filter.Date != "" {
		where.And(`branch_channel_availability_reports.date = ?`, filter.Date)
	}

	if filter.FromCreatedAt != "" {
		where.And(`branch_channel_availability_reports.created_at >= ?`, filter.FromCreatedAt)
	}

	if filter.UntilCreatedAt != "" {
		where.And(`branch_channel_availability_reports.created_at <= ?`, filter.UntilCreatedAt)
	}

	queryWhere, err := bqb.New("?", where).ToRaw()
	if err != nil {
		fmt.Println(err)
	}

	return queryWhere
}

func addSort(filter Filter) string {
	sortKeys := "branch_channel_availability_reports.id"
	sortValue := "DESC"
	if filter.SortKey != "" {
		sortKeys = filter.SortKey
	}
	if strings.ToUpper(filter.SortValue) == "ASC" || strings.ToUpper(filter.SortValue) == "DESC" {
		sortValue = strings.ToUpper(filter.SortValue)
	}
	sortQuery := fmt.Sprintf(" ORDER BY %s %s", sortKeys, sortValue)
	return sortQuery
}

func (r *repository) FindPaginated(filter Filter) (result PaginatedData, err error) {
	// total
	var wg sync.WaitGroup
	queryWhere := generateFilter(filter)

	totalQb := getTotalQueryBuilder()
	formattedTotalQuery := totalQb + queryWhere
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = r.db.GetDB().Get(&result.TotalData, formattedTotalQuery)
		if err != nil {
			fmt.Println(err)
		}
	}()

	// fetch
	currentPage := 1
	offset := 0
	limit := 10

	if filter.Data != nil && *filter.Data > 0 {
		limit = *filter.Data
	}

	if filter.Page != nil && *filter.Page > 1 {
		currentPage = *filter.Page
		offset = currentPage*limit - limit
	}

	queryPagination := fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	queryBuilder := getQueryBuilder()
	formattedQuery := queryBuilder + queryWhere + addSort(filter) + queryPagination
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = r.db.GetDB().Select(&result.Data, formattedQuery)
		if err != nil {
			fmt.Println(err)
		}
	}()

	wg.Wait()

	result.CurrentPage = currentPage
	result.LimitData = limit
	totalPage := float64(result.TotalData) / float64(limit)
	result.TotalPage = int(math.Ceil(totalPage))

	return result, err
}
