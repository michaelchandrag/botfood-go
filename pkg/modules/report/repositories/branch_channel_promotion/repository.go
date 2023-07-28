package repositories

import (
	"fmt"
	"strings"

	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/pkg/modules/report/entities"
	bqb "github.com/nullism/bqb"
)

type Repository interface {
	FindAll(filter Filter) (promotions []entities.BranchChannelPromotion, err error)
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

	Page *int
	Data *int
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
			branch_channel_promotions.branch_channel_id,
			branch_channels.name as branch_channel_name,
			branch_channels.channel as branch_channel_channel,
			branch_channel_promotions.slug,
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
			COUNT(*) as total
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
		where.And("branch_channel_promotions.branch_channel_id = ?", *filter.BranchChannelID)
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

func (r *repository) FindAll(filter Filter) (promotions []entities.BranchChannelPromotion, err error) {
	queryBuilder := getQueryBuilder()
	queryWhere := generateFilter(filter)
	formattedQuery := queryBuilder + queryWhere
	err = r.db.GetDB().Select(&promotions, formattedQuery)
	if err != nil {
		fmt.Println(err)
	}
	return promotions, err
}
