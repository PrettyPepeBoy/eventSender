package main

import (
	cfg "EventSender/config"
	"EventSender/internal/http_server/handlers/users"
	"EventSender/internal/lib/logger/handlers/slogpretty"
	"EventSender/internal/storage/sqlite"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const envLocal = "local"

func main() {
	config := cfg.MustLoad()
	fmt.Println(config)
	logger := setupLogger(config.Env)
	router := chi.NewRouter()
	dataBase, err := sqlite.MustSetupDB(logger, config)
	if err != nil {
		logger.Error("failed to setup database")
		os.Exit(1)
	}

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/users", func(r chi.Router) {
		r.Post("/", users.CreateUser(logger, dataBase))
		r.Post("/cache", users.BuyProduct(logger, dataBase))
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Info("starting server", slog.String("address", config.Server.Port))
	srv := &http.Server{
		Addr:         config.Server.Port,
		Handler:      router,
		ReadTimeout:  config.Server.Timeout,
		WriteTimeout: config.Server.Timeout,
		IdleTimeout:  config.Server.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("failed to start server")
		}
	}()
	<-ctx.Done()
	logger.Info("stopping server")
	time.Sleep(time.Second * 5)
}

func setupLogger(env string) *slog.Logger {
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
