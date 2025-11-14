package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/fatkulllin/gophkeeper/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type PGRepo struct {
	conn *sql.DB
}

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

func (s *PGRepo) Bootstrap(fs embed.FS) error {
	goose.SetBaseFS(fs)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("error set dialect postgres %w", err)
	}

	if err := goose.Up(s.conn, "."); err != nil {
		return fmt.Errorf("error run migrate %w", err)
	}
	return nil
}

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

func (s *PGRepo) CreateUser(ctx context.Context, user model.UserCredentials) (int, error) {

	var id int

	row := s.conn.QueryRowContext(ctx, "INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id", user.Username, user.Password)

	err := row.Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("pg failed to insert new user: %w", err)
	}

	return id, nil
}

func (s *PGRepo) GetUser(ctx context.Context, user model.UserCredentials) (model.User, error) {
	var foundUser model.User
	row := s.conn.QueryRowContext(ctx, "SELECT * FROM users WHERE login = $1", user.Username)
	err := row.Scan(&foundUser.ID, &foundUser.Login, &foundUser.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("user not found")
		}
		return model.User{}, err
	}
	return foundUser, nil
}

// func (s *PGRepo) CreateRecord(ctx context.Context, user model.RecordInput) (int, error) {

// 	var id int

// 	row := s.conn.QueryRowContext(ctx, "INSERT INTO records (login, password_hash) VALUES ($1, $2) RETURNING id", user.Login, user.Password)

// 	err := row.Scan(&id)

// 	if err != nil {
// 		return 0, fmt.Errorf("pg failed to insert new user: %w", err)
// 	}

// 	return id, nil
// }
