package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"vm_coding_challenge/config"
	"vm_coding_challenge/controllers"
	"vm_coding_challenge/db"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	configPath = "./config.json"
)

func main() {
	// Load the config
	err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Config error: ", err)
		return
	}

	// setup db
	err = db.Setup()
	if err != nil {
		fmt.Println("DB error: ", err)
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/request", controllers.Request).Methods("POST")
	router.HandleFunc("/stats/{by:customer}/{id:[0-9]+}", controllers.Statistics).Methods("GET")
	router.HandleFunc("/stats/{by:day}/{day}", controllers.Statistics).Methods("GET")

	srv := &http.Server{
		Handler: handlers.CombinedLoggingHandler(os.Stdout, router),
		Addr:    config.Conf.ServiceURL,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Starting server at ", config.Conf.ServiceURL)
	log.Fatal(srv.ListenAndServe())

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)
}
