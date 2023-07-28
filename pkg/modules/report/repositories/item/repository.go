package repositories

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/pkg/modules/report/entities"
	bqb "github.com/nullism/bqb"
)

type Repository interface {
	FindPaginated(filter Filter) (result PaginatedData, err error)
	FindAll(filter Filter) (items []entities.Item, err error)
	FindVariantByItemIDs(ids []int) (variants []entities.Variant, err error)
}

type repository struct {
	db database.MainDB
}

type Filter struct {
	ID              *int
	BrandID         *int
	BranchChannelID *int
	SortKey         string
	SortValue       string

	BranchIDs            []int
	Name                 string
	Keyword              string
	BranchChannelChannel string
	BranchChannelName    string
	InStock              *int
	IsBundle             *int

	HasSellingPrice *bool

	Page *int
	Data *int
}

type PaginatedData struct {
	Data        []entities.Item `json:"data"`
	CurrentPage int             `json:"current_page"`
	LimitData   int             `json:"limit_data"`
	TotalPage   int             `json:"total_page"`
	TotalData   int             `json:"total_data"`
}

func NewRepository(db database.MainDB) Repository {
	return &repository{
		db: db,
	}
}

func getQueryBuilder() string {
	query := fmt.Sprintf(`
		SELECT
			items.id,
			items.slug,
			items.branch_channel_id,
			branch_channels.name as branch_channel_name,
			branch_channels.channel as branch_channel_channel,
			items.name,
			items.description,
			items.in_stock,
			items.price,
			items.selling_price,
			items.image_url,
			items.last_not_active_at,
			items.created_at,
			items.updated_at
		FROM
			items
		JOIN branch_channels ON branch_channels.id = items.branch_channel_id AND branch_channels.deleted_at IS NULL
	`)
	return query
}

func getTotalQueryBuilder() string {
	query := fmt.Sprintf(`
		SELECT
			count(*) as total
		FROM
			items
		JOIN branch_channels ON branch_channels.id = items.branch_channel_id AND branch_channels.deleted_at IS NULL
	`)
	return query
}

func generateFilter(filter Filter) string {

	where := bqb.Optional("WHERE")

	where.And("items.deleted_at IS NULL")

	if filter.ID != nil {
		where.And("items.id = ?", *filter.ID)
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

	if filter.Name != "" {
		name := "%" + filter.Name + "%"
		where.And("items.name LIKE ?", name)
	}

	if filter.BranchChannelName != "" {
		name := "%" + filter.BranchChannelName + "%"
		where.And("branch_channels.name LIKE ?", name)
	}

	if filter.Keyword != "" {
		q := "%" + filter.Keyword + "%"
		where.And(`(items.name LIKE ? OR branch_channels.channel LIKE ? OR branch_channels.name LIKE ?)`, q, q, q)
	}

	if filter.InStock != nil {
		where.And("items.in_stock = ?", *filter.InStock)
	}

	if filter.IsBundle != nil {
		where.And("items.is_bundle = ?", *filter.IsBundle)
	}

	if len(filter.BranchIDs) > 0 {
		where.And(`branch_channels.branch_id IN (?)`, filter.BranchIDs)
	}

	if filter.HasSellingPrice != nil {
		if *filter.HasSellingPrice {
			where.And("items.selling_price IS NOT NULL")
		} else if !*filter.HasSellingPrice {
			where.And("items.selling_price IS NULL")
		}
	}

	queryWhere, err := bqb.New("?", where).ToRaw()
	if err != nil {
		fmt.Println(err)
	}

	return queryWhere
}

func addSort(filter Filter) string {
	sortKeys := "items.created_at"
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

func (r *repository) FindAll(filter Filter) (items []entities.Item, err error) {
	queryBuilder := getQueryBuilder()
	queryWhere := generateFilter(filter)
	formattedQuery := queryBuilder + queryWhere
	err = r.db.GetDB().Select(&items, formattedQuery)
	if err != nil {
		fmt.Println(err)
	}
	return items, err
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

func (r *repository) FindVariantByItemIDs(ids []int) (variants []entities.Variant, err error) {
	query := fmt.Sprintf(`
			select
				variants.id,
				variants.slug,
				variants.name,
				item_variant_categories.item_id,
				variants.variant_category_slug,
				variant_categories.name as variant_category_name,
				variants.branch_channel_id,
				variants.price,
				variants.in_stock,
				variants.created_at,
				variants.updated_at
			from item_variant_categories
			left join variant_categories on variant_categories.slug = item_variant_categories.variant_category_slug AND variant_categories.deleted_at IS NULL
			left join variants on variants.variant_category_slug = variant_categories.slug AND variants.variant_category_slug = item_variant_categories.variant_category_slug AND variants.deleted_at IS NULL
	`)

	where := bqb.Optional("WHERE")
	where.And(`item_variant_categories.item_id IN (?)`, ids)
	queryWhere, err := bqb.New("?", where).ToRaw()
	if err != nil {
		fmt.Println(err)
	}

	err = r.db.GetDB().Select(&variants, query+queryWhere)
	if err != nil {
		fmt.Println(err)
	}
	return variants, err
}
