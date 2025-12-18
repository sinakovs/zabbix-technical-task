package storage

import (
	"errors"

	"zabbix-technical-task/pkg/userrecord"
)

var (
	errOpenFile     = errors.New("failed to open file")
	errScanFile     = errors.New("scanner error")
	errCreateFile   = errors.New("failed to create file")
	errWriteRecords = errors.New("failed to write records to file")
)

// Storage defines the interface for storage operations.
type Storage interface {
	Init(records map[uint64]userrecord.Record) error
	Save(records map[uint64]userrecord.Record) error
}
