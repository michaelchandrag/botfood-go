package repositories

import (
	"database/sql"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/entities"
	"github.com/michaelchandrag/botfood-go/utils"
	bqb "github.com/nullism/bqb"
)

type Repository interface {
	FindPaginated(filter Filter) (result PaginatedData, err error)
	FindOne(filter Filter) (branchChannel entities.BranchChannel, err error)
	FindWithCurrentShift(listBcIds []int) (branchChannels []entities.BranchChannel, err error)
}

type repository struct {
	db database.MainDB
}

type Filter struct {
	ID *int

	BrandID   *int
	SortKey   string
	SortValue string

	Name    string
	Keyword string
	Channel string
	IsOpen  *int

	Page *int
	Data *int
}

type PaginatedData struct {
	Data        []entities.BranchChannel `json:"data"`
	CurrentPage int                      `json:"current_page"`
	LimitData   int                      `json:"limit_data"`
	TotalPage   int                      `json:"total_page"`
	TotalData   int                      `json:"total_data"`
}

func NewRepository(db database.MainDB) Repository {
	return &repository{
		db: db,
	}
}

func getQueryBuilder() string {
	query := fmt.Sprintf(`
		SELECT
			branch_channels.id,
			branch_channels.brand_id,
			branch_channels.name,
			branch_channels.channel,
			branch_channels.is_open,
			branch_channels.rating
		FROM
			branch_channels
	`)
	return query
}

func getTotalQueryBuilder() string {
	query := fmt.Sprintf(`
		SELECT
			count(*) as total
		FROM
			branch_channels
	`)
	return query
}

func generateFilter(filter Filter) string {

	where := bqb.Optional("WHERE")

	where.And("branch_channels.deleted_at IS NULL")

	if filter.ID != nil {
		where.And("branch_channels.id = ?", *filter.ID)
	}

	if filter.BrandID != nil {
		where.And("branch_channels.brand_id = ?", *filter.BrandID)
	}

	if filter.Channel != "" {
		where.And("branch_channels.channel = ?", filter.Channel)
	}

	if filter.Name != "" {
		name := "%" + filter.Name + "%"
		where.And("branch_channels.name LIKE ?", name)
	}

	if filter.Keyword != "" {
		q := "%" + filter.Keyword + "%"
		where.And(`(branch_channels.name LIKE ? OR branch_channels.channel LIKE ?)`, q, q)
	}

	if filter.IsOpen != nil {
		where.And("branch_channels.is_open = ?", *filter.IsOpen)
	}

	queryWhere, err := bqb.New("?", where).ToRaw()
	if err != nil {
		fmt.Println(err)
	}

	return queryWhere
}

func addSort(filter Filter) string {
	sortKeys := "branch_channels.created_at"
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

func (r *repository) FindOne(filter Filter) (branchChannel entities.BranchChannel, err error) {
	queryBuilder := getQueryBuilder()
	queryWhere := generateFilter(filter)
	formattedQuery := queryBuilder + queryWhere
	err = r.db.GetDB().QueryRowx(formattedQuery).StructScan(&branchChannel)
	if err == sql.ErrNoRows {
		return branchChannel, nil
	} else if err != nil {
		fmt.Println(err)
	}
	return branchChannel, err
}

func (r *repository) FindWithCurrentShift(listBcIds []int) (branchChannels []entities.BranchChannel, err error) {
	if len(listBcIds) <= 0 {
		return branchChannels, err
	}
	now := time.Now()
	query := fmt.Sprintf(`
		SELECT
			branch_channels.id,
			branch_channels.name,
			branch_channels.channel,
			branch_channels.is_open,
			branch_channel_shifts.id as branch_channel_shift_id,
			branch_channel_shifts.day as branch_channel_shift_day,
			branch_channel_shifts.open_time as branch_channel_shift_open_time,
			branch_channel_shifts.close_time as branch_channel_shift_close_time
		FROM branch_channels
		LEFT JOIN branch_channel_shifts ON branch_channel_shifts.branch_channel_id = branch_channels.id AND branch_channel_shifts.day = %d AND branch_channel_shifts.deleted_at IS NULL AND branch_channel_shifts.open_time <= '%s' AND '%s' <= branch_channel_shifts.close_time
	`, utils.GetCurrentDayOfWeek(), now.Format("15:04:05"), now.Format("15:04:05"))
	where := bqb.Optional("WHERE")
	where.And("branch_channels.deleted_at IS NULL")
	if len(listBcIds) > 0 {
		where.And(`branch_channels.id IN (?)`, listBcIds)
	}

	queryWhere, err := bqb.New("?", where).ToRaw()
	formattedQuery := query + queryWhere
	err = r.db.GetDB().Select(&branchChannels, formattedQuery)
	if err != nil {
		fmt.Println(err)
	}
	return branchChannels, err
}
