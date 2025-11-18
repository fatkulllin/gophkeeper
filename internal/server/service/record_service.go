package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fatkulllin/gophkeeper/model"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"go.uber.org/zap"
)

type RecordService struct {
	repo       RecordRepository
	cryptoUtil CryptoUtil
}

func NewRecordService(repo RecordRepository, cryptoUtil CryptoUtil) *RecordService {
	return &RecordService{
		repo:       repo,
		cryptoUtil: cryptoUtil,
	}
}

func (s RecordService) Create(ctx context.Context, userID int, input model.RecordInput) error {
	record := model.Record{
		UserID:   userID,
		Type:     input.Type,
		Metadata: input.Metadata,
		Data:     input.Data,
	}

	encryptedKey, err := s.repo.GetEncryptedKeyUser(ctx, userID)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return err
	}

	decryptUserKey, err := s.cryptoUtil.DecryptWithMasterKey(encryptedKey)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return err
	}

	encryptData, err := s.cryptoUtil.EncryptString(input.Data, decryptUserKey)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return err
	}

	record.Data = []byte(encryptData)

	logger.Log.Debug("decrypt data", zap.ByteString("data", record.Data))
	err = s.repo.CreateRecord(ctx, record)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return fmt.Errorf("failed to create record: %w", err)
	}

	return nil
}

func (s RecordService) GetAll(ctx context.Context, userID int) ([]model.Record, error) {

	records, err := s.repo.GetAllRecords(ctx, userID)
	if err != nil {
		logger.Log.Error("error", zap.Error(err))
		return nil, fmt.Errorf("get records: %w", err)
	}
	logger.Log.Debug("records", zap.Any("records", records))
	return records, nil
}

func (s RecordService) Get(ctx context.Context, userID int, idRecord string) (model.RecordResponse, error) {
	record, err := s.repo.GetRecord(ctx, userID, idRecord)

	if err != nil {
		logger.Log.Error("get record error", zap.Error(err))
		return model.RecordResponse{}, fmt.Errorf("get record: %w", err)
	}

	encryptedKey, err := s.repo.GetEncryptedKeyUser(ctx, userID)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return model.RecordResponse{}, err
	}

	decryptUserKey, err := s.cryptoUtil.DecryptWithMasterKey(encryptedKey)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return model.RecordResponse{}, err
	}

	decryptData, err := s.cryptoUtil.Decrypt(string(record.Data), decryptUserKey)
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

func (s RecordService) Delete(ctx context.Context, userID int, idRecord string) error {
	if err := s.repo.DeleteRecord(ctx, userID, idRecord); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Debug("no rows dor delete", zap.String("record id", idRecord), zap.Int("user id", userID))
			return sql.ErrNoRows
		}
		logger.Log.Error("", zap.Error(err))
		return fmt.Errorf("delete record: %w", err)
	}
	return nil
}

func (s RecordService) Update(ctx context.Context, userID int, idRecord string, input model.RecordUpdateInput) error {
	var record model.Record
	if input.Metadata == nil && input.Data == nil {
		return errors.New("nothing to update")
	}

	if input.Data != nil {
		encryptedKey, err := s.repo.GetEncryptedKeyUser(ctx, userID)
		if err != nil {
			logger.Log.Error("", zap.Error(err))
			return err
		}

		decryptUserKey, err := s.cryptoUtil.DecryptWithMasterKey(encryptedKey)
		if err != nil {
			logger.Log.Error("", zap.Error(err))
			return err
		}

		encryptData, err := s.cryptoUtil.EncryptString(*input.Data, decryptUserKey)
		if err != nil {
			logger.Log.Error("", zap.Error(err))
			return err
		}
		record.Data = []byte(encryptData)
	}
	if input.Metadata != nil {
		record.Metadata = *input.Metadata
	}
	err := s.repo.UpdateRecord(ctx, userID, idRecord, record)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return err
	}
	return nil
}
