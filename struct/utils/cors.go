package utils

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func CorsMiddleware() mux.MiddlewareFunc {
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"},           // Allow requests from any origin
		AllowedMethods:   []string{"GET", "POST"}, // Allow GET and POST requests
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true, // Allow sending of credentials (cookies, headers)
	}

	return cors.New(corsOptions).Handler
}
