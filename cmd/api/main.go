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

	categoryRepo := repository.NewCategoryRepository(pool)
	menuItemRepo := repository.NewMenuItemRepository(pool)
	menuService := services.NewMenuService(categoryRepo, menuItemRepo)
	menuHandler := handlers.NewMenuHandler(menuService)

	adminOnly := middleware.RequireRole("admin")

	// logged-in users
	mux.Handle("GET /api/categories",
		middleware.Chain(http.HandlerFunc(menuHandler.GetCategories), auth))
	mux.Handle("GET /api/menu-items",
		middleware.Chain(http.HandlerFunc(menuHandler.GetItems), auth))

	// admin only
	mux.Handle("POST /api/categories",
		middleware.Chain(http.HandlerFunc(menuHandler.CreateCategories), auth, adminOnly))
	mux.Handle("POST /api/menu-items",
		middleware.Chain(http.HandlerFunc(menuHandler.CreateItem), auth, adminOnly))
	mux.Handle("PATCH /api/menu-items/{id}",
		middleware.Chain(http.HandlerFunc(menuHandler.UpdateItem), auth, adminOnly))

	// admin OR staff = any authenticated user (you only have two roles)
	mux.Handle("PATCH /api/menu-items/{id}/availability",
		middleware.Chain(http.HandlerFunc(menuHandler.SetAvailability), auth))

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	log.Printf("listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
