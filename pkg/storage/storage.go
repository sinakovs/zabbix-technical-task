package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"zabbix-technical-task/pkg/userrecord"
)

var _ Storage = (*FileStorage)(nil)

// FileStorage implements StorageRepo interface for file-based storage.
type FileStorage struct {
	filename string
}

// NewFileStorage creates a new FileStorage instance with the given filename.
func NewFileStorage(filename string) *FileStorage {
	return &FileStorage{
		filename: filename,
	}
}

// InitFromReader initializes the storage by loading records from the provided reader.
func (f *FileStorage) InitFromReader(r io.Reader, records map[uint64]userrecord.Record) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var rec userrecord.Record

		err := json.Unmarshal(scanner.Bytes(), &rec)
		if err != nil {
			log.Printf("failed to unmarshal a record: %v", err)

			continue
		}

		err = rec.Validate()
		if err != nil {
			log.Printf("failed to validate a record: %v", err)

			continue
		}

		id, err := rec.ID()
		if err != nil {
			log.Printf("failed to get record ID: %v", err)

			continue
		}

		records[id] = rec
	}

	err := scanner.Err()
	if err != nil {
		return fmt.Errorf("scanning file %q: %w", f.filename, errScanFile)
	}

	return nil
}

// Init initializes the storage by loading records from the file.
func (f *FileStorage) Init(records map[uint64]userrecord.Record) error {
	file, err := os.Open(f.filename)
	if err != nil {
		return fmt.Errorf("opening file %q: %w", f.filename, errOpenFile)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("failed to close file %q: %v", f.filename, err)
		}
	}()

	err = f.InitFromReader(file, records)
	if err != nil {
		return fmt.Errorf("initializing from file %q: %w", f.filename, err)
	}

	return nil
}

// Save writes all records to the storage file.
func (f *FileStorage) Save(records map[uint64]userrecord.Record) error {
	file, err := os.Create(f.filename)
	if err != nil {
		return fmt.Errorf("creating file %q: %w", f.filename, errCreateFile)
	}

	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			log.Printf("failed to close file %q: %v", f.filename, closeErr)
		}
	}()

	err = saveToWriter(file, records)
	if err != nil {
		return fmt.Errorf("saving to file %q: %w", f.filename, err)
	}

	return nil
}

// saveToWriter writes all records to the provided writer.
func saveToWriter(w io.Writer, records map[uint64]userrecord.Record) error {
	for _, rec := range records {
		data, marshalErr := json.Marshal(rec)
		if marshalErr != nil {
			log.Println(marshalErr)

			continue
		}

		_, err := w.Write(data)
		if err != nil {
			return fmt.Errorf("writing to writer: %w", errWriteRecords)
		}

		_, err = w.Write([]byte("\n"))
		if err != nil {
			return fmt.Errorf("writing newline to writer: %w", errWriteRecords)
		}
	}

	return nil
}
