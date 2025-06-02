package jwt

import (
	"time"

	"github.com/Noviiich/sso/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

// NewToken генерирует новый JWT токен для пользователя и приложения с заданной временем жизни
func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// payload
	claims := token.Claims.(jwt.MapClaims) //преобразование Claims в карту для удобства работы с ключами
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["app_id"] = app.ID
	claims["exp"] = time.Now().Add(duration).Unix() // Срок действия токена

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
