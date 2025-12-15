package router

import (
	"net/http"
	"zabbix-technical-task/internal/handler"
	"zabbix-technical-task/pkg/record"
)

type Router struct {
	Mux *http.ServeMux
}

func New(records *record.RecordStore) Router {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /records", handler.PostData)
	mux.HandleFunc("GET /records/", handler.GetData)
	mux.HandleFunc("PUT /records/", handler.PutData)
	mux.HandleFunc("DELETE /records/", handler.DeleteData)

	return Router{
		Mux: mux,
	}
}
