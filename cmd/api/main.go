package main

import (
	"context"
	"log"
	"net/http"

	"github.com/m15z/restaurantops-api/internal/config"
	"github.com/m15z/restaurantops-api/internal/database"
	"github.com/m15z/restaurantops-api/internal/handlers"
	"github.com/m15z/restaurantops-api/internal/middleware"
	"github.com/m15z/restaurantops-api/internal/repository"
	"github.com/m15z/restaurantops-api/internal/services"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	auth := middleware.Auth(cfg.JWTSecret)           // call layer 1 → get wrapper
	protected := auth(http.HandlerFunc(handlers.Me)) // call layer 2 → get guarded handler

	// wiring: repo → service → handler
	userRepo := repository.NewUserRepository(pool)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)

	healthHandler := handlers.NewHealthHandler(pool)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", healthHandler.Check)

	authHandler := handlers.NewAuthHandler(authService)
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.Handle("GET /api/me", protected) // register layer 3

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	log.Printf("listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
