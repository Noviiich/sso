package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/Noviiich/sso/internal/app/grps"
	"github.com/Noviiich/sso/internal/services/auth"
	"github.com/Noviiich/sso/internal/storage/sqlite"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	storagePath string,
	port int,
	tokenTTL time.Duration,

) *App {

	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, port)

	return &App{
		GRPCServer: grpcApp,
	}
}
