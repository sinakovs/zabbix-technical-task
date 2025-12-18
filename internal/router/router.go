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
func New(records cache.Cache) Routes {
	mux := http.NewServeMux()

	recordHandler := handler.New(records)

	mux.HandleFunc("POST /records", recordHandler.Post)
	mux.HandleFunc("GET /records/", recordHandler.Get)
	mux.HandleFunc("PUT /records/", recordHandler.Put)
	mux.HandleFunc("DELETE /records/", recordHandler.Delete)

	return Routes{
		Mux: mux,
	}
}
