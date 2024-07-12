package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(idUser uint) (string, error) {
	var claim = jwt.MapClaims{}
	claim["id"] = idUser
	claim["iat"] = time.Now().UnixMilli()
	claim["exp"] = time.Now().Add(time.Minute * 2).UnixMilli()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// token := jwt.New(jwt.SigningMethodHS256)
	strToken, err := token.SignedString([]byte("$!1gnK3yyy!!!"))
	if err != nil {
		return "", err
	}

	return strToken, nil
}

func ExtractToken(t *jwt.Token) (uint, error) {
	var userID uint
	// expiredTime, err := t.Claims.GetExpirationTime()
	// if err != nil {
	// 	return 0, err
	// }
	// issuedTime, err := t.Claims.GetIssuedAt()
	// if err != nil {
	// 	return 0, err
	// }
	if t.Valid {
		var tokenClaims = t.Claims.(jwt.MapClaims)
		userID = uint(tokenClaims["id"].(float64))

		return userID, nil
	}

	return 0, errors.New("token tidak valid")
}
