package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	logger "github.com/michaelchandrag/botfood-go/internal/logger"
	"github.com/michaelchandrag/botfood-go/utils"
)

type MainDB interface {
	GetDB() (db *sqlx.DB)
}

type mainDB struct {
	db *sqlx.DB
}

func ConnectMainDB() (db *sqlx.DB, err error) {
	logger.Agent.Info("Initialize Connection to Main DB")

	dbHost := utils.GetEnv("BOTFOOD_DB_HOST", "host")
	dbPort := utils.GetEnv("BOTFOOD_DB_PORT", "3306")
	dbName := utils.GetEnv("BOTFOOD_DB_NAME", "botfood")
	dbUser := utils.GetEnv("BOTFOOD_DB_USER", "user")
	dbPass := utils.GetEnv("BOTFOOD_DB_PASS", "pass")

	sHost := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err = sqlx.Connect("mysql", sHost)
	if err != nil {
		logger.Agent.Error(fmt.Sprintf("Error loading to Database %s", err.Error()))
		panic(err)
	}

	logger.Agent.Info(fmt.Sprintf("Success connect to database %s", dbName))
	return db, nil
}

func NewDB(db *sqlx.DB) MainDB {
	return &mainDB{
		db: db,
	}
}

func (d *mainDB) GetDB() (db *sqlx.DB) {
	return d.db
}
