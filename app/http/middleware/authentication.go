package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kenshaw/envcfg"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type contextKey string

const ctxKey contextKey = "user"

type JwtMiddleware struct {
	SecretKey      string
	UserContextKey string
}

func NewJwtMiddleware(secret string) *JwtMiddleware {
	return &JwtMiddleware{
		SecretKey: secret,
	}
}

func jwtSecret() string {
	cfg, err := envcfg.New(envcfg.ConfigFile("config"))
	if err != nil {
		panic(err)
	}

	secret := cfg.GetString("jwt.jwtSecret")
	if cfg.GetString("jwt.jwtSecret") == "" {
		panic("jwt secret is not set")
	}

	return secret
}

type JWTCustomClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID string) (string, string, int64, error) {
	accessTokenClaims := &JWTCustomClaims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	refreshTokenClaims := &JWTCustomClaims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	expiredAtStr := jwt.NewNumericDate(time.Now().Add(time.Minute * 30)).Unix()
	// expiredAtStr := jwt.NewNumericDate(time.Now().Add(time.Minute * 30)).Format("2006-01-02 15:04:05")

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	jwtAccessToken, err := accessToken.SignedString([]byte(jwtSecret()))
	if err != nil {
		return "", "", 0, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	jwtRefreshToken, err := refreshToken.SignedString([]byte(jwtSecret()))
	if err != nil {
		return "", "", 0, err
	}

	return jwtAccessToken, jwtRefreshToken, expiredAtStr, err
}

func ParseJWT(tokenString string) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret()), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func GetAuthUser(c echo.Context) *JWTCustomClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTCustomClaims)

	return claims
}

func JWTAuth() echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JWTCustomClaims)
		},
		SigningKey: []byte(jwtSecret()),
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"status":     "Unauthorized",
				"statusCode": 401,
				"message":    "Invalid or expired token",
			})
		},
	}

	return echojwt.WithConfig(config)
}

func (m *JwtMiddleware) SetUserContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userToken := c.Get("user").(*jwt.Token)
		claims, ok := userToken.Claims.(*JWTCustomClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		// Simpan user di context
		ctx := context.WithValue(c.Request().Context(), ctxKey, map[string]string{
			"userId": claims.UserID,
		})
		// Update request context di Echo
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
