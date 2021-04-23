package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/anand-kay/linkedin-clone/utils"

	jwt "github.com/dgrijalva/jwt-go"
)

// Authorize - Authorizes the user
func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		tokenStr := req.Header["Auth-Token"][0]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte("secret"), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authorization failed"))
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authorization failed"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authorization failed"))
			return
		}

		ctx1 := context.WithValue(req.Context(), utils.ContextUserIDKey, strconv.FormatFloat(claims["UserID"].(float64), 'f', -1, 64))
		ctx2 := context.WithValue(ctx1, utils.ContextEmailKey, claims["Email"].(string))
		ctx3 := context.WithValue(ctx2, utils.ContextFirstNameKey, claims["FirstName"].(string))
		ctx4 := context.WithValue(ctx3, utils.ContextLastNameKey, claims["LastName"].(string))

		next.ServeHTTP(w, req.WithContext(ctx4))
	})
}
