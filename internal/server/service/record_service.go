package service

import (
	"context"
	"fmt"

	"github.com/fatkulllin/gophkeeper/logger"
	"github.com/fatkulllin/gophkeeper/model"
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
		Data:     input.Data, // можно будет шифровать здесь
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

	// decryptData, err := s.cryptoUtil.Decrypt(encryptData, decryptUserKey)
	// if err != nil {
	// 	logger.Log.Error("", zap.Error(err))
	// 	return err
	// }

	logger.Log.Debug("decrypt data", zap.ByteString("data", record.Data))
	err = s.repo.CreateRecord(ctx, record)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return err
	}
	// if err := s.repo.CreateRecord(ctx, record); err != nil {
	// 	return fmt.Errorf("failed to create record: %w", err)
	// }
	return nil
}

func (s RecordService) GetAll(ctx context.Context, userID int) ([]model.Record, error) {
	records, err := s.repo.GetRecords(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get records: %w", err)
	}
	return records, nil
}

func (s RecordService) Delete(ctx context.Context, userID, recordID int) error {
	if err := s.repo.DeleteRecord(ctx, recordID, userID); err != nil {
		return fmt.Errorf("delete record: %w", err)
	}
	return nil
}
