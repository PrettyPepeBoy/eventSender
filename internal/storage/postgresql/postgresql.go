package postgresql

import (
	cfg "EventSender/config"
	"EventSender/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log/slog"
)

type Storage struct {
	db         *pgxpool.Pool
	createStmt string
	checkStmt  string
}

func MustConnectDB(ctx context.Context, config cfg.Postgresql, logger *slog.Logger) (*Storage, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		logger.Error("failed to connect to postgres", err)
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		logger.Error("database is not response")
		return nil, err
	}

	createStmt := fmt.Sprintf("INSERT INTO products(name, category) VALUES ($1, $2) RETURNING id")
	checkStmt := fmt.Sprintf("SELECT id  FROM products WHERE name = $1 ")
	return &Storage{db: pool, createStmt: createStmt, checkStmt: checkStmt}, nil
}

func (s *Storage) CreateProduct(ctx context.Context, product models.Product) (string, error) {
	var id string
	row := s.db.QueryRow(ctx, s.createStmt, product.Name, product.Category)
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (s *Storage) CheckProduct(ctx context.Context, productName string) (string, error) {
	var id string
	if err := s.db.QueryRow(ctx, s.checkStmt, productName).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
		return "", err
	}
	return id, nil
}
