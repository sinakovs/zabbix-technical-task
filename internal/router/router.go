package router

import (
	"net/http"

	"zabbix-technical-task/internal/handler"
	"zabbix-technical-task/pkg/cache"
)

// Routes holds the HTTP request multiplexer.
type Routes struct {
	Mux *http.ServeMux
}

// New creates a new Routes instance with the given record cache.
func New(records *cache.RecordCache) Routes {
	mux := http.NewServeMux()

	recordHandler := handler.New(records)

	mux.HandleFunc("POST /records", recordHandler.PostData)
	mux.HandleFunc("GET /records/", recordHandler.GetData)
	mux.HandleFunc("PUT /records/", recordHandler.PutData)
	mux.HandleFunc("DELETE /records/", recordHandler.DeleteData)

	return Routes{
		Mux: mux,
	}
}
