package benchmark

import (
	"EventSender/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type Response struct {
	StatusResp int   `json:"statusResp"`
	Err        error `json:"err,omitempty"`
}

func SendManyRequestsInFunction(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userAndProduct []models.UserIdAndProductName

		rawByte, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("failed to read request body", err)
			response(w, r, http.StatusInternalServerError, err)
			return
		}

		err = json.Unmarshal(rawByte, &userAndProduct)
		if err != nil {
			logger.Error("failed to unmarshal json", err)
			response(w, r, http.StatusBadRequest, err)
			return
		}

		wg := sync.WaitGroup{}
		wg.Add(len(userAndProduct))
		for i, elem := range userAndProduct {
			i := i
			elem := elem
			go func() {
				jsonByte, err := json.Marshal(elem)
				if err != nil {
					logger.Error("failed to marshal json", err)
					response(w, r, http.StatusBadRequest, err)
					return
				}

				requestCache, err := http.NewRequest(http.MethodPost, "http://localhost:8081/cache/users", bytes.NewBuffer(jsonByte))

				if err != nil {
					logger.Error("failed to form request")
					response(w, r, http.StatusInternalServerError, err)
					return
				}
				requestCache.Header.Set("Content-Type", "application/json")

				client := http.Client{}
				t := time.Now()
				resp, err := client.Do(requestCache)

				if err != nil {
					logger.Error("failed to send request to cache", err)
					response(w, r, http.StatusBadRequest, err)
					return
				}

				logger.Info(fmt.Sprintf("successfully send request %v", i))

				respByte, err := io.ReadAll(resp.Body)
				if err != nil {
					logger.Error("failed to read response body", err)
					response(w, r, http.StatusInternalServerError, err)
					return
				}
				logger.Info("resp byte : ", slog.String("body", string(respByte)))
				logger.Info("time spend :", slog.String("time", time.Since(t).String()))

				wg.Done()
			}()
		}
		wg.Wait()
		response(w, r, http.StatusOK, err)
	}
}

func response(w http.ResponseWriter, r *http.Request, status int, err error) {
	render.JSON(w, r, Response{
		StatusResp: status,
		Err:        err,
	})
}
