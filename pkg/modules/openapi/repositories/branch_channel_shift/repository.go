package repositories

import (
	"fmt"

	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/pkg/modules/openapi/entities"
	bqb "github.com/nullism/bqb"
)

type Repository interface {
	FindAll(filter Filter) (shifts []entities.BranchChannelShift, err error)
	FindAllGrouped(filter Filter) (grouped entities.GroupedBranchChannelShift, err error)
}

type repository struct {
	db database.MainDB
}

type Filter struct {
	BranchChannelID *int
}

func NewRepository(db database.MainDB) Repository {
	return &repository{
		db: db,
	}
}

func getQueryBuilder() string {
	query := fmt.Sprintf(`
		SELECT
			branch_channel_shifts.id,
			branch_channel_shifts.branch_channel_id,
			branch_channel_shifts.day,
			branch_channel_shifts.open_time,
			branch_channel_shifts.close_time
		FROM
			branch_channel_shifts
	`)
	return query
}

func generateFilter(filter Filter) string {

	where := bqb.Optional("WHERE")

	where.And("branch_channel_shifts.deleted_at IS NULL")

	if filter.BranchChannelID != nil {
		where.And("branch_channel_shifts.branch_channel_id = ?", *filter.BranchChannelID)
	}
	queryWhere, err := bqb.New("?", where).ToRaw()
	if err != nil {
		fmt.Println(err)
	}

	return queryWhere
}

func (r *repository) FindAll(filter Filter) (shifts []entities.BranchChannelShift, err error) {
	queryBuilder := getQueryBuilder()
	queryWhere := generateFilter(filter)
	formattedQuery := queryBuilder + queryWhere
	err = r.db.GetDB().Select(&shifts, formattedQuery)
	if err != nil {
		fmt.Println(err)
	}
	return shifts, err
}

func (r *repository) FindAllGrouped(filter Filter) (grouped entities.GroupedBranchChannelShift, err error) {
	var shifts []entities.BranchChannelShift
	queryBuilder := getQueryBuilder()
	queryWhere := generateFilter(filter)
	formattedQuery := queryBuilder + queryWhere + " ORDER BY open_time ASC"
	err = r.db.GetDB().Select(&shifts, formattedQuery)
	if err != nil {
		fmt.Println(err)
	}

	for _, val := range shifts {
		if val.Day == 1 {
			grouped.Monday = append(grouped.Monday, val.ToModern())
		} else if val.Day == 2 {
			grouped.Tuesday = append(grouped.Tuesday, val.ToModern())
		} else if val.Day == 3 {
			grouped.Wednesday = append(grouped.Wednesday, val.ToModern())
		} else if val.Day == 4 {
			grouped.Thursday = append(grouped.Thursday, val.ToModern())
		} else if val.Day == 5 {
			grouped.Friday = append(grouped.Friday, val.ToModern())
		} else if val.Day == 6 {
			grouped.Saturday = append(grouped.Saturday, val.ToModern())
		} else if val.Day == 7 {
			grouped.Sunday = append(grouped.Sunday, val.ToModern())
		}
	}

	return grouped, err
}
