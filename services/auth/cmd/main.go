package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	dbpkg "copo/auth/internal/db"
	"copo/auth/internal/handler"
	"copo/auth/internal/repository"
	"copo/auth/internal/router"
	"copo/auth/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db, err := dbpkg.ConnectDB()
	if err != nil {
		log.Fatalf("unable to connect Postgres: %v", err)
	}
	defer db.Close()

	//Repository
	userRepo := repository.NewRepository(db)

	//Service
	authSvc := service.NewAuthService(userRepo)
	userSvc := service.NewUserService(userRepo)

	//Handler
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)

	r := router.New(authHandler, userHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Auth service running on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("error shutting down server: %v", err)
	}

	log.Println("Auth service stopped")

}
