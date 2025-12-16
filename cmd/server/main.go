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
)

func main() {
	addr := "localhost"
	port := "8080"

	records := cache.New()

	routes := router.New(records)

	srv := &http.Server{
		Addr:              addr + ":" + port,
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

	err := records.SaveRecords()
	if err != nil {
		log.Printf("Shutdown whit error saving records: %v\n", err)
	}

	log.Println("Shutting down server...")

	err = srv.Shutdown(context.Background())
	if err != nil {
		log.Printf("Shutdown with error: %v\n", err)
	}

	log.Println("Shutdown complete.")
}
