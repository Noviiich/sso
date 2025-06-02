package auth

import (
	"context"

	ssov1 "github.com/Noviiich/sso/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer //заглушка для методов, которые будут реализованы в будущем для запуска gRPC сервера
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	// Обработка ошибка осуществляется с помощью status.Error,  для того, чтобы формат ошибки был понятен любому grpc-клиенту
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	// Здесь должна быть логика аутентификации пользователя
	return &ssov1.LoginResponse{
		Token: "example_token", // Здесь должен быть реальный токен, полученный после аутентификации
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Здесь должна быть логика регистрации пользователя
	return &ssov1.RegisterResponse{UserId: 1}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Здесь должна быть логика проверки, является ли пользователь администратором
	return &ssov1.IsAdminResponse{
		IsAdmin: true,
	}, nil
}
