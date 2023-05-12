package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"syscall"

	logger "github.com/michaelchandrag/botfood-go/internal/logger"
	http "github.com/michaelchandrag/botfood-go/pkg/http"
	utils "github.com/michaelchandrag/botfood-go/utils"

	database "github.com/michaelchandrag/botfood-go/infrastructures/database"
)

func main() {
	var wg sync.WaitGroup
	var errChan = make(chan error, 1)

	// Init Logger
	logger.InitLogger()
	logger.Agent.Info("Running Application BotFood")

	// Init Main DB
	db, err := database.ConnectMainDB()
	if err != nil {
		logger.Agent.Info(err.Error())
	}

	// Init Application
	wg.Add(1)
	go func() {
		defer wg.Done()
		port := utils.GetEnv("BOTFOOD_APP_PORT", "8080")
		fmt.Println("Starting listen address: ", port)
		logger.Agent.Info(fmt.Sprintf("Starting listen address: %s", port))
		errChan <- http.ServeHTTP(port, db)
	}()
	wg.Wait()

	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalChan:
		fmt.Println("Got an interrupt, exiting...")
	case err := <-errChan:
		if err != nil {
			fmt.Println("Error while running api, exiting...: ", err)
		}
	}
}
