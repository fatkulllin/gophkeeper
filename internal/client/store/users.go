package store

import (
	"encoding/base64"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

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
