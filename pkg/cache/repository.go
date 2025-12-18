package cache

import (
	"errors"

	"zabbix-technical-task/pkg/userrecord"
)

var (
	errRecordExists   = errors.New("record already exists")
	errRecordNotFound = errors.New("record not found")
	errIDCannotChange = errors.New("cannot change record ID")
	errSaveRecords    = errors.New("failed to write records to file")
)

// Cache defines the interface for cache operations.
type Cache interface {
	Add(id uint64, record userrecord.Record) error
	Get(id uint64) (userrecord.Record, error)
	Update(id uint64, record userrecord.Record) error
	Delete(id uint64) error
	SaveRecords() error
}
