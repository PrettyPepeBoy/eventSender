package users

import (
	"EventSender/internal/models"
	"EventSender/internal/storage"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

type Response struct {
	StatusResp int   `json:"statusResp"`
	Err        error `json:"err,omitempty"`
}

type userCreator interface {
	CreateUser(mail string) error
}

type userChecker interface {
	CheckUser(id int64) (string, error)
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
			logger.Error("failed to validate JSON")
			response(w, r, http.StatusBadRequest, err)
		}

		err = userCreator.CreateUser(user.Mail)
		if err != nil {
			if errors.Is(err, storage.ErrUserAlreadyExist) {
				logger.Info("user is already exist")
				response(w, r, http.StatusBadRequest, err)
			}
			logger.Error("failed to create user")
			response(w, r, http.StatusBadRequest, err)
		}

		response(w, r, http.StatusOK, nil)
	}
}

func BuyProduct(logger *slog.Logger, userChecker userChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userAndProductId models.UserAndProductId
		rawByte, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("failed to read request body", err)
			response(w, r, http.StatusBadRequest, err)
			return
		}
		err = json.Unmarshal(rawByte, &userAndProductId)
		if err != nil {
			logger.Error("failed to decode JSON", err)
			response(w, r, http.StatusBadRequest, err)
			return
		}
		err = validator.New().Struct(userAndProductId)
		if err != nil {
			logger.Error("failed to validate JSON")
			response(w, r, http.StatusBadRequest, err)
			return
		}

		_, err = userChecker.CheckUser(userAndProductId.UserId)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotExist) {
				logger.Info("user with such id do not exist")
				response(w, r, http.StatusNotFound, err)
				return
			}
			logger.Error("failed to find user")
			response(w, r, http.StatusNotFound, err)
			return
		}

		request, err := http.NewRequest(http.MethodPost, "http://localhost:8081/cache/users", bytes.NewBuffer(rawByte))

		if err != nil {
			logger.Error("failed to form request")
			response(w, r, http.StatusInternalServerError, err)
			return
		}
		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			logger.Error("failed to do request", err)
			response(w, r, http.StatusInternalServerError, err)
			return
		}
		respRawByte, err := io.ReadAll(resp.Body)

		defer resp.Body.Close()
		if err != nil {
			logger.Error("failed to read response")
			logger.Info("response :", string(respRawByte))
			response(w, r, resp.StatusCode, err)
			return
		}
		logger.Info("response :", string(respRawByte))
		response(w, r, resp.StatusCode, err)
	}
}

func response(w http.ResponseWriter, r *http.Request, status int, err error) {
	render.JSON(w, r, Response{
		StatusResp: status,
		Err:        err,
	})
}
