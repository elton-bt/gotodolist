package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/elton-bt/gotodolist/internal/config"
	"github.com/elton-bt/gotodolist/internal/postgres"
	appRuntime "github.com/elton-bt/gotodolist/internal/runtime"
	"github.com/elton-bt/gotodolist/internal/todo"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.Load("gotodolist-monolito", "8080")
	if err != nil {
		logger.Error("configuracao invalida")
		os.Exit(1)
	}

	ctx := context.Background()
	store, err := postgres.NewStore(ctx, cfg.DatabaseURL())
	if err != nil {
		logger.Error("banco indisponivel")
		os.Exit(1)
	}
	defer store.Close()

	if err := store.EnsureSchema(ctx); err != nil {
		logger.Error("falha ao inicializar schema do banco")
		os.Exit(1)
	}

	service := todo.NewService(store)
	server, err := newHTTPServer(cfg, service)
	if err != nil {
		logger.Error("falha ao configurar servidor")
		os.Exit(1)
	}

	logger.Info("servidor iniciado", "app", cfg.AppName, "addr", server.Addr)
	if err := appRuntime.Run(server, cfg.ShutdownTimeout); err != nil {
		logger.Error("servidor encerrado com erro")
		os.Exit(1)
	}

	logger.Info("servidor encerrado")
}
