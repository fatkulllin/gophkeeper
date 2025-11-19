package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fatkulllin/gophkeeper/model"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

// PGRepo предоставляет методы для работы с таблицами пользователей и записей
// в базе данных Postgres.
type RecordRepo struct {
	db *sql.DB
}

// NewPGRepo создаёт подключение к базе данных Postgres по переданному DSN.
// Выполняется проверка соединения через PingContext.
// Возвращает репозиторий или ошибку при инициализации.
func NewRecordRepo(db *sql.DB) *RecordRepo {

	return &RecordRepo{db: db}
}

// CreateRecord добавляет новую запись пользователя.
func (s *RecordRepo) CreateRecord(ctx context.Context, record model.Record) error {

	_, err := s.db.ExecContext(ctx, "INSERT INTO records (user_id, type, metadata, data) VALUES ($1, $2, $3, $4)", record.UserID, record.Type, record.Metadata, record.Data)

	if err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	return nil
}

// GetAllRecords возвращает все записи пользователя.
func (s *RecordRepo) GetAllRecords(ctx context.Context, userID int) ([]model.Record, error) {
	records := make([]model.Record, 0)
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, type, metadata, data
		FROM records
		WHERE user_id = $1
		ORDER BY created_at DESC
		`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var r model.Record
		err = rows.Scan(&r.ID, &r.UserID, &r.Type, &r.Metadata, &r.Data)
		if err != nil {
			return nil, err
		}

		records = append(records, r)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	logger.Log.Debug("records fetched", zap.Int("count", len(records)))
	return records, nil
}

// DeleteRecord удаляет запись по ID.
func (s *RecordRepo) DeleteRecord(ctx context.Context, userID int, idRecord string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM records WHERE id = $1 AND user_id = $2 ", idRecord, userID)

	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("record not found: id=%s", idRecord)
	}
	return nil
}

// GetRecord возвращает запись по её ID.
func (s *RecordRepo) GetRecord(ctx context.Context, userID int, idRecord string) (model.Record, error) {
	var record model.Record
	row := s.db.QueryRowContext(ctx, `
		SELECT id, type, metadata, data
		FROM records
		WHERE user_id = $1 AND id = $2
		`, userID, idRecord)
	err := row.Scan(&record.ID, &record.Type, &record.Metadata, &record.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Record{}, fmt.Errorf("record not found for user %v", userID)
		}
		return model.Record{}, err
	}

	return record, nil
}

// UpdateRecord обновляет метаданные и/или данные записи.
func (s *RecordRepo) UpdateRecord(ctx context.Context, userID int, idRecord string, record model.Record) error {

	if record.Metadata == "" && record.Data == nil {
		return nil
	}

	query := "UPDATE records SET "
	args := []any{}
	idx := 1

	if record.Metadata != "" {
		query += fmt.Sprintf("metadata = $%d", idx)
		args = append(args, record.Metadata)
		idx++
	}

	if record.Data != nil {
		if len(args) > 0 {
			query += ", "
		}
		query += fmt.Sprintf("data = $%d", idx)
		args = append(args, record.Data)
		idx++
	}

	query += fmt.Sprintf(", updated_at = NOW() WHERE id = $%d AND user_id = $%d", idx, idx+1)
	args = append(args, idRecord, userID)
	logger.Log.Debug("run query update", zap.String("query", query), zap.Any("args", args))

	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}
	return nil
}
