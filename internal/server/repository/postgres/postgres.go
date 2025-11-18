package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/fatkulllin/gophkeeper/model"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

// PGRepo предоставляет методы для работы с таблицами пользователей и записей
// в базе данных Postgres.
type PGRepo struct {
	conn *sql.DB
}

// NewPGRepo создаёт подключение к базе данных Postgres по переданному DSN.
// Выполняется проверка соединения через PingContext.
// Возвращает репозиторий или ошибку при инициализации.
func NewPGRepo(dsn string) (*PGRepo, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return &PGRepo{}, fmt.Errorf("failed to open db: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return &PGRepo{}, fmt.Errorf("db ping failed: %w", err)
	}

	return &PGRepo{conn: db}, nil
}

// Bootstrap применяет миграции goose, используя встроенную файловую систему.
// Используется при старте приложения.
func (s *PGRepo) Bootstrap(fs embed.FS) error {
	goose.SetBaseFS(fs)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set postgres dialect: %w", err)
	}

	if err := goose.Up(s.conn, "."); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}

// ExistUser проверяет наличие пользователя по логину.
// Возвращает true, если пользователь существует.
func (s *PGRepo) ExistUser(ctx context.Context, user model.UserCredentials) (bool, error) {
	row := s.conn.QueryRowContext(ctx, "SELECT login FROM users WHERE login = $1", user.Username)
	var userScan string
	err := row.Scan(&userScan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("check existing user: %w", err)
	}
	return true, nil
}

// CreateUser создаёт нового пользователя и возвращает его ID.
func (s *PGRepo) CreateUser(ctx context.Context, user model.UserCredentials) (int, error) {

	var id int

	row := s.conn.QueryRowContext(ctx, "INSERT INTO users (login, password_hash, encrypted_key) VALUES ($1, $2, $3) RETURNING id", user.Username, user.Password, user.EncryptedKey)

	err := row.Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("pg failed to insert new user: %w", err)
	}

	return id, nil
}

// GetUser возвращает пользователя по логину.
func (s *PGRepo) GetUser(ctx context.Context, user model.UserCredentials) (model.User, error) {
	var foundUser model.User
	row := s.conn.QueryRowContext(ctx, "SELECT id, login, password_hash FROM users WHERE login = $1", user.Username)
	err := row.Scan(&foundUser.ID, &foundUser.Login, &foundUser.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("user not found")
		}
		return model.User{}, err
	}

	return foundUser, nil
}

// GetEncryptedKeyUser возвращает зашифрованный ключ пользователя.
func (s *PGRepo) GetEncryptedKeyUser(ctx context.Context, userID int) (string, error) {
	var encryptedKey string
	row := s.conn.QueryRowContext(ctx, "SELECT encrypted_key FROM users WHERE id = $1;", userID)
	err := row.Scan(&encryptedKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user not found")
		}
		return "", err
	}
	return encryptedKey, nil
}

// CreateRecord добавляет новую запись пользователя.
func (s *PGRepo) CreateRecord(ctx context.Context, record model.Record) error {

	_, err := s.conn.ExecContext(ctx, "INSERT INTO records (user_id, type, metadata, data) VALUES ($1, $2, $3, $4)", record.UserID, record.Type, record.Metadata, record.Data)

	if err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	return nil
}

// GetAllRecords возвращает все записи пользователя.
func (s *PGRepo) GetAllRecords(ctx context.Context, userID int) ([]model.Record, error) {
	records := make([]model.Record, 0)
	rows, err := s.conn.QueryContext(ctx, `
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
func (s *PGRepo) DeleteRecord(ctx context.Context, userID int, idRecord string) error {
	result, err := s.conn.ExecContext(ctx, "DELETE FROM records WHERE id = $1 AND user_id = $2 ", idRecord, userID)

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
func (s *PGRepo) GetRecord(ctx context.Context, userID int, idRecord string) (model.Record, error) {
	var record model.Record
	row := s.conn.QueryRowContext(ctx, `
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
func (s *PGRepo) UpdateRecord(ctx context.Context, userID int, idRecord string, record model.Record) error {

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

	_, err := s.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}
	return nil
}

// Close закрывает соединение с базой данных.
func (s *PGRepo) Close() error {
	return s.conn.Close()
}
