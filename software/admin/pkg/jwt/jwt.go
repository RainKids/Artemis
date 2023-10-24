package jwt

import (
	"admin/global"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

// MySecret 定义JWT TOKEN的加密盐
var MySecret = []byte(global.MySignedKey)

type MyClaims struct {
	Username string `json:"username"`
	UserId   int    `json:"userId"`
	jwt.RegisteredClaims
}

func CreateToken(username string, userId int) (string, time.Time, error) {
	expiresTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &MyClaims{
		Username:         username,
		UserId:           userId,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expiresTime)},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(MySecret)
	if err != nil {
		return "", expiresTime, err
	}
	return token, expiresTime, nil
}

func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
