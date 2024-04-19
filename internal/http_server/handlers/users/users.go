package users

import (
	"EventSender/internal/models"
	"context"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"os"
)

type Response struct {
	StatusResp int   `json:"statusResp"`
	Err        error `json:"err,omitempty"`
}

type userCreator interface {
	CreateUser(mail string) error
}

type tableCreator interface {
	CreateTable(ctx context.Context, logger *slog.Logger) error
}

func CreateUser(logger *slog.Logger, userCreator userCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := render.DecodeJSON(r.Body, &user)
		if err != nil {
			logger.Error("failed to decode JSON", err)
			response(w, r, http.StatusBadRequest, err)
		}

		err = validator.New().Struct(user)
		if err != nil {
			logger.Error("failed to decode JSON")
			response(w, r, http.StatusBadRequest, err)
		}

		err = userCreator.CreateUser(user.Mail)
		if err != nil {
			//TODO add mistake
			logger.Info("failed to create user")
			response(w, r, http.StatusBadRequest, err)
		}

		response(w, r, http.StatusOK, nil)
	}
}

func CreateTable(logger *slog.Logger, tableCreator tableCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context
		err := tableCreator.CreateTable(ctx, logger)
		if err != nil {
			logger.Info("failed to create database")
			response(w, r, http.StatusBadRequest, err)
		}
		os.Exit(1)
	}
}

func response(w http.ResponseWriter, r *http.Request, status int, err error) {
	render.JSON(w, r, Response{
		StatusResp: status,
		Err:        err,
	})
}
