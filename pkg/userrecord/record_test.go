package userrecord

import (
	"errors"
	"testing"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		record Record
		err    error
	}{
		{
			name:   "valid id",
			record: Record{"id": 123.0},
			err:    nil,
		},
		{
			name: "valid fields",
			record: Record{
				"id":    456.0,
				"name":  "test",
				"email": "",
				"age":   30,
			},
			err: nil,
		},
		{
			name:   "missing id",
			record: Record{"name": "test"},
			err:    errNoID,
		},
		{
			name:   "id not a number",
			record: Record{"id": "abc"},
			err:    errIDNotNumber,
		},
		{
			name:   "id negative",
			record: Record{"id": -5.0},
			err:    errIDNotUint,
		},
		{
			name:   "id not integer",
			record: Record{"id": 12.34},
			err:    errIDNotUint,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			err := c.record.Validate()
			if !errors.Is(err, c.err) {
				t.Errorf("Expected error %v, got %v", c.err, err)
			}
		})
	}
}

func TestID(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		record Record
		id     uint64
		err    error
	}{
		{
			name:   "valid id",
			record: Record{"id": uint64(123)},
			id:     123,
			err:    nil,
		},
		{
			name:   "id not uint64",
			record: Record{"id": 123.0},
			id:     0,
			err:    errID,
		},
		{
			name:   "missing id",
			record: Record{"name": "test"},
			id:     0,
			err:    errID,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			id, err := c.record.ID()
			if id != c.id || !errors.Is(err, c.err) {
				t.Errorf("Expected id %d and error %v, got id %d and error %v", c.id, c.err, id, err)
			}
		})
	}
}
