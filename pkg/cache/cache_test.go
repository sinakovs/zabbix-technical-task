package cache

import (
	"errors"
	"log"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
	"zabbix-technical-task/pkg/storage/mocks"
	"zabbix-technical-task/pkg/userrecord"
)

func TestNew(t *testing.T) {
	t.Parallel()

	mockStorage := new(mocks.Storage)

	mockStorage.On("Init", mock.Anything).Return(nil).Once()

	cache := New(mockStorage)

	if cache == nil {
		t.Fatal("expected non-nil cache")
	}

	if len(cache.records) != 0 {
		t.Fatalf("expected empty records map, got %d records", len(cache.records))
	}

	mockStorage.On("Init", mock.Anything).Return(errors.New("Some Init error")).Once()

	cache = New(mockStorage)

	if cache != nil {
		t.Fatal("expected non-nil cache even on Init error")
	}

	mockStorage.AssertCalled(t, "Init", mock.Anything)
}

func TestAdd(t *testing.T) {
	t.Parallel()

	mockStorage := new(mocks.Storage)

	mockStorage.On("Init", mock.Anything).Return(nil)
	mockStorage.On("Save", mock.Anything).Return(nil).Once()

	cache := New(mockStorage)
	cache.counter = 48 // Set counter close to limit for testing

	record := userrecord.Record{
		"id":    1,
		"Name":  "John Doe",
		"Email": "fsfs",
		"Age":   30,
	}

	err := cache.Add(1, record)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = cache.Add(1, record)
	if err == nil {
		t.Fatalf("expected error for duplicate record, got nil")
	}

	err = cache.Add(2, record)
	if err != nil {
		t.Fatalf("expected no error when adding record after save, got %v", err)
	}

	mockStorage.On("Save", mock.Anything).Return(errors.New("disk full")).Once()

	cache.counter = 49

	err = cache.Add(3, record)
	if err == nil {
		t.Fatalf("expected error due to storage save failure, got nil")
	}

	mockStorage.AssertCalled(t, "Init", mock.Anything)
	mockStorage.AssertNumberOfCalls(t, "Save", 2)
}

func TestGet(t *testing.T) {
	t.Parallel()

	mockStorage := new(mocks.Storage)

	mockStorage.On("Init", mock.Anything).Return(nil)

	cache := New(mockStorage)

	record := userrecord.Record{
		"id":    1,
		"Name":  "John Doe",
		"Email": "fsfs",
		"Age":   30,
	}

	err := cache.Add(1, record)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, err := cache.Get(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got["Name"] != "John Doe" {
		t.Errorf("expected Name to be 'John Doe', got %v", got["Name"])
	}

	_, err = cache.Get(2)
	if err == nil {
		t.Fatalf("expected error for non-existent record, got nil")
	}

	mockStorage.AssertCalled(t, "Init", mock.Anything)
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	mockStorage := new(mocks.Storage)

	mockStorage.On("Init", mock.Anything).Return(nil)

	cache := New(mockStorage)

	record := userrecord.Record{
		"id":    uint64(1),
		"Name":  "John Doe",
		"Email": "fsfs",
		"Age":   30,
	}

	err := cache.Add(1, record)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedRecord := userrecord.Record{
		"id":    uint64(1),
		"Name":  "James Doe",
		"Phone": "123-456-7890",
		"Age":   25,
	}

	log.Println(updatedRecord.ID())

	err = cache.Update(1, updatedRecord)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, _ := cache.Get(1)

	if got["Name"] != "James Doe" {
		t.Errorf("expected Name to be 'James Doe', got %v", got["Name"])
	}

	if got["Phone"] != "123-456-7890" {
		t.Errorf("expected Phone to be '123-456-7890', got %v", got["Phone"])
	}

	if got["Age"] != 25 {
		t.Errorf("expected Age to be 25, got %v", got["Age"])
	}

	err = cache.Update(2, updatedRecord)
	if err == nil {
		t.Fatalf("expected error for non-existent record, got nil")
	}

	err = cache.Update(1, userrecord.Record{
		"id": uint64(2),
	})
	if err == nil {
		t.Fatalf("expected error for changing record ID, got nil")
	}

	err = cache.Update(1, userrecord.Record{
		"Name":  "James Doe",
		"Phone": "123-456-7890",
		"Age":   25,
	})
	if err == nil {
		t.Fatalf("expected error for missing record ID, got nil")
	}

	mockStorage.AssertCalled(t, "Init", mock.Anything)
}

func TestDelete(t *testing.T) {
	t.Parallel()

	mockStorage := new(mocks.Storage)

	mockStorage.On("Init", mock.Anything).Return(nil)

	cache := New(mockStorage)

	cache.records[1] = userrecord.Record{"id": 1}

	err := cache.Delete(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = cache.Get(1)
	if err == nil {
		t.Fatalf("expected error for deleted record, got nil")
	}

	err = cache.Delete(2)
	if err == nil {
		t.Fatalf("expected error for non-existent record, got nil")
	}

	mockStorage.AssertCalled(t, "Init", mock.Anything)
}

func TestSaveRecords(t *testing.T) {
	t.Parallel()

	mockStorage := new(mocks.Storage)

	mockStorage.On("Init", mock.Anything).Return(nil)
	mockStorage.On("Save", mock.Anything).Return(errors.New("disk full")).Once()

	cache := New(mockStorage)

	err := cache.SaveRecords()
	if err == nil {
		t.Fatalf("expected error from SaveRecords, got nil")
	}

	mockStorage.On("Save", mock.Anything).Return(nil).Once()

	err = cache.SaveRecords()
	if err != nil {
		t.Fatalf("expected no error from SaveRecords, got %v", err)
	}

	mockStorage.AssertNumberOfCalls(t, "Save", 2)
}

func TestRecordCache_Race(t *testing.T) {
	t.Parallel()

	mockStorage := new(mocks.Storage)
	mockStorage.On("Init", mock.Anything).Return(nil)
	mockStorage.On("Save", mock.Anything).Return(nil)

	cache := New(mockStorage)

	var wg sync.WaitGroup

	for i := range uint64(100) {
		wg.Add(1)

		id := i

		go func(id uint64) {
			defer wg.Done()

			rec := userrecord.Record{
				"id":    id,
				"Name":  "John Doe",
				"Email": "fsfs",
				"Age":   30,
			}

			_ = cache.Add(id, rec)
			_, _ = cache.Get(id)
		}(id)
	}

	wg.Wait()
}
