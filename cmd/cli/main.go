package main

import (
	"os"
	"fmt"

	"github.com/michaelchandrag/internal/driver"
)

func main () {
	fmt.Println("Entrypoint CLI Application")
	dbName := os.Getenv("DB_NAME")
	dbUsername := os.GETENV("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connection, err := driver.Connect(dbHost, dbPort, dbUsername, dbPassword, dbName)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}