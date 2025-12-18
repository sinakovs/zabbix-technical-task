package userrecord

import "errors"

var (
	errNoID        = errors.New("missing id")
	errIDNotNumber = errors.New("id must be a number")
	errIDNotUint   = errors.New("'id' must be a non-negative integer")
	errID          = errors.New("id is not exist or is not a uint64")
)

// Record represents a generic record with dynamic fields.
type Record map[string]any // interface{}

// Validate validates the record to ensure it has a valid 'id' field.
func (r Record) Validate() error {
	id, ok := r["id"]
	if !ok {
		return errNoID
	}

	idFloat, ok := id.(float64)
	if !ok {
		return errIDNotNumber
	}

	if idFloat < 0 || idFloat != float64(int(idFloat)) {
		return errIDNotUint
	}

	r["id"] = uint64(idFloat)

	return nil
}

// ID returns the ID of the record as uint64.
func (r Record) ID() (uint64, error) {
	id, ok := r["id"].(uint64)
	if !ok {
		return 0, errID
	}

	return id, nil
}
