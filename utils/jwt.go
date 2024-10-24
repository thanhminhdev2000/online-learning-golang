package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(userId int, role string, expirationTime time.Duration) (string, int64, error) {
	expiration := time.Now().Add(expirationTime)
	var jwtKey = []byte(os.Getenv("JWT_KEY"))

	claims := jwt.MapClaims{
		"userId": strconv.Itoa(userId),
		"role":   role,
		"exp":    expiration.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", 0, err
	}
	return tokenString, int64(expirationTime.Seconds()), nil
}

func CreateAccessToken(userId int, role string) (string, int64, error) {
	return CreateToken(userId, role, 1*time.Hour)
}

func CreateRefreshToken(userId int, role string) (string, int64, error) {
	return CreateToken(userId, role, 7*24*time.Hour)
}

func ValidToken(tokenString string) (int, string, error) {
	var jwtKey = []byte(os.Getenv("JWT_KEY"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return 0, "", fmt.Errorf("token has expired")
			}
		}

		// Kiểm tra và chuyển đổi kiểu dữ liệu của userId
		var userId int
		switch id := claims["userId"].(type) {
		case string:
			userIdInt, err := strconv.Atoi(id)
			if err != nil {
				return 0, "", fmt.Errorf("invalid userId in token")
			}
			userId = userIdInt
		case float64:
			userId = int(id)
		default:
			return 0, "", fmt.Errorf("userId not found in token")
		}

		// Kiểm tra role
		var role string
		if roleVal, ok := claims["role"].(string); ok {
			role = roleVal
		} else {
			return 0, "", fmt.Errorf("role not found in token")
		}

		return userId, role, nil
	}

	return 0, "", fmt.Errorf("invalid token")
}
