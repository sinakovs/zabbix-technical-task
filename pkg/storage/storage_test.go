package storage

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"zabbix-technical-task/pkg/userrecord"
)

func TestNewFileStorage(t *testing.T) {
	t.Parallel()

	filename := "testfile.db"
	storage := NewFileStorage(filename)

	if storage.filename != filename {
		t.Errorf("expected filename %q, got %q", filename, storage.filename)
	}
}

func TestInitFromReader(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	log.SetOutput(&buf)

	defer log.SetOutput(os.Stderr)

	data := []string{
		`{"id":1.0,"Name":"Alice","Age":30}`,
		`invalid-json`, `{"id":"1.0","Name":"Bob","Age":25}`,
		`{"Name":"InvalidRecord","Age":"NaN"}`,
	}

	reader := strings.NewReader(data[0])
	storage := &FileStorage{}
	records := make(map[uint64]userrecord.Record)

	err := storage.InitFromReader(reader, records)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reader = strings.NewReader(data[1])

	err = storage.InitFromReader(reader, records)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logContent := buf.String()

	if !strings.Contains(logContent, "failed to unmarshal a record") {
		t.Errorf("expected unmarshal log, got: %s", logContent)
	}

	reader = strings.NewReader(data[2])

	err = storage.InitFromReader(reader, records)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logContent = buf.String()

	if !strings.Contains(logContent, "failed to validate a record") {
		t.Errorf("expected validate log, got: %s", logContent)
	}
}

func TestSaveToWriter(t *testing.T) {
	t.Parallel()

	records := map[uint64]userrecord.Record{
		1: {"id": 1.0, "Name": "Alice", "Age": 30},
		2: {"id": 2.0, "Name": "Bob", "Age": 25},
	}

	var buf bytes.Buffer

	err := saveToWriter(&buf, records)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := buf.String()
	if !strings.Contains(content, `"Name":"Alice"`) || !strings.Contains(content, `"Name":"Bob"`) {
		t.Errorf("unexpected content: %s", content)
	}

	records = map[uint64]userrecord.Record{}
	buf = bytes.Buffer{}

	err = saveToWriter(&buf, records)
	if err != nil {
		t.Fatalf("unexpected error for empty records: %v", err)
	}

	content = buf.String()
	if content != "" {
		t.Errorf("expected empty content for empty records, got: %s", content)
	}
}
