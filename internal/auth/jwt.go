package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const (
	cookieName = "Auth"
	ttl        = time.Hour * 0
)

type Claims struct {
	jwt.RegisteredClaims
	UserUUID string
}

type UUIDGenerator interface {
	Generate() (string, error)
}

func NewAuthJWTMiddleware(uuidGenerator *UUIDGenerator, secretKey []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var userUUID string
			token, err := c.Cookie(cookieName)
			if err != nil && !errors.Is(err, http.ErrNoCookie) {
				return c.String(http.StatusInternalServerError, fmt.Sprintf("AuthJWTMiddleware c.Cookie %v", err))
			}
			if err == nil {
				userUUID = getUserUUID(token.Value, secretKey)
				if userUUID != "" {
					c.Set("authorized", true)
				}
			}

			if userUUID == "" {
				userUUID, err = (*uuidGenerator).Generate()
				if err != nil {
					return c.String(http.StatusInternalServerError, fmt.Sprintf("AuthJWTMiddleware uuidGenerator.Generate %v", err))
				}
			}

			c.Set("userUUID", userUUID)

			// create or refresh cookie with token
			err = setTokenCookie(c, userUUID, secretKey, ttl)
			if err != nil {
				return c.String(http.StatusInternalServerError, fmt.Sprintf("AuthJWTMiddleware setTokenCookie %v", err))
			}

			return next(c)
		}
	}
}

func getUserUUID(tokenString string, secretKey []byte) string {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return ""
	}

	if !token.Valid {
		return ""
	}
	return claims.UserUUID
}

func generateJWTTokenString(userUUID string, secretKey []byte, ttl time.Duration) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{},
		UserUUID:         userUUID,
	}

	if ttl > 0 { // allow to be infinite if ttl == 0
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(ttl))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func newAuthCookie(token string, ttl time.Duration) *http.Cookie {
	c := &http.Cookie{
		Name:     cookieName,
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	}
	if ttl > 0 { // allow to be infinite if ttl == 0
		c.Expires = time.Now().Add(ttl)
	}
	return c
}

func setTokenCookie(c echo.Context, userUUID string, secretKey []byte, ttl time.Duration) error {
	tokenString, err := generateJWTTokenString(userUUID, secretKey, ttl)
	if err != nil {
		return fmt.Errorf("generateJWTTokenString %v", err)
	}
	token := newAuthCookie(tokenString, ttl)
	c.SetCookie(token)
	return nil
}
