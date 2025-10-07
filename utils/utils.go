package utils

import (
	"errors"
	"gingorm/config"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// 哈希密码
func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	return string(hash), err
}

// 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil

}

// 生成JWT
func GenerateJWT(userID int, username string) (string, error) {
	expireSeconds := config.AppConfig.JWT.Expire
	secret := config.AppConfig.JWT.Secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   userID,
		"username": username,
		"exp":      time.Now().Add(time.Duration(expireSeconds) * time.Second).Unix(),
	})
	signedToken, err := token.SignedString([]byte(secret))
	return "Bearer " + signedToken, err
}

// 解析JWT
func ParseJWT(tokenString string) (int, string, error) {
	if tokenString[:7] == "Bearer " && len(tokenString) > 7 {
		tokenString = tokenString[7:]
	}
	secret := config.AppConfig.JWT.Secret
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected Signing Method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		userID := int(claims["userID"].(float64))
		if !ok {
			return 0, "", errors.New("username claims is not a string")
		}
		return userID, username, nil
	}
	return 0, "", err
}
