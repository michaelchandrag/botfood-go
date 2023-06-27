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
	BranchIDs []int
	SortKey   string
	SortValue string

	BranchChannelID      *int
	BranchChannelName    string
	BranchChannelChannel string
	Keyword              string
	Rating               *int
	WithComment          bool
	WithImages           bool
	IsReviewed           bool
	WithMerchantReply    bool

	FromCreatedAt  string
	UntilCreatedAt string

	Page *int
	Data *int
}

type PaginatedData struct {
	Data        []entities.Review `json:"data"`
	CurrentPage int               `json:"current_page"`
	LimitData   int               `json:"limit_data"`
	TotalPage   int               `json:"total_page"`
	TotalData   int               `json:"total_data"`
}

func NewRepository(db database.MainDB) Repository {
	return &repository{
		db: db,
	}
}

func getQueryBuilder() string {
	query := fmt.Sprintf(`
		SELECT
			reviews.id,
			reviews.slug as slug,
			reviews.branch_channel_id as branch_channel_id,
			branch_channels.name as branch_channel_name,
			branch_channels.channel as branch_channel_channel,
			branch_channels.brand_id as brand_id,
			reviews.rating as rating,
			reviews.order_id,
			reviews.comment as comment,
			reviews.images as images,
			reviews.item_name as item_name,
			reviews.merchant_reply as merchant_reply,
			reviews.merchant_reply_at as merchant_reply_at,
			reviews.is_reviewed,
			reviews.ordered_at,
			reviews.created_at,
			reviews.updated_at,
			reviews.deleted_at
		FROM
			reviews
		JOIN branch_channels ON reviews.branch_channel_id = branch_channels.id AND branch_channels.deleted_at IS NULL
	`)
	return query
}

func getTotalQueryBuilder() string {
	query := fmt.Sprintf(`
		SELECT
			count(*) as total
		FROM
			reviews
		JOIN branch_channels ON reviews.branch_channel_id = branch_channels.id AND branch_channels.deleted_at IS NULL
	`)
	return query
}

func generateFilter(filter Filter) string {

	where := bqb.Optional("WHERE")

	where.And("reviews.deleted_at IS NULL")

	if filter.BrandID != nil {
		where.And("branch_channels.brand_id = ?", *filter.BrandID)
	}

	if filter.BranchChannelID != nil {
		where.And("reviews.branch_channel_id = ?", *filter.BranchChannelID)
	}

	if filter.BranchChannelChannel != "" {
		where.And("branch_channels.channel = ?", filter.BranchChannelChannel)
	}

	if filter.BranchChannelName != "" {
		q := "%" + filter.BranchChannelName + "%"
		where.And("branch_channels.name LIKE ?", q)
	}

	if filter.Rating != nil && *filter.Rating > 0 {
		where.And("reviews.rating = ?", *filter.Rating)
	}

	if filter.Keyword != "" {
		q := "%" + filter.Keyword + "%"
		where.And(`(branch_channels.name LIKE ? OR branch_channels.channel LIKE ? OR reviews.item_name LIKE ? or reviews.comment LIKE ?)`, q, q, q, q)
	}

	if filter.WithComment {
		where.And(`reviews.comment IS NOT NULL`)
	}

	if filter.WithImages {
		where.And(`reviews.images IS NOT NULL`)
	}

	if filter.IsReviewed {
		where.And(`reviews.is_reviewed = 1`)
	}

	if filter.WithMerchantReply {
		where.And(`reviews.merchant_reply IS NOT NULL`)
	}

	if filter.FromCreatedAt != "" {
		where.And(`reviews.created_at >= ?`, filter.FromCreatedAt)
	}

	if filter.UntilCreatedAt != "" {
		where.And(`reviews.created_at <= ?`, filter.UntilCreatedAt)
	}

	if len(filter.BranchIDs) > 0 {
		where.And(`branch_channels.branch_id IN (?)`, filter.BranchIDs)
	}

	queryWhere, err := bqb.New("?", where).ToRaw()
	if err != nil {
		fmt.Println(err)
	}

	return queryWhere
}

func addSort(filter Filter) string {
	sortKeys := "reviews.created_at"
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
