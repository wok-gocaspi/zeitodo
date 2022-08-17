package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
)

func GenerateToken(userid primitive.ObjectID) string {
	claims := jwt.MapClaims{
		"exp":    time.Now().Add(time.Minute * 5).Unix(),
		"iat":    time.Now().Unix(),
		"userID": userid.Hex(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return t
}

func ValidateToken(token string) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	jwtoken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, nil, err
	}
	return jwtoken, claims, nil
}
