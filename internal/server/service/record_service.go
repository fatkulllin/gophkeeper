// Package service содержит бизнес-логику работы с пользователями и их записями.

package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fatkulllin/gophkeeper/model"
	"github.com/fatkulllin/gophkeeper/pkg/cryptoutil"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"go.uber.org/zap"
)

// RecordService отвечает за логику создания, чтения, обновления
// и удаления записей пользователя. Он выполняет шифрование и расшифровку
// данных с использованием пользовательского ключа.
type RecordService struct {
	recordRepo RecordRepositories
	userRepo   UserRepositories
	cryptoUtil CryptoUtil
}

// NewRecordService создаёт новый сервис для работы с записями.
func NewRecordService(recordRepo RecordRepositories, userRepo UserRepositories, cryptoUtil CryptoUtil) *RecordService {
	return &RecordService{
		recordRepo: recordRepo,
		userRepo:   userRepo,
		cryptoUtil: cryptoUtil,
	}
}

func (s *RecordService) Create(ctx context.Context, userID int, input model.RecordInput) error {
	record := model.Record{
		UserID:   userID,
		Type:     input.Type,
		Metadata: input.Metadata,
	}

	encryptedKey, err := s.userRepo.GetEncryptedKeyUser(ctx, userID)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return err
	}

	decryptUserKey, err := s.cryptoUtil.DecryptWithMasterKey(encryptedKey)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return err
	}

	encryptData, err := cryptoutil.Encrypt(input.Data, decryptUserKey)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return err
	}

	record.Data = encryptData

	logger.Log.Debug("encrypted data", zap.ByteString("data", encryptData))
	err = s.recordRepo.CreateRecord(ctx, record)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return fmt.Errorf("failed to create record: %w", err)
	}

	return nil
}

func (s *RecordService) GetAll(ctx context.Context, userID int) ([]model.Record, error) {

	records, err := s.recordRepo.GetAllRecords(ctx, userID)
	if err != nil {
		logger.Log.Error("error", zap.Error(err))
		return nil, fmt.Errorf("get records: %w", err)
	}
	logger.Log.Debug("records", zap.Any("records", records))
	return records, nil
}

func (s *RecordService) Get(ctx context.Context, userID int, idRecord string) (model.RecordResponse, error) {
	record, err := s.recordRepo.GetRecord(ctx, userID, idRecord)

	if err != nil {
		logger.Log.Error("get record error", zap.Error(err))
		return model.RecordResponse{}, fmt.Errorf("get record: %w", err)
	}

	encryptedKey, err := s.userRepo.GetEncryptedKeyUser(ctx, userID)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return model.RecordResponse{}, err
	}

	decryptUserKey, err := s.cryptoUtil.DecryptWithMasterKey(encryptedKey)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return model.RecordResponse{}, err
	}
	decryptData, err := cryptoutil.Decrypt(record.Data, decryptUserKey)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return model.RecordResponse{}, err
	}

	return model.RecordResponse{
		ID:       record.ID,
		Type:     record.Type,
		Metadata: record.Metadata,
		Data:     decryptData,
	}, nil
}

func (s *RecordService) Delete(ctx context.Context, userID int, idRecord string) error {
	if err := s.recordRepo.DeleteRecord(ctx, userID, idRecord); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Debug("no rows for delete", zap.String("record id", idRecord), zap.Int("user id", userID))
			return sql.ErrNoRows
		}
		logger.Log.Error("", zap.Error(err))
		return fmt.Errorf("delete record: %w", err)
	}
	return nil
}

func (s *RecordService) Update(ctx context.Context, userID int, idRecord string, input model.RecordUpdateInput) error {
	var record model.Record
	if input.Metadata == nil && input.Data == nil {
		return errors.New("nothing to update: both metadata and data are nil")
	}

	if input.Data != nil {
		encryptedKey, err := s.userRepo.GetEncryptedKeyUser(ctx, userID)
		if err != nil {
			logger.Log.Error("", zap.Error(err))
			return err
		}

		decryptUserKey, err := s.cryptoUtil.DecryptWithMasterKey(encryptedKey)
		if err != nil {
			logger.Log.Error("", zap.Error(err))
			return err
		}

		encryptData, err := cryptoutil.Encrypt(*input.Data, decryptUserKey)
		if err != nil {
			logger.Log.Error("", zap.Error(err))
			return err
		}
		record.Data = []byte(encryptData)
	}
	if input.Metadata != nil {
		record.Metadata = *input.Metadata
	}
	err := s.recordRepo.UpdateRecord(ctx, userID, idRecord, record)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return err
	}
	return nil
}
