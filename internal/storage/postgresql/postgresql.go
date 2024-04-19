package postgresql

import (
	cfg "EventSender/config"
	"EventSender/internal/util"
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error
}

func MustConnectDB(ctx context.Context, maxAttempts int, delay time.Duration, config cfg.StorageConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s:%s", config.Host, config.Port, config.Database, config.Username, config.Password)
	pool, err := util.DoWithTries(func() (*pgxpool.Pool, error) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err := pgxpool.Connect(ctx, dsn)
		if err != nil {
			return nil, err
		}
		return pool, nil
	}, maxAttempts, delay)

	if err != nil {
		log.Fatal("failed to connect to postgresql db")
	}
	return pool, err
}
