package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"syscall"
	http "github.com/michaelchandrag/botfood-go/pkg/http"
	utils "github.com/michaelchandrag/botfood-go/utils"
)

func main () {
	var wg sync.WaitGroup
	var errChan = make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		port := utils.GetEnv("BOTFOOD_APP_PORT", "8080")
		fmt.Println("Starting listen address: ", port)
		errChan <- http.ServeHTTP(port)
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