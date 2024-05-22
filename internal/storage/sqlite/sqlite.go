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
	db              *sql.DB
	stmtCreate      *sql.Stmt
	stmtCheck       *sql.Stmt
	stmtGetPassword *sql.Stmt
}

func MustSetupDB(logger *slog.Logger, config *config.Config) (*Storage, error) {
	db, err := sql.Open("sqlite3", config.Storage)
	if err != nil {
		logger.Error("failed to create db", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		logger.Error("database is not response")
		return nil, err
	}
	stmtCreate, err := db.Prepare(`INSERT INTO mails(email, password) VALUES (?, ?)`)
	if err != nil {
		logger.Error("invalid statement CREATE", err)
		return nil, err
	}

	stmtCheck, err := db.Prepare(`SELECT email FROM mails WHERE id = ?`)
	if err != nil {
		logger.Error("invalid statement CHECK", err)
		return nil, err
	}

	stmtGetPassword, err := db.Prepare(`SELECT password FROM mails WHERE email = ?`)
	if err != nil {
		logger.Error("invalid statement GET PASSWORD", err)
		return nil, err
	}

	return &Storage{db: db,
			stmtCreate:      stmtCreate,
			stmtCheck:       stmtCheck,
			stmtGetPassword: stmtGetPassword},
		nil
}

func (s *Storage) CreateUser(mail, password string) (int64, error) {
	res, err := s.stmtCreate.Exec(mail, password)
	if err != nil {
		var errSqlite sqlite3.Error
		if errors.As(err, &errSqlite); errors.Is(errSqlite, sqlite3.ErrConstraint) {
			return 0, fmt.Errorf("failed to create user: %w", storage.ErrUserAlreadyExist)
		}
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
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

func (s *Storage) GetUserPassword(mail string) (string, error) {
	var password string
	err := s.stmtGetPassword.QueryRow(mail).Scan(&password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user is not exist : %w", storage.ErrUserNotExist)
		}
		return "", err
	}
	return password, nil
}
