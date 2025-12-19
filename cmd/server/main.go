package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"zabbix-technical-task/internal/router"
	"zabbix-technical-task/pkg/cache"
	"zabbix-technical-task/pkg/storage"
)

func main() {
	port := "8080"

	fileStorage := storage.NewFileStorage("data/data.txt")

	records := cache.New(fileStorage)
	if records == nil {
		log.Fatal("failed to create record cache")

		return
	}

	routes := router.New(records)

	srv := &http.Server{
		Addr:              "" + ":" + port,
		Handler:           routes.Mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Println("Listening on :" + port)

		err := srv.ListenAndServe()
		if err != nil {
			log.Printf("Stopped listening: %v\n", err)
		}
	}()

	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	<-shutdown.Done()

	log.Println("Shutting down server...")

	err := srv.Shutdown(context.Background())
	if err != nil {
		log.Printf("Shutdown with error: %v\n", err)
	}

	err = records.SaveRecords()
	if err != nil {
		log.Printf("Shutdown whit error saving records: %v\n", err)
	}

	log.Println("Shutdown complete.")
}
