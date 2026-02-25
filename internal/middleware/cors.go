package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		fmt.Printf("[%s] %s %s %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.RequestURI,
			duration)
	})
}

func ApplyMiddlewares(r *mux.Router, jwt *JWTMiddleware) {
	r.Use(CORSMiddleware)
	r.Use(LoggingMiddleware)
	r.Use(jwt.AuthMiddleware)
}
