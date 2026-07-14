package helpers

import (
	"fmt"
	"los_andes/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var accessTokenDuration = 10 * time.Hour

func getSecretKey() []byte {

	_ = godotenv.Load()

	key := os.Getenv("Secret_Key")

	if key == "" {
		fmt.Println("⚠️  WARNING: SECRET_KEY no está definido en las variables de entorno.")
	}

	return []byte(key)
}

func GenerateToken(user models.UsuarioJWT) (string, error) {

	var jwtSecret = []byte(os.Getenv("Secret_Key"))

	claims := models.CustomClaims{
		UserId: user.UsuarioId,
		Name:   user.Nombres,
		Rol:    user.Rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    os.Getenv("URL"),
			Audience:  jwt.ClaimStrings{os.Getenv("URL")},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseAndVerifyToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algoritmo de firma inválido")
		}
		return getSecretKey(), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token inválido: %w", err)
	}

	claims, ok := token.Claims.(*models.CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("claims inválidos o token inválido")
	}

	return claims, nil
}
