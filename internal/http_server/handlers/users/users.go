package users

import (
	"EventSender/internal/models"
	"EventSender/internal/storage"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

var (
	userExist        = errors.New("user with such mail is already exist")
	userIdIsNotExist = errors.New("user with such id do not exist")
)

type Response struct {
	StatusResp int   `json:"statusResp"`
	Err        error `json:"err,omitempty"`
}

type userCreator interface {
	CreateUser(mail, password string) (int64, error)
}

type userChecker interface {
	CheckUser(id int64) (string, error)
	GetUserPassword(mail string) (string, error)
}

type productChecker interface {
	CheckProduct(ctx context.Context, productName string) (string, error)
}

func CreateUser(logger *slog.Logger, userCreator userCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := render.DecodeJSON(r.Body, &user)
		if err != nil {
			logger.Error("failed to decode JSON", err)
			response(w, r, http.StatusInternalServerError, err)
			return
		}

		err = validator.New().Struct(user)
		if err != nil {
			logger.Error("failed to validate JSON")
			response(w, r, http.StatusBadRequest, err)
			return
		}

		id, err := userCreator.CreateUser(user.Mail, user.Password)
		if err != nil {
			if errors.Is(err, storage.ErrUserAlreadyExist) {
				logger.Info("user is already exist")
				response(w, r, http.StatusBadRequest, userExist)
				return
			}
			logger.Error("failed to create user")
			response(w, r, http.StatusInternalServerError, err)
			return
		}

		logger.Info("user successfully created with id :", id)
		response(w, r, http.StatusOK, nil)
	}
}

func CheckUser(logger *slog.Logger, userChecker userChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userId models.UserWithId
		rawByte, err := io.ReadAll(r.Body)

		if err != nil {
			logger.Error("failed to read request body", err)
			response(w, r, http.StatusInternalServerError, err)
			return
		}
		err = json.Unmarshal(rawByte, &userId)

		if err != nil {
			logger.Error("failed to decode JSON", err)
			response(w, r, http.StatusBadRequest, err)
			return
		}
		mail, err := userChecker.CheckUser(userId.UserId)

		if err != nil {
			if errors.Is(err, storage.ErrUserNotExist) {
				logger.Info("user with such id is not exist")
				response(w, r, http.StatusNotFound, userIdIsNotExist)
				return
			}
			logger.Error("failed to find user", err)
			response(w, r, http.StatusInternalServerError, err)
			return
		}

		logger.Info("user was found with mail:", mail)
		response(w, r, http.StatusOK, err)
	}
}

func GetUserPassword(logger *slog.Logger, userChecker userChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.UserWithMail
		rawByte, err := io.ReadAll(r.Body)

		if err != nil {
			logger.Error("failed to read request body", err)
			response(w, r, http.StatusInternalServerError, err)
			return
		}

		err = json.Unmarshal(rawByte, &user)
		if err != nil {
			logger.Error("failed to decode JSON", err)
			response(w, r, http.StatusBadRequest, err)
			return
		}

		password, err := userChecker.GetUserPassword(user.Mail)

		if err != nil {
			if errors.Is(err, storage.ErrUserNotExist) {
				logger.Info("user with such mail is not exist")
				response(w, r, http.StatusNotFound, userIdIsNotExist)
				return
			}
			logger.Error("failed to find user", err)
			response(w, r, http.StatusInternalServerError, err)
			return
		}

		logger.Info("user was found with mail:", user.Mail)
		logger.Info("his password is :", password)
		response(w, r, http.StatusOK, err)
	}
}

func BuyProduct(logger *slog.Logger, userChecker userChecker, productChecker productChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userAndProductId models.UserIdAndProductName
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

		_, err = productChecker.CheckProduct(context.Background(), userAndProductId.Name)
		if err != nil {
			logger.Error("failed to find product")
			response(w, r, http.StatusNotFound, err)
			return
		}

		requestCache, err := http.NewRequest(http.MethodPost, "http://localhost:8081/cache/users", bytes.NewBuffer(rawByte))

		if err != nil {
			logger.Error("failed to form request")
			response(w, r, http.StatusInternalServerError, err)
			return
		}
		requestCache.Header.Set("Content-Type", "application/json")

		byteProductName, err := json.Marshal(userAndProductId.Name)
		if err != nil {
			logger.Error("failed to marshal json")
			response(w, r, http.StatusInternalServerError, err)
			return
		}

		requestWordsCache, err := http.NewRequest(http.MethodPost, "http://localhost:8081/cache/words", bytes.NewBuffer(byteProductName))
		if err != nil {
			logger.Error("failed to form request")
			response(w, r, http.StatusInternalServerError, err)
			return
		}

		requestWordsCache.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(requestCache)
		if err != nil {
			logger.Error("failed to do request", err)
			response(w, r, http.StatusInternalServerError, err)
			return
		}
		_, err = client.Do(requestWordsCache)
		if err != nil {
			logger.Error("failed to do request to words", err)
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
