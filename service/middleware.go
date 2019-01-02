package service

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/liuyh73/dailyhub.service/db"
)

var permission = []string{
	"/api",
	"/api/register",
	"/api/login",
}

func permit(uri string) bool {
	for _, u := range permission {
		if u == uri {
			return true
		}
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
			dh_token := ""
			for k, v := range r.Header {
				if strings.ToLower(k) == TokenName {
					dh_token = v[0]
					break
				}
			}
			mapClaims, err := parseToken(dh_token, []byte(SecretKey))
			checkErr(err)
			if err != nil || mapClaims.Valid() != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(writeResp(false, "Unauthorized access to this resource", Token{}))
				return
			}
			log.Println(mapClaims["username"])
			has, err, tokenItem := db.GetUserTokenItem(mapClaims["username"].(string))
			checkErr(err)
			if !has || err != nil || tokenItem.DH_TOKEN != dh_token {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(writeResp(false, "Unauthorized access to this resource", Token{}))
			} else {
				ctx := context.WithValue(r.Context(), "username", mapClaims["username"])
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
