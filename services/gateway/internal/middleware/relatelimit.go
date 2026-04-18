package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/httprate"
)

func RelateLimit() func(http.Handler) http.Handler {
	//TODO: The default value should be 100, but we'll set it to 1000 for testing purposes
	limit := 100
	if os.Getenv("ENV") == "development" {
		limit = 1000
	}
	return httprate.LimitByIP(limit, time.Minute)
}
