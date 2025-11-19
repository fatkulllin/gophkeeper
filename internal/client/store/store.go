package store

import (
	"fmt"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
	berrors "go.etcd.io/bbolt/errors"
)

type BoltStore struct {
	db *bolt.DB
}

var bucketRecords = []byte("records")

var bucketUsers = []byte("users")

// NewBoltDB открывает или создаёт файл BoltDB
func NewBoltDB(cfgDir string) (*BoltStore, error) {

	path := filepath.Join(cfgDir, "data.db")

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
