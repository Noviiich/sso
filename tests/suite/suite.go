package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/Noviiich/sso/internal/config"
	ssov1 "github.com/Noviiich/sso/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/local_tests.yaml")

	ctx, canceCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	// закрытие контекста
	t.Cleanup(func() {
		t.Helper()
		canceCtx()
	})

	// Создаем gRPC клиент
	cc, err := grpc.DialContext(
		context.Background(),
		grpcAddress(cfg),
		// insecure-коннект для тестов
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to connect to gRPC server: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
