package repositories

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/entities"
	bqb "github.com/nullism/bqb"
	"golang.org/x/exp/slices"
)

type Repository interface {
	FindPaginated(filter Filter) (result PaginatedData, err error)
	FindAll(filter Filter) (variants []entities.Variant, err error)
	FindByItemID(itemID int) (vcs []entities.VariantCategory, err error)
	FindByBranchChannelID(branchChannelID int) (variants []entities.Variant, err error)
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

	Name                 string
	Keyword              string
	BranchChannelChannel string
	BranchChannelName    string
	InStock              *int

	VariantCategoryID   *int
	VariantCategoryName string

	Page *int
	Data *int
}

type PaginatedData struct {
	Data        []entities.Variant `json:"data"`
	CurrentPage int                `json:"current_page"`
	LimitData   int                `json:"limit_data"`
	TotalPage   int                `json:"total_page"`
	TotalData   int                `json:"total_data"`
}

func NewRepository(db database.MainDB) Repository {
	return &repository{
		db: db,
	}
}

func getQueryBuilder() string {
	query := fmt.Sprintf(`
		SELECT
			variants.id,
			variants.branch_channel_id,
			branch_channels.name as branch_channel_name,
			branch_channels.channel as branch_channel_channel,
			variant_categories.id as variant_category_id,
			variant_categories.name as variant_category_name,
			variants.name,
			variants.price,
			variants.in_stock
		FROM
			variants
		JOIN branch_channels ON branch_channels.id = variants.branch_channel_id AND branch_channels.deleted_at IS NULL
		JOIN variant_categories ON variant_categories.slug = variants.variant_category_slug AND variant_categories.branch_channel_id = variants.branch_channel_id AND variant_categories.deleted_at IS NULL
	`)
	return query
}

func getTotalQueryBuilder() string {
	query := fmt.Sprintf(`
		SELECT
			COUNT(*) as total
		FROM
			variants
		JOIN branch_channels ON branch_channels.id = variants.branch_channel_id AND branch_channels.deleted_at IS NULL
	`)
	return query
}

func generateFilter(filter Filter) string {

	where := bqb.Optional("WHERE")

	where.And("variants.deleted_at IS NULL")

	if filter.ID != nil {
		where.And("variants.id = ?", *filter.ID)
	}

	if filter.BrandID != nil {
		where.And("branch_channels.brand_id = ?", *filter.BrandID)
	}

	if filter.BranchChannelID != nil {
		where.And("variants.branch_channel_id = ?", *filter.BranchChannelID)
	}

	if filter.BranchChannelChannel != "" {
		where.And("branch_channels.channel = ?", filter.BranchChannelChannel)
	}

	if filter.Name != "" {
		name := "%" + filter.Name + "%"
		where.And("variants.name LIKE ?", name)
	}

	if filter.BranchChannelName != "" {
		name := "%" + filter.BranchChannelName + "%"
		where.And("branch_channels.name LIKE ?", name)
	}

	if filter.InStock != nil {
		where.And("variants.in_stock = ?", *filter.InStock)
	}

	if filter.VariantCategoryID != nil {
		where.And("variant_categories.id = ?", *filter.VariantCategoryID)
	}

	if filter.VariantCategoryName != "" {
		name := "%" + filter.VariantCategoryName + "%"
		where.And("variant_categories.name LIKE ?", name)
	}

	queryWhere, err := bqb.New("?", where).ToRaw()
	if err != nil {
		fmt.Println(err)
	}

	return queryWhere
}

func addSort(filter Filter) string {
	sortKeys := "variants.created_at"
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

func (r *repository) FindAll(filter Filter) (variants []entities.Variant, err error) {
	queryBuilder := getQueryBuilder()
	queryWhere := generateFilter(filter)
	formattedQuery := queryBuilder + queryWhere
	err = r.db.GetDB().Select(&variants, formattedQuery)
	if err != nil {
		fmt.Println(err)
	}
	return variants, err
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

func (r *repository) FindByItemID(itemID int) (vcs []entities.VariantCategory, err error) {
	formattedQuery := fmt.Sprintf(`SELECT
			variants.id,
			variants.name,
			variants.in_stock,
			variants.price,
			variant_categories.id as variant_category_id,
			variant_categories.name as variant_category_name,
			variant_categories.is_required as variant_category_is_required,
			variant_categories.min_quantity as variant_category_min_quantity,
			variant_categories.max_quantity as variant_category_max_quantity
		FROM
			variants
		JOIN item_variant_categories ON variants.variant_category_slug = item_variant_categories.variant_category_slug AND item_variant_categories.item_id = %d AND item_variant_categories.deleted_at IS NULL
		JOIN variant_categories ON variant_categories.slug = variants.variant_category_slug AND variant_categories.slug = item_variant_categories.variant_category_slug AND variant_categories.deleted_at IS NULL
		WHERE variants.deleted_at IS NULL
	`, itemID)
	var variants []entities.Variant
	err = r.db.GetDB().Select(&variants, formattedQuery)
	if err != nil {
		fmt.Println(err)
	}

	for _, val := range variants {
		idxVc := slices.IndexFunc(vcs, func(vc entities.VariantCategory) bool { return vc.ID == val.VariantCategoryID })
		if idxVc == -1 {
			// not found
			isRequired := false
			if val.VariantCategoryIsRequired == 1 {
				isRequired = true
			} else {
				isRequired = false
			}
			var mvs []entities.ModernVariant
			mvs = append(mvs, val.ToModern())
			vc := entities.VariantCategory{
				ID:          val.VariantCategoryID,
				Name:        val.VariantCategoryName,
				IsRequired:  isRequired,
				MinQuantity: val.VariantCategoryMinQuantity,
				MaxQuantity: val.VariantCategoryMaxQuantity,
				Variants:    mvs,
			}
			vcs = append(vcs, vc)
		} else {
			// exists
			vcs[idxVc].Variants = append(vcs[idxVc].Variants, val.ToModern())
		}
	}

	return vcs, err
}

func (r *repository) FindByBranchChannelID(branchChannelID int) (variants []entities.Variant, err error) {
	formattedQuery := fmt.Sprintf(`SELECT
			variants.id,
			variants.name,
			variants.in_stock,
			variants.price,
			item_variant_categories.id as item_variant_category_id,
			item_variant_categories.item_id as item_variant_category_item_id,
			variant_categories.id as variant_category_id,
			variant_categories.name as variant_category_name,
			variant_categories.is_required as variant_category_is_required,
			variant_categories.min_quantity as variant_category_min_quantity,
			variant_categories.max_quantity as variant_category_max_quantity
		FROM
			variants
		JOIN item_variant_categories ON variants.variant_category_slug = item_variant_categories.variant_category_slug AND item_variant_categories.deleted_at IS NULL
		JOIN items on items.id = item_variant_categories.item_id AND items.deleted_at IS NULL
		JOIN variant_categories ON variant_categories.slug = variants.variant_category_slug AND variant_categories.slug = item_variant_categories.variant_category_slug AND variant_categories.deleted_at IS NULL
		WHERE variants.deleted_at IS NULL AND variants.branch_channel_id = %d
	`, branchChannelID)
	err = r.db.GetDB().Select(&variants, formattedQuery)
	if err != nil {
		fmt.Println(err)
	}

	return variants, err
}
