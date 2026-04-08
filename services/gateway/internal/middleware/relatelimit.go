package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

func RelateLimit() func(http.Handler) http.Handler {
	return httprate.LimitByIP(100, time.Minute)
}
