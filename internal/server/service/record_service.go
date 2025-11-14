package service

import (
	"context"
	"fmt"

	"github.com/fatkulllin/gophkeeper/model"
)

type RecordService struct {
	repo RecordRepository
}

func NewRecordService(repo RecordRepository) *RecordService {
	return &RecordService{
		repo: repo,
	}
}

func (s RecordService) Create(ctx context.Context, userID int64, input model.Record) error {
	record := model.Record{
		UserID:   userID,
		Type:     input.Type,
		Metadata: input.Metadata,
		Data:     input.Data, // можно будет шифровать здесь
	}
	if err := s.repo.CreateRecord(ctx, record); err != nil {
		return fmt.Errorf("failed to create record: %w", err)
	}
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
