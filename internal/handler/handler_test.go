package handler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"zabbix-technical-task/pkg/cache/mocks"
	"zabbix-technical-task/pkg/userrecord"
)

func TestNew(t *testing.T) {
	t.Parallel()

	mockCache := new(mocks.Cache)

	handler := New(mockCache)

	if handler == nil {
		t.Fatal("expected non-nil handler")
	}

	if handler.cache != mockCache {
		t.Fatal("expected handler to use the provided cache")
	}
}

func TestPost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		payload        string
		cacheErr       error
		expectedStatus int
		expectedBody   string
	}{
		{"valid record", `{"id":1.0,"Name":"Alice","Age":30}`, nil, http.StatusCreated, "Record created\n"},
		{"invalid JSON", `invalid-json`, nil, http.StatusBadRequest, "Invalid request payload\n"},
		{"validation fail", `{"Name":"","Age":-1}`, nil, http.StatusBadRequest, "missing id\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cache := new(mocks.Cache)
			handler := New(cache)

			cache.On("Add", mock.Anything, mock.Anything).Return(nil)

			req := httptest.NewRequest(http.MethodPost, "/records", strings.NewReader(tt.payload))
			w := httptest.NewRecorder()

			handler.Post(w, req)

			res := w.Result()
			body, _ := io.ReadAll(res.Body)

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			if string(body) != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, string(body))
			}
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		recordID       string
		cacheRecord    userrecord.Record
		cacheErr       error
		expectedStatus int
		expectedBody   string
	}{
		{
			"valid record",
			"1",
			userrecord.Record{"Name": "Alice", "Age": 30},
			nil,
			http.StatusOK,
			`{"Age":30,"Name":"Alice"}` + "\n",
		},
		{
			"nonexistent record",
			"2",
			userrecord.Record{},
			errors.New("not found"),
			http.StatusNotFound,
			"not found\n",
		},
		{
			"invalid ID",
			"abc",
			userrecord.Record{},
			nil,
			http.StatusBadRequest,
			"failed convert str to int: strconv.Atoi: parsing \"abc\": invalid syntax\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cache := new(mocks.Cache)
			handler := New(cache)

			cache.On("Get", mock.Anything).Return(tt.cacheRecord, tt.cacheErr)

			req := httptest.NewRequest(http.MethodGet, "/records/"+tt.recordID, nil)
			w := httptest.NewRecorder()

			handler.Get(w, req)

			res := w.Result()
			body, _ := io.ReadAll(res.Body)

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			if string(body) != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, string(body))
			}
		})
	}
}

func TestPut(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		recordID       string
		payload        string
		cacheErr       error
		expectedStatus int
		expectedBody   string
	}{
		{
			"valid update",
			"1",
			`{"id":1.0,"Name":"Alice","Age":31}`,
			nil,
			http.StatusOK,
			"Record updated\n",
		},
		{
			"invalid JSON",
			"1", `invalid-json`,
			nil,
			http.StatusBadRequest,
			"Invalid request payload\n",
		},
		{
			"validation fail",
			"1",
			`{"Name":"","Age":-1}`,
			nil,
			http.StatusBadRequest,
			"missing id\n",
		},
		{
			"nonexistent record",
			"2",
			`{"id":1.0,"Name":"Bob","Age":25}`,
			errors.New("not found"),
			http.StatusNotFound,
			"not found\n",
		},
		{
			"invalid ID",
			"abc",
			`{"Name":"Bob","Age":25}`,
			nil,
			http.StatusBadRequest,
			"failed convert str to int: strconv.Atoi: parsing \"abc\": invalid syntax\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cache := new(mocks.Cache)
			handler := New(cache)

			cache.On("Update", mock.Anything, mock.Anything).Return(tt.cacheErr)

			req := httptest.NewRequest(http.MethodPut, "/records/"+tt.recordID, strings.NewReader(tt.payload))
			w := httptest.NewRecorder()

			handler.Put(w, req)

			res := w.Result()
			body, _ := io.ReadAll(res.Body)

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			if string(body) != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, string(body))
			}
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		recordID       string
		cacheErr       error
		expectedStatus int
		expectedBody   string
	}{
		{
			"valid delete",
			"1",
			nil,
			http.StatusNoContent,
			"",
		},
		{
			"nonexistent record",
			"2",
			errors.New("not found"),
			http.StatusNotFound,
			"not found\n",
		},
		{
			"invalid ID",
			"abc",
			nil,
			http.StatusBadRequest,
			"failed convert str to int: strconv.Atoi: parsing \"abc\": invalid syntax\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cache := new(mocks.Cache)
			handler := New(cache)

			cache.On("Delete", mock.Anything).Return(tt.cacheErr)

			req := httptest.NewRequest(http.MethodDelete, "/records/"+tt.recordID, nil)
			w := httptest.NewRecorder()

			handler.Delete(w, req)

			res := w.Result()
			body, _ := io.ReadAll(res.Body)

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			if string(body) != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, string(body))
			}
		})
	}
}
