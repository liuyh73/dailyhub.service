package service

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const TokenName = "dh_token"
const Issuer = "sysu-615"
const SecretKey = "DailyHub"

type Token struct {
	DH_TOKEN string `json:"dh_token"`
}

type jwtCustomClaims struct {
	jwt.StandardClaims

	Username string `json:"username"`
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func createToken(secretkey []byte, issuer string, username string) (token Token, err error) {
	claims := &jwtCustomClaims{
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Add(time.Hour * 24 * 365 * 100).Unix()),
			Issuer:    issuer,
		},
		username,
	}
	dh_token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretkey)
	token = Token{
		dh_token,
	}
	return
}

func parseToken(dh_token string, secretKey []byte) (claims jwt.MapClaims, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(dh_token, func(*jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	// fmt.Println(token.Claims)
	claims = token.Claims.(jwt.MapClaims)
	return
}
