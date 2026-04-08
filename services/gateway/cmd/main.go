package main

import (
	"copo/gateway/router"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	r := router.New()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("API Gateway corriendo :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
