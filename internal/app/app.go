package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/Noviiich/sso/internal/app/grps"
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

	grpcApp := grpcapp.New(log, port)

	return &App{
		GRPCServer: grpcApp,
	}
}
