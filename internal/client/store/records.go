package store

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fatkulllin/gophkeeper/model"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

func (s *BoltStore) All() ([]model.Record, error) {
	records := []model.Record{}

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketRecords)
		if b == nil {
			return fmt.Errorf("bucket 'records' not found")
		}
		return b.ForEach(func(k, v []byte) error {
			var rec model.Record
			if err := json.Unmarshal(v, &rec); err != nil {
				return err
			}
			records = append(records, rec)
			return nil
		})
	})

	return records, err
}

func (s *BoltStore) Get(id int64) (model.Record, error) {
	var rec model.Record

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketRecords)
		if b == nil {
			return fmt.Errorf("bucket 'records' not found")
		}
		data := b.Get([]byte(strconv.FormatInt(id, 10)))
		if data == nil {
			return fmt.Errorf("record not found")
		}
		return json.Unmarshal(data, &rec)
	})

	return rec, err
}

func (s *BoltStore) Put(r model.Record) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketRecords)

		data, err := json.Marshal(r)
		if err != nil {
			return err
		}

		key := []byte(strconv.FormatInt(r.ID, 10))
		return b.Put(key, data)
	})
}

func (s *BoltStore) SaveRecords(records []model.Record) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketRecords)
		if b == nil {
			return fmt.Errorf("bucket 'records' not found")
		}

		for _, rec := range records {
			data, err := json.Marshal(rec)
			if err != nil {
				return fmt.Errorf("marshal error: %w", err)
			}

			key := []byte(strconv.FormatInt(rec.ID, 10))

			if err := b.Put(key, data); err != nil {
				return fmt.Errorf("put error for id %d: %w", rec.ID, err)
			}
			logger.Log.Debug("save record", zap.ByteString("record", data))
		}
		logger.Log.Debug("saved all records", zap.Int("count", len(records)))
		return nil
	})
}
