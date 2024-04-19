package util

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

func DoWithTries(fn func() (*pgxpool.Pool, error), attempts int, delay time.Duration) (pool *pgxpool.Pool, err error) {
	for attempts > 0 {
		if pool, err = fn(); err != nil {
			time.Sleep(delay)
			attempts--
			continue
		}
		return pool, nil
	}
	return
}
