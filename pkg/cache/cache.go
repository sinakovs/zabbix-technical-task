package cache

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"zabbix-technical-task/pkg/storage"
	"zabbix-technical-task/pkg/userrecord"
)

var (
	errRecordExists   = errors.New("record already exists")
	errRecordNotFound = errors.New("record not found")
	errIDCannotChange = errors.New("cannot change record ID")
	errSaveRecords    = errors.New("failed to write records to file")
)

// RecordCache provides thread-safe access to a cache of records.
type RecordCache struct {
	mu      sync.RWMutex
	records map[uint64]userrecord.Record
	counter uint8
}

// New creates a new RecordCache instance.
func New() *RecordCache {
	records := make(map[uint64]userrecord.Record)

	err := storage.InitStorage(records)
	if err != nil {
		log.Printf("error initializing storage: %v\n", err)
	}

	return &RecordCache{
		records: records,
		counter: 0,
	}
}

// Add adds a new record to the cache.
func (r *RecordCache) Add(id uint64, record userrecord.Record) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.counter >= 50 {
		log.Println("cache limit reached")

		err := storage.WriteRecordsToFile(r.records)
		if err != nil {
			return fmt.Errorf("writing records to file: %w", errSaveRecords)
		}

		r.counter = 0
	}

	_, exists := r.records[id]
	if exists {
		return fmt.Errorf("record with id %d: %w", id, errRecordExists)
	}

	r.records[id] = record

	return nil
}

// Get retrieves a record by ID from the cache.
func (r *RecordCache) Get(id uint64) (userrecord.Record, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	record, exists := r.records[id]
	if !exists {
		return nil, fmt.Errorf("record with id %d: %w", id, errRecordNotFound)
	}

	return record, nil
}

// Update updates an existing record in the cache.
func (r *RecordCache) Update(id uint64, record userrecord.Record) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.records[id]
	if !exists {
		return fmt.Errorf("record with id %d: %w", id, errRecordNotFound)
	}

	baseID, err := record.ID()
	if err != nil {
		return fmt.Errorf("getting record ID: %w", err)
	}

	if id != baseID {
		return fmt.Errorf("cannot change record's id from %d to %d: %w", id, baseID, errIDCannotChange)
	}

	r.records[id] = record

	return nil
}

// Delete removes a record by ID from the cache.
func (r *RecordCache) Delete(id uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.records[id]
	if !exists {
		return fmt.Errorf("record with id %d: %w", id, errRecordNotFound)
	}

	delete(r.records, id)

	return nil
}

// SaveRecords saves all records from the cache to persistent storage.
func (r *RecordCache) SaveRecords() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	err := storage.WriteRecordsToFile(r.records)
	if err != nil {
		return fmt.Errorf("saving records to file: %w", errSaveRecords)
	}

	return nil
}
