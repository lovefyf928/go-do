package authorization

import (
	"github.com/dgrijalva/jwt-go"
	"go-do/common/conf"
	"time"
)

var TOKEN_HEADER_NAME = ""
var SECRET_KEY = ""

func LoadJwtConfig() {
	TOKEN_HEADER_NAME = conf.ConfigInfo.Jwt.TokenHeaderName
	SECRET_KEY = conf.ConfigInfo.Jwt.SecretKey
}

type UserClaims struct {
	jwt.StandardClaims

	TokenType string `json:"type"`

	MachineId string `json:"machineId"`

	UserId string `json:"uid"`
}

func BuildUserToken(SecretKey []byte, issuer string, tokenType string, machineId string, uid string) (tokenString string, err error) {
	claims := &UserClaims{
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Add(time.Hour * 72).Unix()),
			Issuer:    issuer,
		},
		tokenType,
		machineId,
		uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(SecretKey)
	return
}

func ParseUserToken(tokenSrt string, SecretKey []byte) (claims jwt.Claims, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(tokenSrt, func(*jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	claims = token.Claims

	return
}

func getUserIdByToken(tokenSrt string) (userId string, err error) {
	claims, err := ParseUserToken(tokenSrt, []byte(SECRET_KEY))
	if err == nil {
		userClaims, ok := claims.(UserClaims)
		if ok {
			userId = userClaims.UserId
			return
		}
	}
	return
}
