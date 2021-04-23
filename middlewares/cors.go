package middlewares

import (
	"net/http"
)

func SetCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Auth-Token")

		next.ServeHTTP(w, req)
	})
}
