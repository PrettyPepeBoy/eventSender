package main

import (
	"EventSender/config"
	"EventSender/internal/http_server/handlers/battlegrounds"
	"EventSender/internal/lib/logger/handlers/slogpretty"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"os"
)

const envLocal = "local"

func main() {
	cfg := initConfig()
	logger := initLogger(cfg.Env)
	router := initRouter(cfg, logger)
}

func initRouter(cfg *config.Config, logger *slog.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/battlegrounds", func(r chi.Router) {
		r.Post("/invite", battlegrounds.SendInvite(logger, cfg))
	})

	return router
}

func initConfig() *config.Config {
	return config.MustLoad()
}

func initLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlersOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
