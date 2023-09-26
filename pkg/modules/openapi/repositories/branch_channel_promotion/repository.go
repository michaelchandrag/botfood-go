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
	ID              *int
	BrandID         *int
	BranchChannelID *int
	DiscountType    string
	SortKey         string
	SortValue       string

	Keyword              string
	BranchChannelChannel string
	BranchChannelName    string

	Page *int
	Data *int
}

type PaginatedData struct {
	Data        []entities.BranchChannelPromotion `json:"data"`
	CurrentPage int                               `json:"current_page"`
	LimitData   int                               `json:"limit_data"`
	TotalPage   int                               `json:"total_page"`
	TotalData   int                               `json:"total_data"`
}

func NewRepository(db database.MainDB) Repository {
	return &repository{
		db: db,
	}
}

func getQueryBuilder() string {
	query := fmt.Sprintf(`
		SELECT
			branch_channel_promotions.id,
			branch_channel_promotions.slug,
			branch_channel_promotions.branch_channel_id,
			branch_channels.name as branch_channel_name,
			branch_channels.channel as branch_channel_channel,
			branch_channel_promotions.title,
			branch_channel_promotions.description,
			branch_channel_promotions.tags,
			branch_channel_promotions.discount_type,
			branch_channel_promotions.discount_value,
			branch_channel_promotions.min_spend,
			branch_channel_promotions.max_discount_amount,
			branch_channel_promotions.created_at,
			branch_channel_promotions.updated_at
		FROM
			branch_channel_promotions
		JOIN branch_channels ON branch_channels.id = branch_channel_promotions.branch_channel_id AND branch_channels.deleted_at IS NULL
	`)
	return query
}

func getTotalQueryBuilder() string {
	query := fmt.Sprintf(`
		SELECT
			count(*) as total
		FROM
			branch_channel_promotions
		JOIN branch_channels ON branch_channels.id = branch_channel_promotions.branch_channel_id AND branch_channels.deleted_at IS NULL
	`)
	return query
}

func generateFilter(filter Filter) string {

	where := bqb.Optional("WHERE")

	where.And("branch_channel_promotions.deleted_at IS NULL")

	if filter.ID != nil {
		where.And("branch_channel_promotions.id = ?", *filter.ID)
	}

	if filter.BrandID != nil {
		where.And("branch_channels.brand_id = ?", *filter.BrandID)
	}

	if filter.BranchChannelID != nil {
		where.And("branch_channels.id = ?", *filter.BranchChannelID)
	}

	if filter.BranchChannelChannel != "" {
		where.And("branch_channels.channel = ?", filter.BranchChannelChannel)
	}

	if filter.BranchChannelName != "" {
		name := "%" + filter.BranchChannelName + "%"
		where.And("branch_channels.name LIKE ?", name)
	}

	if filter.DiscountType != "" {
		where.And("branch_channel_promotions.discount_type = ?", filter.DiscountType)
	}

	if filter.Keyword != "" {
		q := "%" + filter.Keyword + "%"
		where.And(`(branch_channel_promotions.title LIKE ? OR branch_channels.channel LIKE ? OR branch_channels.name LIKE ?)`, q, q, q)
	}

	queryWhere, err := bqb.New("?", where).ToRaw()
	if err != nil {
		fmt.Println(err)
	}

	return queryWhere
}

func addSort(filter Filter) string {
	sortKeys := "branch_channel_promotions.created_at"
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
