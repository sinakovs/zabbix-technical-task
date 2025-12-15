package handler

import (
	"encoding/json"
	"net/http"
	"zabbix-technical-task/pkg/record"
)

func PostData(w http.ResponseWriter, r *http.Request) {
	var record record.Record
	record = make(map[string]interface{})

	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetData(w http.ResponseWriter, r *http.Request) {
	// Implementation for getting data
}

func PutData(w http.ResponseWriter, r *http.Request) {
	// Implementation for putting data
}

func DeleteData(w http.ResponseWriter, r *http.Request) {
	// Implementation for deleting data
}
