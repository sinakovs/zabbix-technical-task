package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"zabbix-technical-task/pkg/cache"
	"zabbix-technical-task/pkg/userrecord"
)

var errWrongID = errors.New("wrong ID")

// RecordHandler handles HTTP requests for record operations.
type RecordHandler struct {
	records *cache.RecordCache
}

// New creates a new handler with the given record cache.
func New(records *cache.RecordCache) *RecordHandler {
	return &RecordHandler{
		records: records,
	}
}

// PostData handles POST /records requests to create a new record.
func (h *RecordHandler) PostData(w http.ResponseWriter, r *http.Request) {
	var record userrecord.Record

	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)

		return
	}

	err = record.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	id, err := record.ID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	err = h.records.Add(id, record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)

		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte("Record created\n"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

// GetData handles GET /records/{id} requests to retrieve a record by ID.
func (h *RecordHandler) GetData(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(strings.TrimPrefix(r.URL.Path, "/records/"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	record, err := h.records.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

// PutData handles PUT /records/{id} requests to update an existing record.
func (h *RecordHandler) PutData(w http.ResponseWriter, r *http.Request) {
	var record userrecord.Record

	id, err := parseID(strings.TrimPrefix(r.URL.Path, "/records/"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	err = json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)

		return
	}

	err = record.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	err = h.records.Update(id, record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte("Record updated\n"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

// DeleteData andles DELETE /records/{id} requests to delete a record by ID.
func (h *RecordHandler) DeleteData(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(strings.TrimPrefix(r.URL.Path, "/records/"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	err = h.records.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(s string) (uint64, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed convert str to int: %w", err)
	}

	if id < 0 {
		return 0, fmt.Errorf("invalid request: %w", errWrongID)
	}

	return uint64(id), nil
}
