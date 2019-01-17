package service

import (
	"context"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/liuyh73/dailyhub.service/db"
)

var permission = []string{
	"/api/register",
	"/api/login",
	"/api/users",
}

func permit(uri string) bool {
	for _, u := range permission {
		if strings.HasPrefix(uri, u) {
			return true
		}
	}
	if uri == "/api" {
		return true
	}
	return false
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, sw_token,sign")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Content-Type", "application/json")
		if !permit(r.RequestURI) {
			var mapClaims jwt.MapClaims
			dh_token := ""
			// 如果token存在于Authorization中
			token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
				return []byte(SecretKey), nil
			})
			checkErr(err)
			if token != nil {
				var ok bool
				mapClaims, ok = token.Claims.(jwt.MapClaims)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write(writeResp(false, "Unauthorized access to this resource", nil))
					return
				}
				dh_token = strings.Split(r.Header["Authorization"][0], " ")[1]
			} else {
				// 如果token存在于header中
				for k, v := range r.Header {
					if strings.ToLower(k) == TokenName {
						dh_token = v[0]
						break
					}
				}

				if dh_token == "" {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write(writeResp(false, "Unauthorized access to this resource", nil))
					return
				}

				mapClaims, err = parseToken(dh_token, []byte(SecretKey))
				checkErr(err)
				if err != nil || mapClaims.Valid() != nil {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write(writeResp(false, "Unauthorized access to this resource", nil))
					return
				}
			}
			// log.Println(mapClaims["username"])
			has, err, tokenItem := db.GetUserTokenItem(mapClaims["username"].(string))
			checkErr(err)
			if !has || err != nil || tokenItem.DH_TOKEN != dh_token {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(writeResp(false, "Unauthorized access to this resource", nil))
			} else {
				ctx := context.WithValue(r.Context(), "username", mapClaims["username"])
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
