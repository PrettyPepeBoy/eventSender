package products

import (
	"EventSender/internal/models"
	"context"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Response struct {
	StatusResp int   `json:"statusResp"`
	Err        error `json:"err,omitempty"`
}

type productCreator interface {
	CreateProduct(ctx context.Context, product models.Product) (string, error)
}

func CreateProduct(logger *slog.Logger, ctx context.Context, creator productCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		err := render.DecodeJSON(r.Body, &product)
		if err != nil {
			logger.Error("failed to decode JSON", err)
			response(w, r, http.StatusInternalServerError, err)
			return
		}

		id, err := creator.CreateProduct(ctx, product)
		if err != nil {
			logger.Error("failed to create product", err)
			response(w, r, http.StatusBadRequest, err)
			return
		}
		logger.Info("created new product with id:", id)
		response(w, r, http.StatusOK, nil)
	}
}

func response(w http.ResponseWriter, r *http.Request, status int, err error) {
	render.JSON(w, r, Response{
		StatusResp: status,
		Err:        err,
	})
}
