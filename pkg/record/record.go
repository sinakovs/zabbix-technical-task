package record

type Record map[string]interface{}

type RecordStore struct {
	records map[string]Record
}

func NewRecordStore() *RecordStore {
	return &RecordStore{
		records: make(map[string]Record),
	}
}

func (r *RecordStore) Add(id string, record Record) {
	r.records[id] = record
}

func (r *RecordStore) Get(id string) (Record, bool) {
	record, exists := r.records[id]
	return record, exists
}
