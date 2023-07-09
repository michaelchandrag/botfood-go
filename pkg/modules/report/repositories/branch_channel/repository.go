package repositories

import (
	"database/sql"
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
	FindOne(filter Filter) (branchChannel entities.BranchChannel, err error)
	FindAll(filter Filter) (branchChannels []entities.BranchChannel, err error)
}

type repository struct {
	db database.MainDB
}

type Filter struct {
	ID        *int
	BrandID   *int
	BranchIDs []int
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
		JOIN branchs on branchs.id = branch_channels.branch_id AND branchs.deleted_at IS NULL
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
	sortKeys := "branch_channels.id"
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

func (r *repository) FindAll(filter Filter) (branchChannels []entities.BranchChannel, err error) {
	queryBuilder := getQueryBuilder()
	queryWhere := generateFilter(filter)
	formattedQuery := queryBuilder + queryWhere + addSort(filter)
	err = r.db.GetDB().Select(&branchChannels, formattedQuery)
	if err != nil {
		fmt.Println(err)
	}
	return branchChannels, err
}
