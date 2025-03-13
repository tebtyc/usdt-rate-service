package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"usdt-rate-service/internal/domain"

	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
}

func New(dbURL string) (*Storage, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("postgres.New: %w", err)
	}

	return &Storage{DB: db}, nil
}

func (s *Storage) Close() error {
	return s.DB.Close()
}

func (s *Storage) SaveRates(ctx context.Context, rate domain.Rate) error {
	const query = `INSERT INTO rates (ask, bid, timestamp) VALUES ($1, $2, $3)`

	_, err := s.DB.ExecContext(ctx, query, rate.Ask, rate.Bid, rate.Timestamp)
	if err != nil {
		return fmt.Errorf("postgres.SaveRates: %w", err)
	}

	return nil
}
