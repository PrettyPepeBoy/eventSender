package sqlite

import (
	"EventSender/config"
	"EventSender/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"

	"log/slog"
)

type Storage struct {
	db         *sql.DB
	stmtCreate *sql.Stmt
	stmtCheck  *sql.Stmt
}

func MustSetupDB(logger *slog.Logger, config *config.Config) (*Storage, error) {
	db, err := sql.Open("sqlite3", config.Storage)
	if err != nil {
		logger.Error("failed to create db", err)
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users(
		   id INTEGER PRIMARY KEY,
		   mail TEXT UNIQUE NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_mail ON users(mail)
	`)
	if err != nil {
		logger.Error("failed to connect with db", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		logger.Error("database is not response")
		return nil, err
	}
	stmtCreate, err := db.Prepare(`INSERT INTO users(mail) VALUES (?)`)
	if err != nil {
		logger.Error("invalid statement CREATE", err)
		return nil, err
	}

	stmtCheck, err := db.Prepare(`SELECT mail FROM users WHERE id = ?`)
	if err != nil {
		logger.Error("invalid statement CHECK", err)
		return nil, err
	}

	return &Storage{db: db,
			stmtCreate: stmtCreate,
			stmtCheck:  stmtCheck},
		nil
}

func (s *Storage) CreateUser(mail string) error {
	_, err := s.stmtCreate.Exec(mail)
	if err != nil {
		var errSqlite sqlite3.Error
		if errors.As(err, &errSqlite); errors.Is(errSqlite, sqlite3.ErrConstraint) {
			return fmt.Errorf("failed to create user: %w", storage.ErrUserAlreadyExist)
		}
		return err
	}
	return nil
}

func (s *Storage) CheckUser(id int64) (string, error) {
	var mail string
	err := s.stmtCheck.QueryRow(id).Scan(&mail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user is not exist : %w", storage.ErrUserNotExist)
		}
		return "", err
	}
	return mail, nil
}
