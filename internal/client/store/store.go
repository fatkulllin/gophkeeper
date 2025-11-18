package store

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fatkulllin/gophkeeper/logger"
	"github.com/fatkulllin/gophkeeper/model"
	bolt "go.etcd.io/bbolt"
	berrors "go.etcd.io/bbolt/errors"
	"go.uber.org/zap"
)

type BoltStore struct {
	db *bolt.DB
}

var bucketRecords = []byte("records")

var bucketUsers = []byte("users")

// NewBoltDB открывает или создаёт файл BoltDB
func NewBoltDB() (*BoltStore, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(configDir, "gophkeeper", "data.db")

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}

	db, err := bolt.Open(path, 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to open BoltDB: %w", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketRecords)

		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(bucketUsers)

		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}

	store := &BoltStore{
		db: db,
	}

	return store, nil
}

func (s *BoltStore) Clear() error {
	return s.db.Update(func(tx *bolt.Tx) error {

		if err := tx.DeleteBucket(bucketRecords); err != nil && err != berrors.ErrBucketNotFound {
			return err
		}
		if _, err := tx.CreateBucket(bucketRecords); err != nil {
			return err
		}

		if err := tx.DeleteBucket(bucketUsers); err != nil && err != berrors.ErrBucketNotFound {
			return err
		}
		if _, err := tx.CreateBucket(bucketUsers); err != nil {
			return err
		}

		return nil
	})
}

func (s *BoltStore) All() ([]model.Record, error) {
	records := []model.Record{}

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketRecords)

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

func (s *BoltStore) PutUserKey(userKey string) error {
	raw, err := base64.StdEncoding.DecodeString(userKey)
	if err != nil {
		return fmt.Errorf("decode key before save: %w", err)
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketUsers)
		return b.Put([]byte("userKey"), raw)
	})
}

func (s *BoltStore) GetUserKey() ([]byte, error) {
	var value []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketUsers)
		v := b.Get([]byte("userKey"))
		if v == nil {
			return fmt.Errorf("key not found")
		}
		value = append(value, v...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return value, nil
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

		return nil
	})
}
