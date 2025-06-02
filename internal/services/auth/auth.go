package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Noviiich/sso/internal/domain/models"
	"github.com/Noviiich/sso/internal/lib/jwt"
	"github.com/Noviiich/sso/internal/lib/logger/sl"
	"github.com/Noviiich/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials") //
)

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

func New(
	log *slog.Logger,
	usrSaver UserSaver,
	usrProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    usrSaver,
		usrProvider: usrProvider,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

// RegisterNewUser регистрирует нового пользователя в системе и возвращает его ID.
// Если регистрация не удалась, возвращает ошибку.
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	pass string,
) (int64, error) {
	// op - метка для логирования
	const op = "Auth.RegisterNewUser"

	// Создание контекста для логирования
	log := a.log.With(
		slog.String("op", op),       //имя операции
		slog.String("email", email), // email пользователя
	)

	log.Info("registering new user")

	// генерация хэша и соли для пароля
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Созраняем пароль в БД
	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string, // пароль в открытом виде
	appID int, // ID приложения, в которое пользователь хочет войти
) (string, error) {
	const op = "Auth.Login"

	// Создание контекста для логирования
	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("user login attempt")

	// Получаем пользователя из БД
	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		a.log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Проверям корректность полученного пароля
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Warn("invalid password", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	// Получение информации о приложении
	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	// Генерация JWT токена
	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("token generated successfully", slog.String("token", token))
	return token, nil
}

func (a *Auth) IsAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	const op = "Auth.IsAdmin"

	// Создание контекста для логирования
	log := a.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("checking if user is admin")

	// Проверяем, является ли пользователь администратором
	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		log.Error("failed to check if user is admin", sl.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
