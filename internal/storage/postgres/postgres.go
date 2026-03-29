package postgres

import (
	"context"
	"errors"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	CreateTableQuery = `
	CREATE TABLE IF NOT EXISTS url(
		id BIGSERIAL PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`

	InsertQuery = `
	INSERT INTO url (alias, url)
	VALUES ($1, $2)
	`

	SelectQuery = `
	SELECT url FROM url WHERE alias = $1;
	`

	DeleteQuery = `
	DELETE FROM url WHERE alias = $1;
	`
)

type Storage struct {
	conn *pgx.Conn
}

func NewStorage(ctx context.Context, databaseUrl string) (*Storage, error) {
	const op = "storage.Postgres.NewStorage"

	pgConn, err := pgx.ParseConfig(databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	conn, err := pgx.ConnectConfig(ctx, pgConn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = conn.Exec(ctx, CreateTableQuery)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{conn: conn}, nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *Storage) SaveUrl(ctx context.Context, urlToSave string, alias string) error {
	const op = "storage.Postgres.SaveUrl"

	_, err := s.conn.Exec(ctx, InsertQuery, alias, urlToSave)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, storage.ErrorAliasAlreadyExists)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetUrl(ctx context.Context, alias string) (string, error) {
	const op = "storage.Postgres.GetUrl"

	var url string

	err := s.conn.QueryRow(ctx, SelectQuery, alias).Scan(&url)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrorAliasNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

func (s *Storage) DeleteUrl(ctx context.Context, alias string) error {
	const op = "storage.Postgres.DeleteUrl"

	_, err := s.conn.Exec(ctx, DeleteQuery, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
