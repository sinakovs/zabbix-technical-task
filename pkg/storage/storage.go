package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"zabbix-technical-task/pkg/userrecord"
)

const filename = "data/data.txt"

// InitStorage initializes the storage by loading records from the file.
func InitStorage(records map[uint64]userrecord.Record) error {
	file, err := os.Open(filename)
	if err != nil {
		errOpenFile := fmt.Errorf("failed to open file %q: %w", filename, err)

		return errOpenFile
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("failed to close file %q: %v", filename, err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var rec userrecord.Record

		err = json.Unmarshal(scanner.Bytes(), &rec)
		if err != nil {
			log.Printf("failed to unmarshal a record: %v", err)

			continue
		}

		err = rec.Validate()
		if err != nil {
			log.Printf("failed to validate a record: %v", err)

			continue
		}

		id, errID := rec.ID()
		if errID != nil {
			log.Printf("failed to get record ID: %v", errID)

			continue
		}

		records[id] = rec
	}

	errScanFile := fmt.Errorf("scanner error for file %q: %w", filename, err)

	return errScanFile
}

// WriteRecordsToFile writes all records to the storage file.
func WriteRecordsToFile(records map[uint64]userrecord.Record) error {
	file, err := os.Create(filename)
	if err != nil {
		errCreateFile := fmt.Errorf("failed to write to file %q: %w", filename, err)

		return errCreateFile
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("failed to close file %q: %v", filename, err)
		}
	}()

	lines := make([]byte, 0, len(records))

	for _, rec := range records {
		data, marshalErr := json.Marshal(rec)
		if marshalErr != nil {
			log.Println(marshalErr)

			continue
		}

		lines = append(lines, data...)

		lines = append(lines, '\n')
	}

	_, err = file.Write(lines)
	if err != nil {
		errWriteRecords := fmt.Errorf("failed to write to file %q: %w", filename, err)

		return errWriteRecords
	}

	return nil
}
