package main

import (
	"fmt"
	"net/http"
	"zabbix-technical-task/internal/router"
	"zabbix-technical-task/pkg/record"
)

var addr string = "localhost:8080"

func main() {

	records := record.NewRecordStore()

	router := router.New(records)

	err := http.ListenAndServe(addr, router.Mux)
	if err != nil {
		fmt.Errorf("failed to start server: %w", err)
	}
}
